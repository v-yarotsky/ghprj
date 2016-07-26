package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/v-yarotsky/ghprj"
	"github.com/v-yarotsky/ghprj/github"
)

func main() {
	var expire bool
	flag.BoolVar(&expire, "expire", false, "Expire caches")
	flag.Parse()

	accessToken := github.NewAuthenticator(func(c *github.Credentials, twoFactor bool) error {
		return fmt.Errorf("Not authenticated. Please log in with ghlogin")
	}, false).AccessToken()

	c, _ := github.NewCachingClient(accessToken)

	if expire {
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
