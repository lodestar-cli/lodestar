package lodestarDir

import (
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