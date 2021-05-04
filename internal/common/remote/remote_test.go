package remote

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/goccy/go-yaml"
	"github.com/lodestar-cli/lodestar/internal/common/auth"
	sops "go.mozilla.org/sops/v3/decrypt"
	"os"
	"testing"
)

type TestUser struct {
	Username string `yaml:"gitUser"`
	Token    string `yaml:"gitToken"`
}

var (
	testUser           = TestUser{}
	testRepository     = LodestarRepository{}
	testManagementFile = ManagementFile{}
)

func TestMain(m *testing.M) {
	err := setTestUser()
	if err != nil {
		fmt.Printf("error setting test user: %s", err)
	}

	err = setTestLodestarRepository()
	if err != nil {
		fmt.Printf("error setting test repository: %s", err)
	}

	setTestManagementFile()

	os.Exit(m.Run())
}

func setTestLodestarRepository() error {
	var err error

	testRepository.Url = "https://github.com/lodestar-cli/lodestar-folder-app-example.git"
	testRepository.Credentials = &auth.TokenCredentials{
		Username: testUser.Username,
		Token:    testUser.Token,
	}
	testRepository.FileSystem = memfs.New()
	testRepository.Storage = memory.NewStorage()

	a := http.BasicAuth{
		Username: testUser.Username,
		Password: testUser.Token,
	}

	cloneOptions := &git.CloneOptions{
		URL:  testRepository.Url,
		Auth: &a,
	}

	testRepository.Repository, err = git.CloneContext(context.TODO(), testRepository.Storage, testRepository.FileSystem, cloneOptions)
	if err != nil {
		return err
	}

	testRepository.Worktree, err = testRepository.Repository.Worktree()
	if err != nil {
		return err
	}

	return nil
}

func setTestUser() error {
	path := "../../../sops-secrets/test-credentials.enc.yaml"
	content, err := sops.File(path, "yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &testUser)
	if err != nil {
		return err
	}

	if testUser.Username == "" || testUser.Token == "" {
		return errors.New("couldn't load sops user from local repository")
	}

	return nil
}

func setTestManagementFile() {
	testYaml := `
test: "fail"
  test1: fail1
test2: fail2
`
	testManagementFile.Name = "test"
	testManagementFile.Path = "/test/test.yaml"
	testManagementFile.ByteContent = []byte(testYaml)
	testManagementFile.StringContent = testYaml

}
