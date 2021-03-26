package app

import (
	"fmt"
	auth "github.com/lodestar-cli/lodestar/internal/common/auth"
	repo "github.com/lodestar-cli/lodestar/internal/common/repo"
	yaml "github.com/lodestar-cli/lodestar/internal/common/yaml"

)

func Push(username string, token string, url string, configPath string, tag string) error {

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

	oldTag, err := yaml.GetTag(values)
	if err != nil {
		return err
	}

	fmt.Printf("Updating %s to %s...\n",oldTag,tag)
	newConfig, err := yaml.ReplaceTag(values, tag)
	if err != nil {
		return err
	}

	fmt.Printf("Pushing changes to %s as %s...\n",url,username)
	err = repo.UpdateAndPush(repository, configPath, newConfig, auth, oldTag, tag)
	if err != nil {
		return err
	}

	fmt.Println("Push complete!")
	return nil

}
