package main

// cache
// alfred-compatible output
// selecta-compatible output
// get token

import (
	"./formatter"
	"./fuzz"
	"./github"
	"log"
	"os"
)

func main() {
	c, _ := github.NewCachingClient()
	repos, err := c.UserAndOrgRepos()

	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		repos = fuzz.FilterRepos(repos, os.Args[1])
	}

	results, err := (&formatter.Alfred{}).FormattedResults(repos)

	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(results)
}
