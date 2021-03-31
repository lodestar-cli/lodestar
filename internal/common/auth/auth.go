package auth

import (
	http "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CreateAuth(username string, token string) (*http.BasicAuth, error) {

	auth := new(http.BasicAuth)

	auth.Username=username
	auth.Password=token

	return auth, nil
}
