package github

import (
	"encoding/json"
)

const clientId = "5351a4cf6969f32fe1c6"
const clientSecret = "c3c8cf8e2c35e7c9406618a6dec0abd0d35125d8"

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

func NewClient(accessToken string) *Client {
	return &Client{
		&HttpApi{accessToken: &accessToken, username: nil, password: nil},
	}
}

func NewBasicAuthClient(username, password string) *Client {
	return &Client{
		&HttpApi{accessToken: nil, username: &username, password: &password},
	}
}

type AuthenticationData struct {
	Secret string   `json:"client_secret"`
	Scopes []string `json:"scopes"`
	Note   string   `json:"note"`
}

type Authorization struct {
	Token string `json:"token"`
}

func (c *Client) GetOrCreateAuthorization(scopes []string, note string) (*Authorization, error) {
	resp, err := c.api.Put("/authorizations/clients/"+clientId,
		&AuthenticationData{
			Secret: clientSecret,
			Scopes: scopes,
			Note:   note,
		})

	if err != nil {
		return nil, err
	}

	authorization := &Authorization{}
	json.Unmarshal(resp.Body, authorization)
	return authorization, err
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
