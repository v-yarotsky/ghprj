package ghprj

import (
	"encoding/xml"

	"github.com/v-yarotsky/ghprj/github"
)

type Alfred struct{}

type items struct {
	XMLName xml.Name `xml:"items"`
	Items   []item   `xml:"item"`
}

type item struct {
	XMLName      xml.Name `xml:"item"`
	UID          string   `xml:"uid,attr"`
	Autocomplete string   `xml:"autocomplete,attr"`
	Valid        bool     `xml:"valid,attr"`
	Title        string   `xml:"title"`
	Subtitle     string   `xml:"subtitle"`
	Icon         string   `xml:"icon"`
	Arg          string   `xml:"arg"`
}

func (a *Alfred) FormattedResults(repos []github.Repo) ([]byte, error) {
	itemsArr := make([]item, len(repos))

	for i, repo := range repos {
		itemsArr[i] = item{
			Title:        repo.Name,
			Subtitle:     repo.FullName,
			Icon:         "repo.png",
			UID:          repo.HtmlUrl,
			Valid:        true,
			Arg:          repo.HtmlUrl,
			Autocomplete: "true",
		}
	}

	items := &items{Items: itemsArr}

	result, err := xml.MarshalIndent(items, "", "    ")

	if err != nil {
		return nil, err
	}

	return result, nil
}
