package lodestarDir

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path"
	"strings"
)

func GetConfigPath(command string, commandName string) (string, error)  {
	dir, err := getUserHomeDir()
	if err != nil {
		return "", err
	}
	lodestarDirectory := path.Join(dir, ".lodestar/",command, commandName+".yaml")

	return lodestarDirectory, nil
}

func GetConfigFileNames(command string) ([]string, error){
	var fileNames []string

	dir, err := getUserHomeDir()
	files, err := ioutil.ReadDir(path.Join(dir, ".lodestar/",command))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if strings.Contains(f.Name(), ".yaml"){
			name := strings.ReplaceAll(f.Name(),".yaml","")
			fileNames = append(fileNames, name)
		}
	}

	return fileNames, nil
}

func AddConfig(content []byte, command string, commandName string) error {
	path, err := GetConfigPath(command, commandName)

	err = ioutil.WriteFile(path,content,644)
	if err != nil {
		return err
	}

	fmt.Println("Added config file at: "+path)
	return nil
}

func GetConfigContent(path string) ([]byte, error){

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func getUserHomeDir() (string, error){
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}