package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	auth "github.com/lodestar-cli/lodestar/internal/common/auth"
	repo "github.com/lodestar-cli/lodestar/internal/common/repo"
	yaml "github.com/lodestar-cli/lodestar/internal/common/yaml"
)

func main() {
	username := *flag.String("username", os.Getenv("GIT_USER"), "the username for authenticaiton")
	token := *flag.String("token", os.Getenv("GIT_TOKEN"), "the token for authenticaiton")
	url := *flag.String("repoUrl", "", "repository to pull")
	srcPath := *flag.String("srcPath", "", "path to source config file")
	destPath := *flag.String("destPath", "", "path to destination config file")

	auth, err := auth.CreateAuth(username, token)
	if err != nil {
		log.Fatal(err)
	}

	repository, err := repo.GetRepository(url, auth)
	if err != nil {
		log.Fatal(err)
	}

	src, err := repo.GetFileContent(srcPath)
	if err != nil {
		log.Fatal(err)
	}

	newTag, err := yaml.GetTag(src)
	if err != nil {
		log.Fatal(err)
	}

	dest, err := repo.GetFileContent(destPath)
	if err != nil {
		log.Fatal(err)
	}

	oldTag, err := yaml.GetTag(dest)
	if err != nil {
		log.Fatal(err)
	}

	newConfig, err := yaml.ReplaceTag(dest, newTag)
	if err != nil {
		log.Fatal(err)
	}

	output, err := repo.UpdateAndPush(repository, destPath, newConfig, auth, oldTag, newTag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

}
