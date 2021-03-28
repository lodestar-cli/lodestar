package app

import (
	"errors"
	"fmt"
	"github.com/go-git/go-billy/v5"
	auth "github.com/lodestar-cli/lodestar/internal/common/auth"
	"github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
	repo "github.com/lodestar-cli/lodestar/internal/common/repo"
	tag "github.com/lodestar-cli/lodestar/internal/common/tag"
)

func Push(username string, to string, name string, configPath string, environment string, t string, output bool) error {
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
				err := push(username, to, config.AppInfo.RepoUrl, environment, env.SrcPath, config.AppInfo.StatePath, t, output)
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
				err := push(username, to, config.AppInfo.RepoUrl, environment ,env.SrcPath, config.AppInfo.StatePath, t, output)
				if err != nil {
					return err
				}
				break
			}
		}
		return nil
	}
}

func push(username string, token string, url string, environment string, configPath string, statePath string, t string, output bool) error {
	var fs billy.Filesystem

	auth, err := auth.CreateAuth(username, token)
	if err != nil {
		return err
	}

	fmt.Printf("Cloning %s as %s...\n", url, username)
	repository, fs, err := repo.GetRepository(url, auth)
	if err != nil {
		return err
	}

	stateGraph, err := GetEnvironmentStateGraph(fs, statePath)
	if err != nil {
		return err
	}

	m, err := CompareEnvironmentStateTag(stateGraph, environment, t)
	if err != nil {
		return err
	}
	if !m{
		return nil
	}

	fmt.Printf("Retrieving environment configuration at %s...\n", configPath)
	values, err := repo.GetFileString(fs, configPath)
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

	repository, stateGraph, err = UpdateEnvironmentStateTag(fs, repository, stateGraph, statePath, environment, t)
	if err != nil {
		return err
	}

	fmt.Printf("Pushing changes to %s as %s...\n",url,username)

	err = repo.UpdateAndPush(fs, repository, configPath, newConfig, auth, oldTag, t)
	if err != nil {
		return err
	}

	fmt.Println("Push complete!")

	if output {
		err = OutputEnvironmentStateGraph(stateGraph)
		if err != nil {
			return err
		}
	}

	return nil
}
