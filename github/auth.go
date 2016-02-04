package github

import (
	"io/ioutil"
	"log"
	"os"
)

type usernamePasswordCallback func() (string, string, error)

type Authenticator struct {
	StoreDir               string
	GetUsernameAndPassword usernamePasswordCallback
}

func NewAuthenticator(cb usernamePasswordCallback) *Authenticator {
	storeDir := os.Getenv("HOME") + "/.gh-prj"
	return &Authenticator{StoreDir: storeDir, GetUsernameAndPassword: cb}
}

func (a *Authenticator) AccessToken() string {
	storeFile := a.StoreDir + "/auth_token"
	token, err := ioutil.ReadFile(storeFile)
	if err != nil {
		token = []byte(a.obtainAccessToken())
		os.MkdirAll(a.StoreDir, 0700)

		if err = ioutil.WriteFile(storeFile, token, 0600); err != nil {
			log.Printf("Could not store authentication token: %s", err)
		}
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
