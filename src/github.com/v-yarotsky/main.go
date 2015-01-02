package main

// cache
// alfred-compatible output
// selecta-compatible output
// get token

import (
	"./formatter"
	"./fuzz"
	"./github"
	"flag"
	"log"
	"os"
)

func main() {
	expirePtr := flag.Bool("expire", false, "Expire caches")
	flag.Parse()

	c, _ := github.NewCachingClient()

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
