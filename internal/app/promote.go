package app

import (
	"errors"
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/goccy/go-yaml"
	auth "github.com/lodestar-cli/lodestar/internal/common/auth"
	"github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
	repo "github.com/lodestar-cli/lodestar/internal/common/repo"
	tag "github.com/lodestar-cli/lodestar/internal/common/tag"
)

func Promote(username string, token string, name string, configPath string, srcEnv string, destEnv string, output bool) error {
	var config *LodestarAppConfig
	var srcPath string
	var destPath string
	if name == "" && configPath == "" {
		return errors.New("Must provide an App name or a path to a configuration file. For more information, run: lodestar app push --help")
	} else if configPath != ""{
		config, err := GetAppConfig(configPath)
		if err != nil {
			return err
		}
		if len(config.EnvGraph) < 1 {
			return errors.New("No environments are provided for "+config.AppInfo.Name)
		}

		for _, env := range config.EnvGraph {
			if env.Name == srcEnv {
				srcPath=env.SrcPath
			}else if env.Name == destEnv {
				destPath=env.SrcPath
			}
			if srcPath != "" && destPath != "" {
				break
			}
		}
		err = promote(username,token,config.AppInfo.RepoUrl,srcPath,destPath, destEnv, config.AppInfo.StatePath, config.AppInfo.Name, output)
		return err
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
			if env.Name == srcEnv {
				srcPath=env.SrcPath
			}else if env.Name == destEnv {
				destPath=env.SrcPath
			}
			if srcPath != "" && destPath != "" {
				break
			}
		}
		err = promote(username,token,config.AppInfo.RepoUrl,srcPath,destPath, destEnv, config.AppInfo.StatePath, config.AppInfo.Name, output)
		return err
	}
}


func promote(username string, token string, url string, srcPath string, destPath string, destEnv string, statePath string, name string, output bool) error {
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

	fmt.Printf("Retrieving tag from configuration file %s...\n",srcPath)
	src, err := repo.GetFileString(fs, srcPath)
	if err != nil {
		return err
	}
	newTag, err := tag.Get(src)
	if err != nil {
		return err
	}
	dest, err := repo.GetFileString(fs, destPath)
	if err != nil {
		return err
	}
	oldTag, err := tag.Get(dest)
	if err != nil {
		return err
	}

	m, err := CompareEnvironmentStateTag(stateGraph, destEnv, newTag)
	if err != nil {
		return err
	}
	if !m{
		if output {
			err = OutputEnvironmentStateGraph(stateGraph, name)
			if err != nil {
				return err
			}
		}

		s, err := yaml.Marshal(stateGraph)
		if err != nil {
			return err
		}
		fmt.Println(string(s))

		return nil
	}

	fmt.Printf("Updating %s to %s at %s...\n",oldTag,newTag, destPath)
	newConfig, err := tag.Replace(dest, newTag)
	if err != nil {
		return err
	}

	repository, stateGraph, err = UpdateEnvironmentStateTag(fs, repository, stateGraph, statePath, destEnv, newTag)
	if err != nil {
		return err
	}

	fmt.Printf("Pushing changes to %s as %s...\n",url,username)
	err = repo.UpdateAndPush(fs,repository, destPath, newConfig, auth, oldTag, newTag)
	if err != nil {
		return err
	}

	fmt.Println("Promote complete!")

	if output {
		err = OutputEnvironmentStateGraph(stateGraph, name)
		if err != nil {
			return err
		}
	}

	s, err := yaml.Marshal(stateGraph)
	if err != nil {
		return err
	}
	fmt.Println(string(s))

	return nil

}
