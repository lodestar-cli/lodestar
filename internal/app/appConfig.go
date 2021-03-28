package app
import (
	"github.com/goccy/go-yaml"
	lodestarDir "github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
)

type appEnvironment struct {
	Name string `yaml:"name"`
	SrcPath string `yaml:"srcPath"`
}

type LodestarAppConfig struct {
	AppInfo struct {
		Name string   `yaml:"name"`
		Type string `yaml:"type"`
		Description string `yaml:"description"`
		RepoUrl string `yaml:"repoUrl"`
		Target string `yaml:"target"`
		StatePath string `yaml:"statePath"`
	} `yaml:"appInfo"`
	EnvGraph []appEnvironment `yaml:"envGraph,flow"`
}


//retrieves object from config yaml
func GetAppConfig(path string) (*LodestarAppConfig, error){
	//read from app configuration yaml
	content, err := lodestarDir.GetConfigContent(path)
	if err != nil {
		return nil, err
	}
	app := LodestarAppConfig{}

	err = yaml.Unmarshal(content, &app)
	if err != nil {
		return nil, err
	}

	return &app, nil
}