package repo

import (
	"context"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/lodestar-cli/lodestar/internal/cli/files"
	"github.com/lodestar-cli/lodestar/internal/common/auth"
)

type LodestarRepository struct {
	Url         string
	Credentials auth.GitCredentials
	FileSystem  billy.Filesystem
	Storer      *memory.Storage
	Repository  *git.Repository
	Worktree    *git.Worktree
}

func NewLodestarRepository(url string, credentials auth.GitCredentials) (*LodestarRepository, error) {
	r := LodestarRepository{
		Url: url,
		FileSystem: memfs.New(),
		Credentials: credentials,
		Storer: memory.NewStorage(),
	}

	err := r.setRepository()
	if err != nil {
		return nil, err
	}
	r.setWorktree()
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *LodestarRepository) CommitFiles(commitMessage string, updatedFiles ...files.LodestarFile) error {

	for _, file := range updatedFiles {

		switch f := file.(type) {
		case *files.ManagementFile:
			_, err := r.Worktree.Remove(f.Path)
			if err != nil {
				return err
			}

			configFile, err := r.FileSystem.Create(f.Path)
			if err != nil {
				return err
			}

			_, err = configFile.Write([]byte(f.StringContent))
			if err != nil {
				return err
			}

			_, err = r.Worktree.Add(f.Path)
			if err != nil {
				return err
			}
		default:
			f.Print()
		}
	}

	commitOptions, err := r.Credentials.CreateCommitOptions(r.Url)
	if err != nil {
		return err
	}
	_, err = r.Worktree.Commit(commitMessage, commitOptions)
	if err != nil {
		return err
	}

	return nil
}

func (r *LodestarRepository) Push() error{
	pushOption, err := r.Credentials.CreatePushOptions(r.Url)
	if err != nil {
		return err
	}

	err = r.Repository.Push(pushOption)
	if err != nil {
		return err
	}

	return nil

}


func (r *LodestarRepository) setRepository() error {

	cloneOptions, err := r.Credentials.CreateCloneOptions(r.Url)
	if err != nil{
		return err
	}

	r.Repository, err = git.CloneContext(context.TODO(),r.Storer, r.FileSystem, cloneOptions)
	if err != nil {
		return err
	}

	return nil
}


func (r *LodestarRepository) setWorktree() error {
	var err error

	r.Worktree, err = r.Repository.Worktree()
	if err != nil {
		return err
	}

	return nil
}