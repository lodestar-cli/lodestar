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
	configPath := *flag.String("configPath", "", "path to config file")
	tag := *flag.String("tag", "", "new image tag")

	auth, err := auth.CreateAuth(username, token)
	if err != nil {
		log.Fatal(err)
	}

	repository, err := repo.GetRepository(url, auth)
	if err != nil {
		log.Fatal(err)
	}

	values, err := repo.GetFileContent(configPath)
	if err != nil {
		log.Fatal(err)
	}

	oldTag, err := yaml.GetTag(values)
	if err != nil {
		log.Fatal(err)
	}

	newConfig, err := yaml.ReplaceTag(values, tag)
	if err != nil {
		log.Fatal(err)
	}

	output, err := repo.UpdateAndPush(repository, configPath, newConfig, auth, oldTag, tag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

}
