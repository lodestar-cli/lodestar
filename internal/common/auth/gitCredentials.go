package auth

import "github.com/go-git/go-git/v5"

type GitCredentials interface{
	CreateCloneOptions(url string) (*git.CloneOptions, error)
	CreateCommitOptions(url string) (*git.CommitOptions, error)
	CreatePushOptions(url string) (*git.PushOptions, error)
}
