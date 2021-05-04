package remote

import (
	"bytes"
	"github.com/goccy/go-yaml"
	"github.com/lodestar-cli/lodestar/internal/common/environment"
	"testing"
	"time"
)

func TestAppStateFile_GetByteContent(t *testing.T) {
	a, err := setTestAppStateFile()
	if err != nil {
		t.Errorf("Error settint test app state file: %s", err)
	}
	b := a.GetByteContent()

	if !bytes.Equal(a.ByteContent, b) {
		t.Error("Error Updating File: byte content does not update")
	}
}

func TestAppStateFile_GetStringContent(t *testing.T) {
	a, err := setTestAppStateFile()
	if err != nil {
		t.Errorf("Error settint test app state file: %s", err)
	}
	s := a.GetStringContent()

	if a.StringContent != s {
		t.Error("Error Updating File: string content does not update")
	}
}

func TestAppStateFile_UpdateEnvironmentStateGraph(t *testing.T) {
	a, err := setTestAppStateFile()
	if err != nil {
		t.Errorf("Error settint test app state file: %s", err)
	}
	m := map[string]string{
		"tag":  "newTag",
		"tag1": "alsoNewTag",
	}

	testTable := []struct {
		EnvName       string
		EnvMap        map[string]string
		ExpectUpdated bool
		ExpectError   bool
	}{
		{"dev", m, true, false},
		{"badEnv", m, false, true},
		{"dev", a.EnvironmentStateGraph[0].YamlKeys, false, false},
	}

	for _, test := range testTable {
		a, err := setTestAppStateFile()
		if err != nil {
			t.Errorf("Error settint test app state file: %s", err)
		}
		var n int

		for i, e := range a.EnvironmentStateGraph {
			if e.Name == test.EnvName {
				n = i
			}
		}

		u, err := a.UpdateEnvironmentStateGraph(test.EnvName, test.EnvMap)
		if err != nil {
			if test.ExpectError {
			} else {
				t.Errorf("Error updating EnvironmentStateGraph: %s", err)
			}
		}

		if u != test.ExpectUpdated {
			t.Errorf("Bad Update: Expected %v and got %v", test.ExpectUpdated, u)
		}

		if test.ExpectUpdated {
			for k, _ := range test.EnvMap {
				if a.EnvironmentStateGraph[n].YamlKeys[k] != test.EnvMap[k] {
					t.Errorf("Error after updating EnvironmentStateGraph: for key %s: expected %s but got %s", k, test.EnvMap[k], a.EnvironmentStateGraph[n].YamlKeys[k])
				}
			}
		}
	}
}

func TestAppStateFile_UpdateFile(t *testing.T) {
	a, err := setTestAppStateFile()
	if err != nil {
		t.Errorf("Error settint test app state file: %s", err)
	}
	// a content
	b := a.ByteContent
	s := a.StringContent
	a.ByteContent = []byte("")
	a.StringContent = ""

	a.UpdateFile()

	if !bytes.Equal(a.ByteContent, b) {
		t.Error("Error Updating File: byte content does not update")
	}
	if a.StringContent != s {
		t.Error("Error Updating File: string content does not update")
	}
	// a time

	now := time.Now()
	time.Sleep(time.Second)
	a.UpdateFile()

	newTime, err := time.Parse(time.RFC3339, a.Updated)
	if err != nil {
		t.Errorf("error setting time: %s", err)
	}

	if !now.Before(newTime) {
		t.Error("Error Updating File: time is not updated")
	}

}

func setTestAppStateFile() (*AppStateFile, error) {
	testGraph := AppStateGraph{
		Updated: time.Now().Format(time.RFC3339),
		EnvironmentStateGraph: []environment.EnvironmentState{
			{
				Name:     "dev",
				YamlKeys: map[string]string{"tag": "123456", "tag1": "4567899"},
			},
			{
				Name:     "qa",
				YamlKeys: map[string]string{"tag": "123456", "tag1": "4567899"},
			},
			{
				Name:     "staging",
				YamlKeys: map[string]string{"tag": "123456", "tag1": "4567899"},
			},
			{
				Name:     "prod",
				YamlKeys: map[string]string{"tag": "123456", "tag1": "4567899"},
			},
		},
	}
	bytes, err := yaml.Marshal(testGraph)
	if err != nil {
		return nil, err
	}

	a := AppStateFile{
		Name:                  "test",
		Path:                  "/test/test.yaml",
		Updated:               testGraph.Updated,
		EnvironmentStateGraph: testGraph.EnvironmentStateGraph,
		ByteContent:           bytes,
		StringContent:         string(bytes),
	}

	return &a, nil
}
