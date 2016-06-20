package github

import (
	"io/ioutil"
	"log"
	"os"
)

type Credentials struct {
	Username       string
	Password       string
	TwoFactorToken string
}

type credentialsCallback func(*Credentials, bool) error

type Authenticator struct {
	StoreDir       string
	GetCredentials credentialsCallback
	Force          bool
}

func NewAuthenticator(cb credentialsCallback, force bool) *Authenticator {
	storeDir := alfredGithubDir("")
	return &Authenticator{StoreDir: storeDir, GetCredentials: cb, Force: force}
}

func (a *Authenticator) AccessToken() string {
	storeFile := a.StoreDir + "/auth_token"
	token, err := ioutil.ReadFile(storeFile)
	if err != nil || a.Force {
		token = []byte(a.obtainAccessToken())
		os.MkdirAll(a.StoreDir, 0700)

		if err = ioutil.WriteFile(storeFile, token, 0600); err != nil {
			log.Printf("Could not store authentication token: %s", err)
		}
	}
	return string(token)
}

func (a *Authenticator) obtainAccessToken() string {
	credentials := Credentials{}
	a.mustGetCredentials(&credentials, false)

	authorization, err := doObtainAccessToken(credentials)
	for err != nil {
		switch err {
		case err2FAOTPRequired:
			a.mustGetCredentials(&credentials, true)
			authorization, err = doObtainAccessToken(credentials)
		case nil:
			break
		default:
			log.Fatalf("Failed to create authorization: %s", err)
		}
	}
	return authorization.Token
}

func (a *Authenticator) mustGetCredentials(c *Credentials, twoFactor bool) {
	if err := a.GetCredentials(c, twoFactor); err != nil {
		log.Fatalf("Could not get credentials: %s", err)
	}
}

func doObtainAccessToken(c Credentials) (*Authorization, error) {
	client := NewBasicAuthClient(c)
	return client.ForceCreateAuthorization([]string{"repo"}, "gh-prj")
}
