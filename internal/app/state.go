package app

import (
	"errors"
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/goccy/go-yaml"
	"github.com/lodestar-cli/lodestar/internal/common/repo"
	"io/ioutil"
)

type LodestarAppStateConfig struct {
	State []struct {
		Environment string `yaml:"environment"`
		Tag string `yaml:"tag"`
	} `yaml:"state"`
}



func GetEnvironmentStateGraph(fs billy.Filesystem, path string) (*LodestarAppStateConfig, error){
	content, err := repo.GetFileByte(fs, path)
	if err != nil {
		return nil, err
	}
	stateGraph := LodestarAppStateConfig{}


	err = yaml.Unmarshal(content, &stateGraph)
	if err != nil {
		return nil, err
	}


	return &stateGraph, nil
}

func CompareEnvironmentStateTag(stateGraph *LodestarAppStateConfig, env string, newTag string) (bool, error){
	var stateTag string
	for _, state := range stateGraph.State{
		if state.Environment == env{
			stateTag = state.Tag
			break
		}
	}
	if stateTag == ""{
		return false, errors.New(env+" does not exist in this App!")
	}
	if stateTag == newTag{
		fmt.Printf("%s already is using tag: %s\n",env,newTag)
		return false, nil
	}
	return true, nil
}

func UpdateEnvironmentStateTag(fs billy.Filesystem, repository *git.Repository, stateGraph *LodestarAppStateConfig, path string, env string, tag string) (*git.Repository, *LodestarAppStateConfig,error){
	updateTag := false

	for i, state := range stateGraph.State{
		if state.Environment == env{
			stateGraph.State[i].Tag = tag
			updateTag = true
			break
		}
	}
	if updateTag != true {
		return nil, nil, errors.New("Couldn't update environment tag")
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return nil, nil, err
	}

	_, err = worktree.Remove(path)
	if err != nil {
		return nil, nil, err
	}

	configFile, err := fs.Create(path)
	if err != nil {
		return nil, nil, err
	}

	s ,err := yaml.Marshal(stateGraph)
	if err != nil {
		return nil, nil, err
	}

	_, err = configFile.Write(s)
	if err != nil {
		return nil, nil, err
	}

	_, err = worktree.Add(path)
	if err != nil {
		return nil, nil, err
	}

	return repository, stateGraph, nil
}

func OutputEnvironmentStateGraph(stateGraph *LodestarAppStateConfig) error{
	s ,err := yaml.Marshal(stateGraph)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("state.yaml", s, 0755)
	if err != nil {
		return err
	}

	return nil
}