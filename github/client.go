package github

import (
	"encoding/json"
	"fmt"
)

const clientId = "5351a4cf6969f32fe1c6"
const clientSecret = "c3c8cf8e2c35e7c9406618a6dec0abd0d35125d8"

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

func NewBasicAuthClient(c Credentials) *Client {
	return &Client{&HttpApi{credentials: c}}
}

type AuthenticationData struct {
	Secret string   `json:"client_secret"`
	Scopes []string `json:"scopes"`
	Note   string   `json:"note"`
}

type Authorization struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func (c *Client) ForceCreateAuthorization(scopes []string, note string) (*Authorization, error) {
	authorization, err := c.doGetOrCreateOrganizations(scopes, note)
	if err != nil {
		return nil, err
	}

	if authorization.Token == "" {
		_, err := c.api.Delete(fmt.Sprintf("/authorizations/%d", authorization.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to delete lost authorization: %s", err)
		}
		authorization, err = c.doGetOrCreateOrganizations(scopes, note)
	}

	return authorization, err
}

func (c *Client) doGetOrCreateOrganizations(scopes []string, note string) (*Authorization, error) {
	resp, err := c.api.Put(fmt.Sprintf("/authorizations/clients/%s", clientId),
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
	return authorization, nil
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
