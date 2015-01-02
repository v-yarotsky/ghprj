package main

// get token
// alfred-compatible output
// selecta-compatible output
// cache

import (
	"./github"
	"fmt"
)

func main() {
	c, _ := github.NewClient()
	repos, err := c.UserAndOrgRepos()
	fmt.Println("resulting repos", repos, err, "count: ", len(repos))
}
