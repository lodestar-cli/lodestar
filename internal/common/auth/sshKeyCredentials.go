package auth

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"io/ioutil"
	"time"
)

type SSHKeyCredentials struct {
	Username string
	KeyPath string
	KeyPassword string
	PublicKeys *ssh.PublicKeys
}

func (s *SSHKeyCredentials) CreateCloneOptions(url string) (*git.CloneOptions, error){
	var err error

	if s.PublicKeys == nil {
		s.PublicKeys, err = s.getLocalPublicKeys()
		if err != nil {
			return nil, err
		}
	}

	cloneOptions := &git.CloneOptions{
		URL : url,
		Auth: s.PublicKeys,
	}

	return cloneOptions, nil
}

func (s *SSHKeyCredentials) CreateCommitOptions() (*git.CommitOptions, error){
	signature := &object.Signature{
		Name: s.Username,
		Email: s.Username,
		When: time.Now(),
	}

	commitOptions := git.CommitOptions{
		Author: signature,
	}

	return &commitOptions, nil
}

func (s *SSHKeyCredentials) CreatePushOptions() (*git.PushOptions, error){
	var err error

	if s.PublicKeys == nil {
		s.PublicKeys, err = s.getLocalPublicKeys()
		if err != nil {
			return nil, err
		}
	}

	pushOptions := &git.PushOptions{
		RemoteName: "origin",
		Auth: s.PublicKeys,
	}

	return pushOptions, nil
}

func (s *SSHKeyCredentials) getLocalPublicKeys() (*ssh.PublicKeys, error) {
	sshKey, _ := ioutil.ReadFile(s.KeyPath)
	auth, err := ssh.NewPublicKeys(s.Username, []byte(sshKey), s.KeyPassword)
	if err != nil {
		return nil, err
	}

	return auth, nil
}
