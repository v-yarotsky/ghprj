package formatter

import (
	"encoding/xml"

	"github.com/v-yarotsky/gh-prj/github"
)

type Alfred struct{}

type items struct {
	Items []item `xml:"item"`
}

type item struct {
	Title    string `xml:"title"`
	Subtitle string `xml:"subtitle"`
	Icon     string `xml:"icon"`
	Uid      string `xml:"uid",attr`
	Valid    bool   `xml:"valid",attr`
	Arg      string `xml:"arg",attr`
}

func (a *Alfred) FormattedResults(repos []github.Repo) ([]byte, error) {
	itemsArr := make([]item, len(repos))

	for i, repo := range repos {
		itemsArr[i] = item{
			Title:    repo.Name,
			Subtitle: repo.FullName,
			Icon:     "repo.png",
			Uid:      repo.HtmlUrl,
			Valid:    true,
			Arg:      repo.HtmlUrl,
		}
	}

	items := &items{Items: itemsArr}

	result, err := xml.MarshalIndent(items, "", "    ")

	if err != nil {
		return nil, err
	}

	return result, nil
}
