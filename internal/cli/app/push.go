package app

import (
	"errors"
	"fmt"
	"github.com/lodestar-cli/lodestar/internal/cli/file"
	"github.com/lodestar-cli/lodestar/internal/cli/home"
	"github.com/lodestar-cli/lodestar/internal/common/auth"
	"github.com/lodestar-cli/lodestar/internal/common/environment"
	"github.com/lodestar-cli/lodestar/internal/common/remote"
	"strings"
	"sync"
)

type PushCliOptions struct{
	Username        string
	Token           string
	App             string
	ConfigPath      string
	EnvironmentName string
	YamlKeys        string
}

type Push struct {
	CliOptions           PushCliOptions
	KeysMap              map[string]string
	GitAuth              auth.GitCredentials
	Environment          *environment.Environment
	Repository           *remote.LodestarRepository
	AppConfigurationFile *file.AppConfigurationFile
	AppStateFile         *remote.AppStateFile
}

func NewPush(username string, token string, app string, configPath string, environment string, yamlKeys string) (*Push, error){
	var err error
	var wg sync.WaitGroup
	fatalErrors := make(chan error)
	finish := make(chan bool)
	defer close(fatalErrors)
	defer close(finish)

	cli := PushCliOptions{
		Username: username,
		Token: token,
		App: app,
		ConfigPath: configPath,
		EnvironmentName: environment,
		YamlKeys: yamlKeys,
	}

	p := Push{
		CliOptions: cli,
	}

	//1. get app config file
	err = p.setAppConfigurationFile()
	if err != nil {
		return nil, err
	}

	//2. Get Environment for AppConfig as well as Auth.  Check to make sure keys given match the ones in the AppConfig file
	wg.Add(3)

	go p.setEnvironment(fatalErrors, &wg)
	go p.setAuth(fatalErrors, &wg)
	go p.compareKeys(fatalErrors, &wg)
	go func(){
		wg.Wait()
		finish <- true
	}()

	select {
	case <- finish:
		break
	case err = <-fatalErrors:
		return nil, err
	}

	//3. Clone Manifest Repository
	fmt.Printf("Cloning %s as %s...\n", p.AppConfigurationFile.Info.RepoUrl, p.CliOptions.Username)
	p.Repository, err = remote.NewLodestarRepository(p.AppConfigurationFile.Info.RepoUrl, p.GitAuth)

	//4. Fetch App State File from Repository
	p.AppStateFile, err = remote.NewAppStateFile(p.Repository,p.AppConfigurationFile.Info.StatePath, p.AppConfigurationFile.Info.Name)
	if err != nil {
		return nil, err
	}

	return &p, nil

}

func (p *Push) Execute() error {
	var err error
	var wg sync.WaitGroup
	fatalErrors := make(chan error)
	finish := make(chan bool)
	defer close(fatalErrors)
	defer close(finish)

	//1. Update StateGraph and ManagementFiles
	fileChannel := make(chan remote.LodestarFile, 2)
	var updatedFiles []remote.LodestarFile
	wg.Add(2)

	go p.updateAppStateFile(fatalErrors, fileChannel, &wg)
	go p.updateManagementFile(fatalErrors, fileChannel, &wg)

	go func(){
		wg.Wait()
		close(fileChannel)
		finish <- true
	}()

	select {
	case <- finish:
		for f := range fileChannel{
			updatedFiles = append(updatedFiles, f)
		}

		switch len(updatedFiles) {
		case 0:
			fmt.Printf("%s environmnet's state and management files are up to date!", p.Environment.Name)
		case 1:
			fmt.Printf("WARNING: %s environmnet's state and management files were out of sync. Syncing files to newest push", p.Environment.Name)
			err = p.Repository.CommitFiles(fmt.Sprintf("Lodestar updated %v in %s environment", p.AppConfigurationFile.YamlKeys, p.Environment.Name), updatedFiles...)
			if err != nil{
				return err
			}
			fmt.Printf("Pushing changes to %s as %s...\n",p.AppConfigurationFile.Info.RepoUrl,p.CliOptions.Username)
			err = p.Repository.Push()
			if err != nil {
				return err
			}
			fmt.Println("Push complete!")
		default:
			err = p.Repository.CommitFiles(fmt.Sprintf("Lodestar updated %v in %s environment", p.AppConfigurationFile.YamlKeys, p.Environment.Name), updatedFiles...)
			if err != nil{
				return err
			}
			fmt.Printf("Pushing changes to %s as %s...\n",p.AppConfigurationFile.Info.RepoUrl,p.CliOptions.Username)
			err = p.Repository.Push()
			if err != nil {
				return err
			}
			fmt.Println("Push complete!")
		}

	case err = <-fatalErrors:
		return err
	}

	return nil
}

func (p *Push) Output(b bool) error{
	if b{
		err := p.AppStateFile.Output()
		if err != nil {
			return err
		}
	}

	p.AppStateFile.Print()

	return nil
}

func (p *Push) setAppConfigurationFile() error {
	if p.CliOptions.App == "" && p.CliOptions.ConfigPath == "" {
		return errors.New("must provide an App name or a path to a configuration file. For more information, run: lodestar app push --help")
	}else if p.CliOptions.ConfigPath != "" {
		var err error
		p.AppConfigurationFile, err = file.NewAppConfigurationFile(p.CliOptions.ConfigPath)
		if err != nil {
			return err
		}
	} else {
		path, err := home.GetPath("app", p.CliOptions.App)
		if err != nil {
			return err
		}
		p.AppConfigurationFile, err = file.NewAppConfigurationFile(path)
	}
	return nil
}

func (p *Push) setEnvironment(fatalErrors chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	p.Environment, err = p.AppConfigurationFile.GetEnvironment(p.CliOptions.EnvironmentName)
	if err != nil {
		fatalErrors <- err
		return
	}
}

func (p *Push) compareKeys(fatalErrors chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	keys := p.setKeyMap(fatalErrors)
	if len(keys) == 0 {
		return
	}

	for _, k := range keys{
		b := false
		for _ , i := range p.AppConfigurationFile.YamlKeys{
			if i == k{
				b = true
			}
		}
		if b {
			continue
		}
		fatalErrors <- errors.New("a key pair supplied had a key that does not match the keys in the AppConfig File.  Please make sure all key pairs match.")
		return
	}
}

func (p *Push) setKeyMap(fatalErrors chan error) []string {
	s := strings.Split( p.CliOptions.YamlKeys, ",")
	var keys []string

	for _, keyPair := range s {
		kv := strings.Split( keyPair, "=")

		if len(kv) != 2 {
			fatalErrors <- errors.New("yamlKeys are incorrectly formatted.  Must be a single string of <key>=<value> pairs delimited by commas")
			return nil
		}

		p.KeysMap[kv[0]]=kv[1]
		keys = append(keys, kv[0])
	}

	return keys
}

func (p *Push) setAuth(fatalErrors chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	var a auth.GitCredentials
	if p.CliOptions.Username == "" {
		fatalErrors <- errors.New("username is not set. Try lodestar app push --help")
	} else if p.CliOptions.Token == "" {
		fatalErrors <- errors.New("token is not set. Try lodestar app push --help")
	} else {
		a = &auth.TokenCredentials{
			Username: p.CliOptions.Username,
			Token: p.CliOptions.Token,
		}
		p.GitAuth = a
	}
	//need to add option that creates ssh credentials
}

func (p *Push) updateAppStateFile(fatalErrors chan error, fileChannel chan remote.LodestarFile, wg *sync.WaitGroup) {
	defer wg.Done()
	var l remote.LodestarFile
	updated, err := p.AppStateFile.UpdateEnvironmentGraph(p.Environment.Name,p.KeysMap)
	if err != nil{
		fatalErrors <- err
	}
	if updated{
		err = p.AppStateFile.UpdateFile()
		if err != nil {
			fatalErrors <- err
		}
		l = p.AppStateFile
		fileChannel <- l
	}
}

func (p *Push) updateManagementFile(fatalErrors chan error, fileChannel chan remote.LodestarFile, wg *sync.WaitGroup) {
	defer wg.Done()
	var l remote.LodestarFile

	fmt.Printf("Retrieving environment configuration at %s...\n", p.Environment.SrcPath)

	m, err := remote.NewManagementFile(p.Environment, p.Repository)
	if err != nil{
		fatalErrors <- err
		return
	}

	updated, err := m.UpdateFileContents(p.KeysMap)
	if err != nil{
		fatalErrors <- err
	}
	if updated{
		l = m
		fileChannel <- l
	}
}