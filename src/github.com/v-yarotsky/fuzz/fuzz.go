package fuzz

import (
	"errors"
	"math"
	"strings"
)

func Score(choice, query string) float64 {
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

	score := float64(len(query)) / float64(matchLength) // Penalize longer matches.
	return score / float64(len(choice))                 // Normalize vs. the length of the choice, panalizing longer strings.
	return 0.0
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
