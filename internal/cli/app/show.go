package app

import (
	"errors"
	"fmt"
	"github.com/lodestar-cli/lodestar/internal/cli/files"
	"github.com/lodestar-cli/lodestar/internal/cli/home"
	"github.com/lodestar-cli/lodestar/internal/common/auth"
	"github.com/lodestar-cli/lodestar/internal/common/repo"
)

type ShowCliOptions struct{
	Username       string
	Token          string
	App            string
	ConfigPath     string
}

type Show struct{
	CliOptions           ShowCliOptions
	GitAuth              auth.GitCredentials
	Repository           *repo.LodestarRepository
	AppConfigurationFile *files.AppConfigurationFile
	AppStateFile         *files.AppStateFile
}

func NewShow(username string, token string, app string, configPath string) (*Show, error){
	var err error

	cli := ShowCliOptions{
		Username: username,
		Token: token,
		App: app,
		ConfigPath: configPath,
	}

	s := Show{
		CliOptions: cli,
	}

	//1. get app config file
	err = s.setAppConfigurationFile()
	if err != nil {
		return nil, err
	}

	s.setAuth()

	//3. Clone Manifest Repository
	fmt.Printf("Cloning %s as %s...\n", s.AppConfigurationFile.Info.RepoUrl, s.CliOptions.Username)
	s.Repository, err = repo.NewLodestarRepository(s.AppConfigurationFile.Info.RepoUrl, s.GitAuth)

	//4. Fetch App State File from Repository
	s.AppStateFile, err = files.NewAppStateFile(s.Repository,s.AppConfigurationFile.Info.StatePath, s.AppConfigurationFile.Info.Name)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Show) Execute() {
	fmt.Println(" AppConfiguration")
	fmt.Println("------------------")
	s.AppConfigurationFile.Print()
	fmt.Println(" AppState")
	fmt.Println("----------")
	s.AppStateFile.Print()
}

func (s *Show) setAppConfigurationFile() error {
	if s.CliOptions.App == "" && s.CliOptions.ConfigPath == "" {
		return errors.New("must provide an App name or a path to a configuration file. For more information, run: lodestar app push --help")
	}else if s.CliOptions.ConfigPath != "" {
		var err error
		s.AppConfigurationFile, err = files.NewAppConfigurationFile(s.CliOptions.ConfigPath)
		if err != nil {
			return err
		}
	} else {
		path, err := home.GetPath("app", s.CliOptions.App)
		if err != nil {
			return err
		}
		s.AppConfigurationFile, err = files.NewAppConfigurationFile(path)
	}
	return nil
}

func (s *Show) setAuth() error{

	var a auth.GitCredentials
	if s.CliOptions.Username == "" {
		return errors.New("username is not set. Try lodestar app push --help")
	} else if s.CliOptions.Token == "" {
		return errors.New("token is not set. Try lodestar app push --help")
	} else {
		a = &auth.TokenCredentials{
			Username: s.CliOptions.Username,
			Token: s.CliOptions.Token,
		}
		s.GitAuth = a
	}

	return nil
	//need to add option that creates ssh credentials
}