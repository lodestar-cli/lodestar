package app

import (
	"fmt"
	"github.com/lodestar-cli/lodestar/internal/cli/file"
	"github.com/lodestar-cli/lodestar/internal/cli/home"
)

type List struct {
	AppFilePaths          []string
	AppConfigurationFiles []*file.AppConfigurationFile
}

func NewList() (*List, error) {
	var err error
	l := List{}
	l.AppFilePaths, err = home.GetConfigFilePaths("app")
	if err != nil {
		return nil, err
	}
	for _, p := range l.AppFilePaths {
		f, err := file.NewAppConfigurationFile(p)
		if err != nil {
			return nil, err
		}
		l.AppConfigurationFiles = append(l.AppConfigurationFiles, f)
	}

	return &l, nil

}

func (l *List) Execute() {
	fmt.Println(" -App-")
	for _, app := range l.AppConfigurationFiles {
		fmt.Printf("* %s\n", app.Name)
	}
}
