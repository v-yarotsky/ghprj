package github

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"

	"github.com/cli/oauth/device"
	"github.com/pkg/browser"
)

const OauthClientID = "Iv1.47ce08a23a551fc7"

type accessTokenObtainer func() (accessToken string, err error)

func DeviceFlowAccessTokenObtainer() (string, error) {
	scopes := []string{"repo"}

	code, err := device.RequestCode(http.DefaultClient, "https://github.com/login/device/code", OauthClientID, scopes)
	if err != nil {
		return "", fmt.Errorf("Error while requesting device authorization code: %w", err)
	}

	log.Printf("Enter the following code in the browser: %s", code.UserCode)

	u, _ := url.Parse(code.VerificationURI)
	q := u.Query()
	q.Add("code", code.UserCode)
	u.RawQuery = q.Encode()
	log.Printf("Navigating to %s", u)
	err = browser.OpenURL(u.String())
	if err != nil {
		return "", fmt.Errorf("Failed to launch browser: %w", err)
	}

	res, err := device.Wait(context.TODO(), http.DefaultClient, "https://github.com/login/oauth/access_token", device.WaitOptions{
		ClientID:   OauthClientID,
		DeviceCode: code,
	})
	if err != nil {
		return "", fmt.Errorf("Failed to obtain access token: %w", err)
	}

	return res.Token, nil
}

func StaticToken(token string) func() (string, error) {
	return func() (string, error) {
		return token, nil
	}
}

type Authenticator struct {
	StoreDir string
	GetToken accessTokenObtainer
	Force    bool
}

func NewAuthenticator(cb accessTokenObtainer, force bool) *Authenticator {
	storeDir := alfredGithubDir("")
	return &Authenticator{StoreDir: storeDir, GetToken: cb, Force: force}
}

func (a *Authenticator) AccessToken() string {
	storeFile := a.StoreDir + "/auth_token"
	token, err := os.ReadFile(storeFile)
	if err != nil || a.Force {
		token = []byte(a.obtainAccessToken())
		os.MkdirAll(a.StoreDir, 0700)

		if err = os.WriteFile(storeFile, token, 0600); err != nil {
			log.Printf("Could not store authentication token: %s", err)
		}
	}
	return string(token)
}

func (a *Authenticator) obtainAccessToken() string {
	accessToken, err := a.GetToken()
	if err != nil {
		log.Fatalf("Could not obtain access token: %s", err)
	}
	return accessToken
}

func getNoteForAuthToken() string {
	user, _ := user.Current()
	username := "<unknown>"
	if user != nil {
		username = user.Username
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "<unknown>"
	}
	return fmt.Sprintf("gh-prj (%s@%s)", username, hostname)
}
