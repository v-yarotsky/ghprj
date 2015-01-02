package fuzz

import (
	"../github"
	"errors"
	"math"
	"sort"
	"strings"
)

func FilterRepos(repos []github.Repo, query string) []github.Repo {
	scoredRepos := NewScoredRepos(repos, query)
	sort.Sort(sort.Reverse(scoredRepos))
	var filtered []github.Repo

	for i, r := range scoredRepos.Repos {
		if scoredRepos.Scores[i] > 0.0 {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

type ScoredRepos struct {
	Repos  []github.Repo
	Scores []float64
	Query  string
}

func NewScoredRepos(repos []github.Repo, query string) ScoredRepos {
	scores := make([]float64, len(repos))

	for i, r := range repos {
		scores[i] = score(r.Name, query)
	}

	return ScoredRepos{Repos: repos, Scores: scores, Query: query}
}

func (a ScoredRepos) Len() int           { return len(a.Repos) }
func (a ScoredRepos) Less(i, j int) bool { return a.Scores[i] < a.Scores[j] }
func (a ScoredRepos) Swap(i, j int) {
	a.Repos[i], a.Repos[j] = a.Repos[j], a.Repos[i]
	a.Scores[i], a.Scores[j] = a.Scores[j], a.Scores[i]
}

func score(choice, query string) float64 {
	if len(query) == 0 {
		return 1.0
	}

	if len(choice) == 0 {
		return 0.0
	}

	choice = strings.ToLower(choice)
	query = strings.ToLower(query)

	matchLength, err := computeMatchLength(choice, query)

	if err != nil {
		return 0.0
	}

	stringScore := float64(len(query)) / float64(matchLength) // Penalize longer matches.
	return stringScore / float64(len(choice))                 // Normalize vs. the length of the choice, panalizing longer strings.
}

func computeMatchLength(str, chars string) (int, error) {
	runes := []rune(chars)
	firstChar := runes[0]
	restChars := runes[1:]

	firstIndexes := findCharInString(firstChar, str)

	matchLength := math.MaxInt32

	for _, i := range firstIndexes {
		lastIndex := findEndOfMatch(str, restChars, i)
		if lastIndex != -1 {
			newMatchLength := lastIndex - i + 1
			if matchLength > newMatchLength {
				matchLength = newMatchLength
			}
		}
	}

	if matchLength == math.MaxInt32 {
		return -1, errors.New("did not match")
	}
	return matchLength, nil
}

func findCharInString(chr rune, str string) []int {
	indexes := []int{}

	for i, cur := range []rune(str) {
		if chr == cur {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

func findEndOfMatch(str string, chars []rune, firstIndex int) int {
	lastIndex := firstIndex
	runes := []rune(str)

	for _, chr := range chars {
		index := -1
		for i, r := range runes[(lastIndex + 1):] {
			if chr == r {
				index = i
				break
			}
		}

		if index == -1 {
			return -1
		}
		lastIndex = index
	}

	return lastIndex
}
