package app
import (
	"github.com/goccy/go-yaml"
	lodestarDir "github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
)

type appEnvironment struct {
	Name string `tag:"name"`
	SrcPath string `tag:"srcPath"`
}

type LodestarAppConfig struct {
	AppInfo struct {
		Name string   `tag:"name"`
		Type string `tag:"type"`
		Description string `tag:"description"`
		RepoUrl string `tag:"repoUrl"`
		Target string `tag:"target"`
		StatePath string `tag:"statePath"`
	} `tag:"appInfo"`
	EnvGraph []appEnvironment `tag:"envGraph,flow"`
}


//retrieves object from config tag
func GetAppConfig(path string) (*LodestarAppConfig, error){
	//read from app configuration tag
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