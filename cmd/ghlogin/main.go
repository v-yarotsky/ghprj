package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/v-yarotsky/ghprj/github"
)

func main() {
	var personalToken string
	flag.StringVar(&personalToken, "token", "", "Personal Access Token with `repo` scope")
	flag.Parse()

	accessToken := github.NewAuthenticator(github.StaticToken(personalToken), true).AccessToken()

	if _, err := os.Stat("info.plist"); err == nil {
		log.Println("Patching info.plist with github username")
		c, _ := github.NewCachingClient(accessToken)
		username, err := c.Client.GetLogin()
		if err != nil {
			log.Fatalf("Failed to get current username: %s", err)
		}
		log.Printf("Detected username: %s", username)
		exec.Command("/usr/libexec/PlistBuddy", "-c", fmt.Sprintf("Set :variables:username \"%s\"", username), "info.plist")
	} else {
		log.Println("info.plist not found; not patching it with github username")
	}
}
