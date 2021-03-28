package app

import (
	"errors"
	"fmt"
	"github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
)

func Show(appName string) error {
	if appName == "" {
		return errors.New("Must provide an App name. \n For more information, run: lodestar app push --help")
	}

	path, err := lodestarDir.GetConfigPath("app", appName)
	content, err := lodestarDir.GetConfigContent(path)
	if err != nil {
		return err
	}

	text := string(content)
	fmt.Println(text)

	return nil
}