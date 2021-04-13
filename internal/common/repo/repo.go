package repo

import (
	"bytes"
	"context"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"time"
)


func GetRepository(url string, auth *http.BasicAuth) (*git.Repository, billy.Filesystem, error) {
	var (
		storer = memory.NewStorage()
		fs     = memfs.New()
	)

	cloneOptions := new(git.CloneOptions)
	cloneOptions.URL=url
	cloneOptions.Auth=auth

	repository, err := git.CloneContext(context.TODO(),storer, fs, cloneOptions)
	if err != nil {
		return nil, nil, err
	}

	return repository, fs, nil

}

func GetFileString(fs billy.Filesystem, path string) (string, error){
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

func GetFileByte(fs billy.Filesystem, path string) ([]byte, error){
	//Create array the size of the file content
	stat, err := fs.Stat(path)
	if err != nil{
		return nil, err
	}
	bytes := make([]byte, stat.Size())

	//get file from memory
	file, err := fs.Open(path)
	if err != nil{
		return nil, err
	}
	file.Read(bytes)

	return bytes, nil
}

func UpdateAndPush(fs billy.Filesystem, repository *git.Repository, configPath string, newConfig string, auth *http.BasicAuth, oldTag string, newTag string) error{

	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Remove(configPath)
	if err != nil {
		return err
	}

	configFile, err := fs.Create(configPath)
	if err != nil {
		return err
	}
	
	b := []byte(newConfig)
	b = bytes.Trim(b, "\x00")
	_, err = configFile.Write(b)
	if err != nil {
		return err
	}

	_, err = worktree.Add(configPath)
	if err != nil {
		return err
	}

	//should be edited when context is added
	signature := &object.Signature{
		Name: auth.Username,
		Email: auth.Username,
		When: time.Now(),
	}

	_, err = worktree.Commit("Config file updated with Bazel: "+oldTag+" ---> "+newTag, &git.CommitOptions{
		Author: signature,
	})
	if err != nil {
		return err
	}

	pushOption := &git.PushOptions{
		RemoteName: "origin",
		Auth: auth,
	}

	err = repository.Push(pushOption)
	if err != nil {
		return err
	}

	return nil

}
