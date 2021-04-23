package file

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/lodestar-cli/lodestar/internal/cli/home"
	"github.com/lodestar-cli/lodestar/internal/common/environment"
	"io/ioutil"
)

type Info struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	RepoUrl     string `yaml:"repoUrl"`
	Target      string `yaml:"target"`
	StatePath   string `yaml:"statePath"`
}

type AppConfiguration struct {
	Info     Info                      `yaml:"info"`
	EnvGraph []environment.Environment `yaml:"environmentGraph"`
	YamlKeys []string                  `yaml:"yamlKeys"`
}

//LodestarFile
type AppConfigurationFile struct {
	Path             string
	Name             string
	Info             Info
	EnvironmentGraph []environment.Environment
	YamlKeys         []string
	ByteContent      []byte
	StringContent    string
}

func NewAppConfigurationFile(path string) (*AppConfigurationFile, error) {
	a := AppConfiguration{}
	content, err := home.GetContent(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &a)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	f := AppConfigurationFile{
		Path:             path,
		Name:             a.Info.Name,
		Info:             a.Info,
		EnvironmentGraph: a.EnvGraph,
		YamlKeys:         a.YamlKeys,
		ByteContent:      content,
		StringContent:    string(content),
	}

	return &f, nil
}

func (a *AppConfigurationFile) Output() error {
	err := ioutil.WriteFile(a.Name+".yaml", a.ByteContent, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (a *AppConfigurationFile) Print() {
	fmt.Println(a.StringContent)
}

func (a *AppConfigurationFile) GetStringContent() string {
	return a.StringContent
}

func (a *AppConfigurationFile) GetByteContent() []byte {
	return a.ByteContent
}

func (a *AppConfigurationFile) GetEnvironment(environment string) (*environment.Environment, error) {
	if len(a.EnvironmentGraph) < 1 {
		return nil, fmt.Errorf("No environments are provided for " + a.Info.Name)
	}

	for _, e := range a.EnvironmentGraph {
		if e.Name == environment {
			if e.SrcPath != "" {
				return &e, nil
			}
			return nil, fmt.Errorf("configuration file path for environment %s is empty in app %s", environment, a.Info.Name)
		}
	}

	return nil, fmt.Errorf("failed to find environment %s in app %s", environment, a.Info.Name)
}

func (a *AppConfiguration) UpdateAppConfiguration() {

}
