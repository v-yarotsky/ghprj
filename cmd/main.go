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

	accessToken := github.NewAuthenticator(func(c *github.Credentials, twoFactor bool) error {
		if !terminal.IsTerminal(int(os.Stdin.Fd())) {
			return fmt.Errorf("Can not ask for credentials - not an interactive terminal")
		}

		var err error
		if twoFactor {
			fmt.Printf("Two-Factor OTP: ")
			fmt.Scanf("%s", &c.TwoFactorToken)
		} else {
			fmt.Printf("Username: ")
			fmt.Scanf("%s", &c.Username)
			fmt.Printf("Password: ")
			var password []byte
			password, err = gopass.GetPasswd()
			c.Password = string(password)
		}
		return err
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
