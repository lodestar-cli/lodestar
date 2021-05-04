package remote

import (
	"github.com/goccy/go-yaml"
	"github.com/lodestar-cli/lodestar/internal/common/auth"
	"github.com/lodestar-cli/lodestar/internal/common/environment"
	"testing"
	"time"
)

func TestNewLodestarRepository(t *testing.T) {
	var cPass auth.GitCredentials
	var cFail auth.GitCredentials

	cPass = &auth.TokenCredentials{
		Username: testUser.Username,
		Token:    testUser.Token,
	}

	cFail = &auth.TokenCredentials{
		Username: "badUser",
		Token:    "badToken",
	}

	testTable := []struct {
		Url         string
		Credentials auth.GitCredentials
		ExpectError bool
	}{
		{"", cPass, true},
		{"https://github.com/lodestar-cli/lodestar-folder-app-example.git", cPass, false},
		{"github.com/lodestar-cli/lodestar-folder-app-example.git", cPass, true},
		{"https://github.com/lodestar-cli/lodestar-folder-app-example.git", cFail, true},
	}

	for _, test := range testTable {
		r, err := NewLodestarRepository(test.Url, test.Credentials)
		if err != nil {
			if test.ExpectError {
				continue
			} else {
				t.Errorf("error creating lodestar repository: %s", err)
			}
		}

		if r.Repository == nil {
			t.Error("NewLodestarRepository did not clone a repository")
		}

		if r.Worktree == nil {
			t.Error("NewLodestarRepository did not create a new Worktree")
		}
	}
}

/*func TestLodestarRepository_CommitFiles(t *testing.T) {
	aPass, err := setTestAppStateFileForCommit("lodestar-folder-app-example.yaml")
	if err != nil{
		t.Errorf("error creating test app state file: %s", err)
	}



	aFail, err := setTestAppStateFileForCommit("bad-path.yaml")
	if err != nil{
		t.Errorf("error creating test app state file: %s", err)
	}


	testTable := []struct {
		Files       []LodestarFile
		ExpectError bool
	}{
		{[]LodestarFile{aPass,aFail}},
	}

	for _, test := range testTable{
		r := testRepository
		err := r.CommitFiles("Test Commit Go Testing", test.Files...)
		if err != nil{
			if test.ExpectError{
				continue
			}

			t.Errorf("error commiting files: %s", err)
		}
	}
}*/

func setTestAppStateFileForCommit(path string) (*AppStateFile, error) {
	testGraph := AppStateGraph{
		Updated: time.Now().Format(time.RFC3339),
		EnvironmentStateGraph: []environment.EnvironmentState{
			{
				Name:     "dev",
				YamlKeys: map[string]string{"tag": "123456"},
			},
			{
				Name:     "qa",
				YamlKeys: map[string]string{"tag": "123456"},
			},
			{
				Name:     "staging",
				YamlKeys: map[string]string{"tag": "123456"},
			},
			{
				Name:     "prod",
				YamlKeys: map[string]string{"tag": "123456"},
			},
		},
	}
	bytes, err := yaml.Marshal(testGraph)
	if err != nil {
		return nil, err
	}

	a := AppStateFile{
		Name:                  "test",
		Path:                  path,
		Updated:               testGraph.Updated,
		EnvironmentStateGraph: testGraph.EnvironmentStateGraph,
		ByteContent:           bytes,
		StringContent:         string(bytes),
	}

	return &a, nil
}

func setManagementFileForCommit(path string) {
	m := ManagementFile{
		Name: "dev-config",
		Path: path,
	}

	m.StringContent = `
tag: "testCommit"
`
	m.ByteContent = []byte(m.StringContent)
}
