package ghprj

import (
	"bytes"
	"errors"
	"math"
	"sort"
	"strings"

	"github.com/v-yarotsky/ghprj/github"
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

	score := float64(len(query)) / float64(matchLength)
	normalizationFactor := float64(len(query)) / float64(len(choice))
	normalizedScore := score * normalizationFactor
	return normalizedScore
}

func computeMatchLength(str, chars string) (int, error) {
	runes := []rune(chars)
	firstChar := runes[0]
	restChars := string(runes[1:])

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

func findEndOfMatch(str, chars string, firstIndex int) int {
	lastIndex := firstIndex
	byteStr := []byte(str)
	for _, chr := range chars {
		i := bytes.IndexRune(byteStr[(lastIndex+1):], chr)
		if i == -1 {
			return -1
		}
		lastIndex += i
	}

	return lastIndex + 1
}
