package github

import (
	"io/ioutil"
	"log"
	"os"
)

type usernamePasswordCallback func() (string, string, error)

type Authenticator struct {
	StoreFile              string
	GetUsernameAndPassword usernamePasswordCallback
}

func NewAuthenticator(cb usernamePasswordCallback) *Authenticator {
	storeFile := os.Getenv("HOME") + "/.gh-prj/auth_token"
	return &Authenticator{StoreFile: storeFile, GetUsernameAndPassword: cb}
}

func (a *Authenticator) AccessToken() string {
	token, err := ioutil.ReadFile(a.StoreFile)
	if err != nil {
		token = []byte(a.obtainAccessToken())
		ioutil.WriteFile(a.StoreFile, token, 0600)
	}
	return string(token)
}

func (a *Authenticator) obtainAccessToken() string {
	username, password, err := a.GetUsernameAndPassword()

	if err != nil {
		log.Fatalf("Failed to create authorization: %s", err)
	}

	client := NewBasicAuthClient(username, password)

	authorization, err := client.GetOrCreateAuthorization([]string{"repo"}, "gh-prj")

	if err != nil {
		log.Fatalf("Failed to create authorization: %s", err)
	}

	return authorization.Token
}
