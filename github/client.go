package github

import (
	"encoding/json"
	"fmt"
)

type Client struct {
	api *HttpApi
}

type Repo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HtmlUrl  string `json:"html_url"`
}

func NewClient(accessToken string) *Client {
	return &Client{
		&HttpApi{accessToken: accessToken},
	}
}

func (c *Client) GetLogin() (string, error) {
	res, err := c.api.Get("/user")
	if err != nil {
		return "", fmt.Errorf("Failed to get authenticated user: %w", err)
	}
	var me struct {
		Login string `json:"login"`
	}
	err = json.Unmarshal(res.Body, &me)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal get current user response: %w", err)
	}
	return me.Login, nil
}

func (c *Client) UserAndOrgRepos() ([]Repo, error) {
	return c.paginatedRepos("/user/repos?per_page=100")
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
