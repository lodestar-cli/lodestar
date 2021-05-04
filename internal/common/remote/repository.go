package remote

import (
	"context"
	"errors"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/lodestar-cli/lodestar/internal/common/auth"
	"strings"
)

type LodestarRepository struct {
	Url         string
	Credentials auth.GitCredentials
	FileSystem  billy.Filesystem
	Storage     *memory.Storage
	Repository  *git.Repository
	Worktree    *git.Worktree
}

func NewLodestarRepository(url string, credentials auth.GitCredentials) (*LodestarRepository, error) {
	if url == ""{
		return nil, errors.New("cannot clone repository. Url cannot be blank")
	} else if !strings.Contains(url, "https://"){
		if !strings.Contains(url, "http://") {
			return nil, errors.New("cannot clone repository. Url must be of schema http or https")
		}
	}

	r := LodestarRepository{
		Url:         url,
		FileSystem:  memfs.New(),
		Credentials: credentials,
		Storage:     memory.NewStorage(),
		Repository: new(git.Repository),
		Worktree:   new(git.Worktree),
	}

	err := r.setRepository()
	if err != nil {
		return nil, err
	}
	err = r.setWorktree()
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *LodestarRepository) CommitFiles(commitMessage string, updatedFiles ...LodestarFile) error {

	for _, file := range updatedFiles {

		switch f := file.(type) {
		case *ManagementFile:
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
		case *AppStateFile:
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
		}
	}

	commitOptions, err := r.Credentials.CreateCommitOptions()
	if err != nil {
		return err
	}
	_, err = r.Worktree.Commit(commitMessage, commitOptions)
	if err != nil {
		return err
	}

	return nil
}

func (r *LodestarRepository) Push() error {
	pushOption, err := r.Credentials.CreatePushOptions()
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
	if err != nil {
		return err
	}

	r.Repository, err = git.CloneContext(context.TODO(), r.Storage, r.FileSystem, cloneOptions)
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
