package home
import (
	"io/ioutil"
	"os/user"
	"path"
	"strings"
)

func GetPath(command string, commandName string) (string, error)  {
	dir, err := getUserHomeDir()
	if err != nil {
		return "", err
	}
	lodestarDirectory := path.Join(dir, ".lodestar/",command, commandName+".yaml")

	return lodestarDirectory, nil
}

func GetConfigFilePaths(command string) ([]string, error){
	var filePaths []string

	dir, err := getUserHomeDir()
	files, err := ioutil.ReadDir(path.Join(dir, ".lodestar/",command))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if strings.Contains(f.Name(), ".yaml"){
			fp := path.Join(dir, f.Name())
			filePaths = append(filePaths, fp)
		}
	}

	return filePaths, nil
}

func GetContent(path string) ([]byte, error){

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