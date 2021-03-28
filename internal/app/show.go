package app

import (
	"fmt"
	"github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
)

func Show(appName string) error {
	path, err := lodestarDir.GetConfigPath("app", appName)
	content, err := lodestarDir.GetConfigContent(path)
	if err != nil {
		return err
	}

	text := string(content)
	fmt.Println(text)

	return nil
}