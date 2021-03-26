package app

import (
	"fmt"
	auth "github.com/lodestar-cli/lodestar/internal/common/auth"
	repo "github.com/lodestar-cli/lodestar/internal/common/repo"
	yaml "github.com/lodestar-cli/lodestar/internal/common/yaml"
)

func Promote(username string, token string, url string, srcPath string, destPath string) error {

	auth, err := auth.CreateAuth(username, token)
	if err != nil {
		return err
	}

	fmt.Printf("Cloning %s as %s...\n", url, username)
	repository, err := repo.GetRepository(url, auth)
	if err != nil {
		return err
	}

	fmt.Printf("Retrieving tag from configuration file %s...\n",srcPath)
	src, err := repo.GetFileContent(srcPath)
	if err != nil {
		return err
	}
	newTag, err := yaml.GetTag(src)
	if err != nil {
		return err
	}
	dest, err := repo.GetFileContent(destPath)
	if err != nil {
		return err
	}
	oldTag, err := yaml.GetTag(dest)
	if err != nil {
		return err
	}

	fmt.Printf("Updating %s to %s at %s...\n",oldTag,newTag, destPath)
	newConfig, err := yaml.ReplaceTag(dest, newTag)
	if err != nil {
		return err
	}

	fmt.Printf("Pushing changes to %s as %s...\n",url,username)
	err = repo.UpdateAndPush(repository, destPath, newConfig, auth, oldTag, newTag)
	if err != nil {
		return err
	}

	fmt.Println("Promote complete!")
	return nil

}
