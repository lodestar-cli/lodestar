package auth

import "github.com/go-git/go-git/v5"

type GitCredentials interface{
	CreateCloneOptions(url string) (*git.CloneOptions, error)
	CreateCommitOptions() (*git.CommitOptions, error)
	CreatePushOptions() (*git.PushOptions, error)
}
