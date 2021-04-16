package files

import (
	"errors"
	"fmt"
	"github.com/lodestar-cli/lodestar/internal/common/environment"
	"github.com/lodestar-cli/lodestar/internal/common/repo"
	"io/ioutil"
	"strings"
)

type ManagementFile struct{
	Name          string
	Path          string
	ByteContent   []byte
	StringContent string
}

func NewManagementFile(env *environment.Environment, repository *repo.LodestarRepository) (*ManagementFile, error){
	//Create array the size of the file content
	stat, err := repository.FileSystem.Stat(env.SrcPath)
	if err != nil{
		return nil, err
	}
	bytes := make([]byte, stat.Size())

	//get file from memory
	file, err := repository.FileSystem.Open(env.SrcPath)
	if err != nil{
		return nil, err
	}

	//get file content as string
	file.Read(bytes)

	c := ManagementFile{
		Name: env.Name+"-config",
		Path: env.SrcPath,
		ByteContent: bytes,
		StringContent: string(bytes),
	}

	return &c, nil
}

func (m *ManagementFile) Print(){
	fmt.Println(m.StringContent)
}

func (m *ManagementFile) Output() error {
	err := ioutil.WriteFile(m.Name+".yaml", m.ByteContent, 0755)
	if err != nil {
		return err
	}

	return err
}

func (m *ManagementFile) GetStringContent() string {
	return m.StringContent
}

func (m *ManagementFile) GetByteContent() []byte {
	return m.ByteContent
}

func (m *ManagementFile) UpdateFileContents(keysMap map[string]string) (bool,error) {
	lines := strings.Split(m.StringContent, "\n")
	var usedKeys []string
	updated := false

	for j, line := range lines {
		for k, v := range keysMap {
			if strings.Contains(line, k) {
				usedKeys = append(usedKeys, k)
				txt := strings.Split(line, " ")
				for i, w := range txt {
					if strings.Contains(w, k) {
						cv := strings.Join(txt[i+1:], " ")
						if cv != v {
							updated = true
							nl := txt[:i+1]
							nl = append(nl, v)
							t := strings.Join(nl, " ")
							lines[j] = t
						}
					}
				}
			}
		}
		if len(usedKeys) == len(keysMap) {
			break
		}
	}

	if updated{
		m.StringContent = strings.Join(lines, "\n")
		m.ByteContent = []byte(m.StringContent)
		return updated, nil
	}

	return updated, nil
}

func (m *ManagementFile) GetKeyValues(keys []string) (map[string]string, error){
	lines := strings.Split(m.StringContent, "\n")
	keyMap := map[string]string{}

	for _, line := range lines {
		for _, k := range keys {
			if strings.Contains(line, k) {
				txt := strings.Split(line, " ")
				for i, w := range txt {
					if strings.Contains(w, k) {
						v := strings.Join(txt[i+1:], " ")
						keyMap[k] = v
					}
				}
			}
		}
	}

	if len(keyMap) != len(keys){
		return nil, errors.New("couldn't find all keys listed in AppConfiguration file in Destination Environment - Cannot do promote")
	}

	return keyMap, nil
}
