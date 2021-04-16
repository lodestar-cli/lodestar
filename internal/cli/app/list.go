package app

import (
	"fmt"
	"github.com/lodestar-cli/lodestar/internal/cli/files"
	"github.com/lodestar-cli/lodestar/internal/cli/home"
)

type List struct{
	AppFilePaths          []string
	AppConfigurationFiles []*files.AppConfigurationFile
}

func NewList() (*List, error){
	var err error
	l := List{}
	l.AppFilePaths, err = home.GetConfigFilePaths("app")
	if err != nil{
		return nil, err
	}
	for _, p := range l.AppFilePaths{
		f, err := files.NewAppConfigurationFile(p)
		if err != nil{
			return nil, err
		}
		l.AppConfigurationFiles = append(l.AppConfigurationFiles, f)
	}

	return &l, nil

}

func (l *List) Execute() {
	fmt.Println(" -Name-\t\t-Description-")
	for _, app := range l.AppConfigurationFiles {
		fmt.Printf("* %s\t\t%s\n",app.Name,app.Info.Description)
	}
}