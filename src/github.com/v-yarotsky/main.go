package main

// get token
// alfred-compatible output
// selecta-compatible output
// cache

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	c := &Client{
		&GithubHttpClient{os.Getenv("GH_PRJ_GITHUB_TOKEN")},
	}

	repos, err := c.UserAndOrgRepos()

	fmt.Println("resulting repos", repos, err, "count: ", len(repos))
}

type Client struct {
	api *GithubHttpClient
}

type Org struct {
	Login string `json:"login"`
}

type Repo struct {
	Name    string `json:"name"`
	HtmlUrl string `json:"html_url"`
}

func (c *Client) UserAndOrgRepos() ([]Repo, error) {
	orgs, err := c.Orgs()

	if err != nil {
		return nil, err
	}

	ownRepos, err := c.Repos()

	if err != nil {
		return nil, err
	}

	remainingOrgs := len(orgs)
	orgReposChan := make(chan []Repo)
	orgReposErrChan := make(chan error)

	for _, org := range orgs {
		go func(org Org) {
			orgRepos, err := c.OrgRepos(org.Login)
			if err != nil {
				orgReposErrChan <- err
				return
			}
			orgReposChan <- orgRepos
		}(org)
	}

	repos := ownRepos
	for {
		select {
		case orgRepos := <-orgReposChan:
			repos = append(repos, orgRepos...)
			remainingOrgs--
		case err = <-orgReposErrChan:
			return nil, err
		}

		if remainingOrgs == 0 {
			break
		}
	}

	return repos, nil
}

func (c *Client) Repos() ([]Repo, error) {
	return c.paginatedRepos("/user/repos?per_page=100")
}

func (c *Client) OrgRepos(orgLogin string) ([]Repo, error) {
	return c.paginatedRepos("/orgs/" + orgLogin + "/repos?per_page=100")
}

func (c *Client) paginatedRepos(initialPath string) ([]Repo, error) {
	result := []Repo{}

	for nextPagePath := initialPath; len(nextPagePath) > 0; {
		resp, err := c.api.Get(nextPagePath)

		if err != nil {
			return result, err
		}

		repos := []Repo{}
		err = json.Unmarshal(resp.Body, &repos)

		if err != nil {
			return result, err
		}

		result = append(result, repos...)

		nextPagePath = resp.Paging.NextPagePath
	}

	return result, nil
}

func (c *Client) Orgs() ([]Org, error) {
	resp, err := c.api.Get("/user/orgs")

	if err != nil {
		return nil, err
	}

	orgs := []Org{}
	err = json.Unmarshal(resp.Body, &orgs)

	return orgs, err
}
