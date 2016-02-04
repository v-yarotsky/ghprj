package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/howeyc/gopass"
	"github.com/v-yarotsky/ghprj"
	"github.com/v-yarotsky/ghprj/github"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	expirePtr := flag.Bool("expire", false, "Expire caches")
	flag.Parse()

	accessToken := github.NewAuthenticator(func() (string, string, error) {
		if !terminal.IsTerminal(int(os.Stdin.Fd())) {
			log.Fatal("Can not ask for username/password - not an interactive terminal")
		}

		fmt.Printf("Username: ")
		var username string
		fmt.Scanf("%s", &username)
		fmt.Printf("Password: ")
		password, err := gopass.GetPasswd()
		return username, string(password), err
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
		repos = ghprj.FilterRepos(repos, flag.Arg(0))
	}

	results, err := (&ghprj.Alfred{}).FormattedResults(repos)
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(results)
}
