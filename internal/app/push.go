package app

import (
	"errors"
	"fmt"
	auth "github.com/lodestar-cli/lodestar/internal/common/auth"
	"github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
	repo "github.com/lodestar-cli/lodestar/internal/common/repo"
	 tag "github.com/lodestar-cli/lodestar/internal/common/tag"

)

func Push(username string, to string, name string, configPath string, environment string, t string) error {
	var config *LodestarAppConfig

	if name == "" && configPath == "" {
		return errors.New("Must provide an App name or a path to a configuration file. For more information, run: lodestar app push --help")
	} else if configPath != "" {
		config, err := GetAppConfig(configPath)
		if err != nil {
			return err
		}
		if len(config.EnvGraph) < 1 {
			return errors.New("No environments are provided for " + config.AppInfo.Name)
		}

		for _, env := range config.EnvGraph {
			if env.Name == environment {
				err := push(username, to, config.AppInfo.RepoUrl, env.SrcPath, t)
				if err != nil {
					return err
				}
				break
			}
		}
		return nil
	} else {
		path, err := lodestarDir.GetConfigPath("app", name)
		if err != nil {
			return err
		}
		fmt.Printf("Retrieving config for %s...\n", name)
		config, err = GetAppConfig(path)
		if err != nil {
			return err
		}
		if len(config.EnvGraph) < 1 {
			return errors.New("No environments are provided for " + name)
		}

		for _, env := range config.EnvGraph {
			if env.Name == environment {
				err := push(username, to, config.AppInfo.RepoUrl, env.SrcPath, t)
				if err != nil {
					return err
				}
				break
			}
		}
		return nil
	}
}

func push(username string, token string, url string, configPath string, t string) error {

	auth, err := auth.CreateAuth(username, token)
	if err != nil {
		return err
	}
	fmt.Printf("Cloning %s as %s...\n", url, username)
	repository, err := repo.GetRepository(url, auth)
	if err != nil {
		return err
	}

	fmt.Printf("Retrieving environment configuration at %s...\n", configPath)
	values, err := repo.GetFileContent(configPath)
	if err != nil {
		return err
	}

	oldTag, err := tag.Get(values)
	if err != nil {
		return err
	}

	fmt.Printf("Updating %s to %s...\n",oldTag,t)
	newConfig, err := tag.Replace(values, t)
	if err != nil {
		return err
	}

	fmt.Printf("Pushing changes to %s as %s...\n",url,username)
	err = repo.UpdateAndPush(repository, configPath, newConfig, auth, oldTag, t)
	if err != nil {
		return err
	}

	fmt.Println("Push complete!")
	return nil

}
