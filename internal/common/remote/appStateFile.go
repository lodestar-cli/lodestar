package remote

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/lodestar-cli/lodestar/internal/common/environment"
	"io/ioutil"
	"time"
)

type AppStateGraph struct {
	Updated               string             `yaml:"updated"`
	EnvironmentStateGraph []environment.EnvironmentState `yaml:"environmentStateGraph"`
}

//LodestarFile
type AppStateFile struct {
	Path             string
	Name             string
	Updated          string
	EnvironmentGraph []environment.EnvironmentState
	ByteContent      []byte
	StringContent    string
}

func NewAppStateFile(repository *LodestarRepository, path string, name string) (*AppStateFile, error){
	stat, err := repository.FileSystem.Stat(path)
	if err != nil{
		return nil, err
	}
	bytes := make([]byte, stat.Size())

	//get file from memory
	file, err := repository.FileSystem.Open(path)
	if err != nil{
		return nil, err
	}
	//get file content as string
	_, err = file.Read(bytes)
	if err != nil{
		return nil, err
	}

	a := AppStateFile{
		Name: name+"-state",
		Path: path,
		ByteContent: bytes,
		StringContent: string(bytes),
	}

	err = a.setAppStateGraphFields()
	if err != nil{
		return nil, err
	}

	return &a, nil
}

//LodestarFile Interface Functions

func (a *AppStateFile) Print(){
	fmt.Println(a.StringContent)
}

func (a *AppStateFile) Output() error {
	err := ioutil.WriteFile(a.Name+".yaml", a.ByteContent, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (a *AppStateFile)GetStringContent() string {
	return a.StringContent
}

func (a *AppStateFile)GetByteContent() []byte {
	return a.ByteContent
}

//AppState Specific Functions

func (a *AppStateFile) UpdateEnvironmentGraph(env string, keys map[string]string) (bool, error) {
	for _, e := range a.EnvironmentGraph{
		if e.Name == env {
			update := e.UpdateKeys(keys)
			if update{
				return true, nil
			} else {
				return false, nil
			}
		}
	}
	return false, fmt.Errorf("%s does not exist in this app", env)
}

func (a *AppStateFile)UpdateFile() error {

	a.Updated = time.Now().Format(time.RFC3339)

	s := AppStateGraph{
		Updated: a.Updated,
		EnvironmentStateGraph: a.EnvironmentGraph,
	}

	bytes ,err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	a.ByteContent = bytes
	a.StringContent = string(bytes)

	return nil
}

func (a *AppStateFile) setAppStateGraphFields() error {
	stateGraph := AppStateGraph{}
	err := yaml.Unmarshal(a.ByteContent, &stateGraph)
	if err != nil {
		return err
	}

	a.EnvironmentGraph = stateGraph.EnvironmentStateGraph
	a.Updated = stateGraph.Updated

	return nil
}