package github

import (
	"encoding/json"
	"fmt"
	"reflect"
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

func NewBasicAuthClient(c Credentials) *Client {
	return &Client{&HttpApi{credentials: c}}
}

type AuthenticationData struct {
	Scopes []string `json:"scopes"`
	Note   string   `json:"note"`
}

type Authorization struct {
	ID     int      `json:"id"`
	Token  string   `json:"token"`
	Scopes []string `json:"scopes"`
	Note   string   `json:"note"`
}

func (c *Client) ForceCreateAuthorization(scopes []string, note string) (*Authorization, error) {
	err := c.findAndDeleteExistingAuthorization(scopes, note)
	if err != nil {
		return nil, err
	}

	resp, err := c.api.Post("/authorizations", &AuthenticationData{
		Scopes: scopes,
		Note:   note,
	})
	if err != nil {
		return nil, err
	}

	authorization := &Authorization{}
	err = json.Unmarshal(resp.Body, authorization)
	if err != nil {
		return nil, err
	}
	return authorization, nil
}

func (c *Client) findAndDeleteExistingAuthorization(scopes []string, note string) error {
	resp, err := c.api.Get("/authorizations")
	if err != nil {
		return err
	}
	authorizations := []Authorization{}
	err = json.Unmarshal(resp.Body, &authorizations)
	if err != nil {
		return err
	}
	for _, a := range authorizations {
		if reflect.DeepEqual(a.Scopes, scopes) && a.Note == note {
			_, err := c.api.Delete(fmt.Sprintf("/authorizations/%d", a.ID))
			return err
		}
	}
	return nil
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
