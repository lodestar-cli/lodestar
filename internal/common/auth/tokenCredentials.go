package auth

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"time"
)

type TokenCredentials struct {
	Username string
	Token    string
}

func (t *TokenCredentials) CreateCloneOptions(url string) (*git.CloneOptions, error){

	auth := http.BasicAuth{
		Username: t.Username,
		Password: t.Token,
	}


	cloneOptions := &git.CloneOptions{
		URL : url,
		Auth: &auth,
	}

	return cloneOptions, nil
}

func (t *TokenCredentials) CreateCommitOptions(url string) (*git.CommitOptions, error){
	signature := &object.Signature{
		Name: t.Username,
		Email: t.Username,
		When: time.Now(),
	}

	commitOptions := git.CommitOptions{
		Author: signature,
	}

	return &commitOptions, nil
}

func (t *TokenCredentials) CreatePushOptions(url string) (*git.PushOptions, error){

	auth := &http.BasicAuth{
		Username: t.Username,
		Password: t.Token,
	}

	pushOptions := &git.PushOptions{
		RemoteName: "origin",
		Auth: auth,
	}

	return pushOptions, nil
}
