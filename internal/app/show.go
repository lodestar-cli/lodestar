package app

import (
	"errors"
	"fmt"
	"github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
)

func Show(appName string, configPath string) error {
	var content []byte
	if appName == "" && configPath == ""{
		return errors.New("Must provide an App name or path to an App configuration file. \n For more information, run: lodestar app push --help")
	}else if configPath != ""{
		content, err := lodestarDir.GetConfigContent(configPath)
		if err != nil {
			return err
		}
		text := string(content)
		fmt.Println(text)
	} else {
		path, err := lodestarDir.GetConfigPath("app", appName)
		content, err = lodestarDir.GetConfigContent(path)
		if err != nil {
			return err
		}
		text := string(content)
		fmt.Println(text)
	}
	return nil
}