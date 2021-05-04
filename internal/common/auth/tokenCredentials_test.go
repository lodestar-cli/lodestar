package auth

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"reflect"
	"testing"
)

func TestTokenCredentials_CreateCloneOptions(t *testing.T) {

	tc := TokenCredentials{
		Username: "testUser",
		Token:    "testToken",
	}
	var ta transport.AuthMethod

	ta = &http.BasicAuth{
		Username: tc.Username,
		Password: tc.Token,
	}

	url := "https:testurl.com"

	co, err := tc.CreateCloneOptions(url)
	if err != nil {
		t.Error("Failed Creating Clone Options for TokenCredentials")
	}

	if !reflect.DeepEqual(co.Auth, ta) {
		t.Error("Clone Option created do not reflect the username and token given")
	}
}

func TestTokenCredentials_CreateCommitOptions(t *testing.T) {

	tc := TokenCredentials{
		Username: "testUser",
		Token:    "testToken",
	}

	co, err := tc.CreateCommitOptions()
	if err != nil {
		t.Error("Failed Creating Commit Options for TokenCredentials")
	}

	if co.Author.Name != tc.Username || co.Author.Email != tc.Username {
		t.Error("Commit Option created do not reflect the username and token given")
	}
}

func TestTokenCredentials_CreatePushOptions(t *testing.T) {

	tc := TokenCredentials{
		Username: "testUser",
		Token:    "testToken",
	}
	var ta transport.AuthMethod

	ta = &http.BasicAuth{
		Username: tc.Username,
		Password: tc.Token,
	}

	co, err := tc.CreatePushOptions()
	if err != nil {
		t.Error("Failed Creating Clone Options for TokenCredentials")
	}

	if !reflect.DeepEqual(co.Auth, ta) {
		t.Error("Clone Option created do not reflect the username and token given")
	}
}
