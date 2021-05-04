package home

import (
	"os/user"
	"path"
	"strings"
)

func GetPath(command string, commandName string) (string, error) {
	dir, err := getUserHomeDir()
	if err != nil {
		return "", err
	}
	lodestarDirectory := path.Join(dir, ".lodestar/", command, commandName+".yaml")

	return lodestarDirectory, nil
}

func GetConfigFilePaths(command string, iou IoUtil) ([]string, error) {
	var filePaths []string
	h, err := getUserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := path.Join(h, ".lodestar/", command)

	files, err := iou.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if strings.Contains(f.Name(), ".yaml") {
			fp := path.Join(dir, f.Name())
			filePaths = append(filePaths, fp)
		}
	}

	return filePaths, nil
}

func GetContent(path string, iou IoUtil) ([]byte, error) {

	content, err := iou.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func getUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
