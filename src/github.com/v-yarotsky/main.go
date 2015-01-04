package main

import (
	"./formatter"
	"./fuzz"
	"./github"
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

func main() {
	expirePtr := flag.Bool("expire", false, "Expire caches")
	flag.Parse()

	accessToken := github.NewAuthenticator(func() (string, string, error) {
		if !terminal.IsTerminal(int(os.Stdin.Fd())) {
			log.Fatal("Can not ask for username/password - not an interactive terminal")
		}

		fmt.Printf("Username: ")
		var username, password string
		fmt.Scanf("%s", &username)
		fmt.Printf("Password: ")
		password = string(gopass.GetPasswd())
		return username, password, nil
	}).AccessToken()

	c, _ := github.NewCachingClient(accessToken)

	if *expirePtr {
		err := c.Expire()
		if err != nil {
			log.Printf("Failed to expire cache: %s", err)
		}
	}

	repos, err := c.UserAndOrgRepos()

	if err != nil {
		log.Fatal(err)
	}

	if len(flag.Args()) > 0 {
		repos = fuzz.FilterRepos(repos, flag.Arg(0))
	}

	results, err := (&formatter.Alfred{}).FormattedResults(repos)

	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(results)
}
