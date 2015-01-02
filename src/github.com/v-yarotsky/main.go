package main

// alfred-compatible output
// selecta-compatible output
// get token
// cache

import (
	"./formatter"
	"./fuzz"
	"./github"
	"log"
	"os"
	"sort"
)

func main() {
	c, _ := github.NewClient()
	repos, err := c.UserAndOrgRepos()

	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		repos = filterRepos(repos, os.Args[1])
	}

	results, err := (&formatter.Alfred{}).FormattedResults(repos)

	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(results)
}

type ByScore struct {
	Repos  []github.Repo
	Scores []float64
	Query  string
}

func NewByScore(repos []github.Repo, query string) ByScore {
	scores := make([]float64, len(repos))

	for i, r := range repos {
		scores[i] = fuzz.Score(r.Name, query)
	}

	return ByScore{Repos: repos, Scores: scores, Query: query}
}

func (a ByScore) Len() int           { return len(a.Repos) }
func (a ByScore) Less(i, j int) bool { return a.Scores[i] < a.Scores[j] }
func (a ByScore) Swap(i, j int) {
	a.Repos[i], a.Repos[j] = a.Repos[j], a.Repos[i]
	a.Scores[i], a.Scores[j] = a.Scores[j], a.Scores[i]
}

func filterRepos(repos []github.Repo, query string) []github.Repo {
	scoredRepos := NewByScore(repos, query)
	sort.Sort(sort.Reverse(scoredRepos))
	var filtered []github.Repo

	for i, r := range scoredRepos.Repos {
		if scoredRepos.Scores[i] > 0.0 {
			filtered = append(filtered, r)
		}
	}

	return filtered
}
