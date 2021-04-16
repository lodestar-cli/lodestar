package app

import (
	"errors"
	"fmt"
	"github.com/lodestar-cli/lodestar/internal/cli/files"
	"github.com/lodestar-cli/lodestar/internal/cli/home"
	"github.com/lodestar-cli/lodestar/internal/common/auth"
	"github.com/lodestar-cli/lodestar/internal/common/environment"
	"github.com/lodestar-cli/lodestar/internal/common/repo"
	"sync"
)

type PromoteCliOptions struct{
	Username             string
	Token                string
	App                  string
	ConfigPath           string
	SrcEnvironment       string
	DestEnvironment      string
}

type Promote struct{
	CliOptions PromoteCliOptions
	GitAuth              auth.GitCredentials
	SrcEnvironment       *environment.Environment
	DestEnvironment      *environment.Environment
	Repository           *repo.LodestarRepository
	AppConfigurationFile *files.AppConfigurationFile
	AppStateFile         *files.AppStateFile
}

func NewPromote(username string, token string, app string, configPath string, srcEnv string, destEnv string) (*Promote,error){
	var err error
	var wg sync.WaitGroup
	fatalErrors := make(chan error)
	finish := make(chan bool)
	defer close(fatalErrors)
	defer close(finish)

	cli := PromoteCliOptions{
		Username: username,
		Token: token,
		App: app,
		SrcEnvironment: srcEnv,
		DestEnvironment: destEnv,
		ConfigPath: configPath,
	}

	p := Promote{
		CliOptions: cli,
	}

	//1. get app config file
	err = p.setAppConfigurationFile()
	if err != nil {
		return nil, err
	}

	//2. Get Environments from AppConfig as well as Auth.  Check to make sure keys given match the ones in the AppConfig file
	wg.Add(3)
	go p.setEnvironment(fatalErrors, "src",&wg)
	go p.setEnvironment(fatalErrors, "dest",&wg)
	go p.setAuth(fatalErrors, &wg)
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
	p.Repository, err = repo.NewLodestarRepository(p.AppConfigurationFile.Info.RepoUrl, p.GitAuth)

	//4. Fetch App State File from Repository
	p.AppStateFile, err = files.NewAppStateFile(p.Repository,p.AppConfigurationFile.Info.StatePath, p.AppConfigurationFile.Info.Name)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Promote) Execute() error {
	fatalErrors := make(chan error)
	finish := make(chan bool)
	defer close(fatalErrors)
	defer close(finish)
	var wg sync.WaitGroup
	var err error
	var keysMap map[string]string

	fmt.Printf("Retrieving key values from configuration file %s...\n",p.SrcEnvironment.Name)

	smf, err := files.NewManagementFile(p.SrcEnvironment,p.Repository)
	if err != nil{
		return err
	}

	keysMap, err = smf.GetKeyValues(p.AppConfigurationFile.YamlKeys)
	if err != nil{
		return err
	}

	//1. Update StateGraph and ManagementFiles
	fileChannel := make(chan files.LodestarFile, 2)
	var updatedFiles []files.LodestarFile
	wg.Add(2)

	fmt.Printf("Updating %s environment to %s environment's keys", p.DestEnvironment.Name, p.SrcEnvironment.Name)
	go p.updateAppStateFile(fatalErrors, fileChannel, keysMap, &wg)
	go p.updateManagementFile(fatalErrors, fileChannel,smf, keysMap, &wg)

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
			fmt.Printf("%s environmnet's state and management files are up to date!", p.DestEnvironment.Name)
		case 1:
			fmt.Printf("WARNING: %s environmnet's state and management files were out of sync. Syncing files to newest push", p.DestEnvironment.Name)
			err = p.Repository.CommitFiles(fmt.Sprintf("Lodestar updated %v in %s environment", p.AppConfigurationFile.YamlKeys, p.DestEnvironment.Name), updatedFiles...)
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
			err = p.Repository.CommitFiles(fmt.Sprintf("Lodestar updated %v in %s environment", p.AppConfigurationFile.YamlKeys, p.DestEnvironment.Name), updatedFiles...)
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

func (p *Promote) Output(b bool) error{
	if b{
		err := p.AppStateFile.Output()
		if err != nil {
			return err
		}
	}

	p.AppStateFile.Print()

	return nil
}

func (p *Promote) updateAppStateFile(fatalErrors chan error, fileChannel chan files.LodestarFile, keysMap map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	var l files.LodestarFile
	updated, err := p.AppStateFile.UpdateEnvironmentGraph(p.DestEnvironment.Name, keysMap)
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

func (p *Promote) updateManagementFile(fatalErrors chan error, fileChannel chan files.LodestarFile, smf *files.ManagementFile, keysMap map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	var l files.LodestarFile

	updated, err := smf.UpdateFileContents(keysMap)
	if err != nil{
		fatalErrors <- err
	}
	if updated{
		l = smf
		fileChannel <- l
	}
}


func (p *Promote) setAppConfigurationFile() error {
	if p.CliOptions.App == "" && p.CliOptions.ConfigPath == "" {
		return errors.New("must provide an App name or a path to a configuration file. For more information, run: lodestar app push --help")
	}else if p.CliOptions.ConfigPath != "" {
		var err error
		p.AppConfigurationFile, err = files.NewAppConfigurationFile(p.CliOptions.ConfigPath)
		if err != nil {
			return err
		}
	} else {
		path, err := home.GetPath("app", p.CliOptions.App)
		if err != nil {
			return err
		}
		p.AppConfigurationFile, err = files.NewAppConfigurationFile(path)
	}
	return nil
}

func (p *Promote) setEnvironment(fatalErrors chan error, env string, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	switch env {
	case "src":
		p.SrcEnvironment, err = p.AppConfigurationFile.GetEnvironment(p.CliOptions.SrcEnvironment)
		if err != nil {
			fatalErrors <- err
		}
	case "dest":
		p.DestEnvironment, err = p.AppConfigurationFile.GetEnvironment(p.CliOptions.DestEnvironment)
		if err != nil {
			fatalErrors <- err
		}
	}
}

func (p *Promote) setAuth(fatalErrors chan error, wg *sync.WaitGroup) {
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