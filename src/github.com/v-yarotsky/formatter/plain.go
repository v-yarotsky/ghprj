package formatter

import (
	"../github"
	"bytes"
)

type Plain struct{}

func (a *Plain) FormattedResults(repos []github.Repo) ([]byte, error) {
	buf := bytes.NewBufferString("")

	for _, repo := range repos {
		buf.WriteString(repo.HtmlUrl + "\n")
	}

	return buf.Bytes(), nil
}
