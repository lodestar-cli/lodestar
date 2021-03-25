package repo
import (
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

var (
	storer = memory.NewStorage()
	fs     = memfs.New()
)

func GetRepository(repository string, auth *http.BasicAuth) (*git.Repository, error) {

	clone, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:  repository,
		Auth: auth,
	})
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	fmt.Println("Repository cloned")
	return clone, nil

}

func GetFileContent(path string) (string, error){
	//Create array the size of the file content
	stat, err := fs.Stat(path)
	if err != nil{
		return "", err
	}
	bytes := make([]byte, stat.Size()+1)

	//get file from memory
	file, err := fs.Open(path)
	if err != nil{
		return "", err
	}

	//get file content as string
	file.Read(bytes)
	content := string(bytes)

	return content, nil
}

func UpdateAndPush(repository *git.Repository, configPath string, newConfig string, auth *http.BasicAuth, oldTag string, newTag string) (string, error){
	worktree, err := repository.Worktree()
	if err != nil {
		fmt.Printf("%v", err)
		return "", err
	}

	worktree.Remove(configPath)

	configFile, err := fs.Create(configPath)
	if err != nil {
		return "", err
	}
	configFile.Write([]byte(newConfig))
	worktree.Add(configPath)

	worktree.Commit("Config file updated with Bazel: "+oldTag+" ---> "+newTag, &git.CommitOptions{})


	err = repository.Push(&git.PushOptions{
		Auth:       auth,
	})
	if err != nil {
		return "" , err
	}

	return "Commit Pushed!", nil

}
