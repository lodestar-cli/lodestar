package app

import (
	"fmt"
	"github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
)

func List() error {
	fileNames, err := lodestarDir.GetConfigFileNames("app")
	if err != nil {
		return err
	}
	fmt.Println(" -Name-\t\t-Description-")
	for _, name := range fileNames {
		path, err := lodestarDir.GetConfigPath("app", name)
		a, err := GetAppConfig(path)
		if err != nil {
			return err
		}
		fmt.Printf("* %s\t\t%s\n",name,a.AppInfo.Description)
	}
	return nil
}