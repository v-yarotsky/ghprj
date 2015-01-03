package github

import (
	"encoding/json"
)

type Client struct {
	api *HttpApi
}

type Org struct {
	Login string `json:"login"`
}

type Repo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HtmlUrl  string `json:"html_url"`
}

func NewClient() (*Client, error) {
	return &Client{
		&HttpApi{AccessToken()},
	}, nil
}

func (c *Client) UserAndOrgRepos() (*[]Repo, error) {
	orgs, err := c.orgs()

	if err != nil {
		return nil, err
	}

	ownRepos, err := c.repos()

	if err != nil {
		return nil, err
	}

	remainingOrgs := len(orgs)
	orgReposChan := make(chan []Repo)
	orgReposErrChan := make(chan error)

	for _, org := range orgs {
		go func(org Org) {
			orgRepos, err := c.orgRepos(org.Login)
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

	return &repos, nil
}

func (c *Client) repos() ([]Repo, error) {
	return c.paginatedRepos("/user/repos?per_page=100")
}

func (c *Client) orgRepos(orgLogin string) ([]Repo, error) {
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

func (c *Client) orgs() ([]Org, error) {
	resp, err := c.api.Get("/user/orgs")

	if err != nil {
		return nil, err
	}

	orgs := []Org{}
	err = json.Unmarshal(resp.Body, &orgs)

	return orgs, err
}
