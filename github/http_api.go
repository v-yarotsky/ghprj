package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const GITHUB_API_ROOT = "https://api.github.com"

type HttpApi struct {
	accessToken string
}

type Response struct {
	Body   []byte
	Paging *PaginationInfo
}

type PaginationInfo struct {
	NextPagePath string
}

func (c *HttpApi) Get(requestPath string) (*Response, error) {
	return c.request("GET", requestPath, nil)
}

func (c *HttpApi) Delete(requestPath string) (*Response, error) {
	return c.request("DELETE", requestPath, nil)
}

func (c *HttpApi) Post(requestPath string, body interface{}) (*Response, error) {
	return c.request("POST", requestPath, body)
}

func (c *HttpApi) request(requestType string, requestPath string, body interface{}) (*Response, error) {
	var bodyJson io.Reader

	if body != nil {
		body, err := json.Marshal(&body)
		if err != nil {
			return nil, fmt.Errorf("Failed to prepare request: %s", err)
		}

		bodyJson = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(requestType, c.fullUrl(requestPath), bodyJson)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("User-Agent", "gh-prj")

	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	response := &Response{}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response.Body = responseBody

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return response, fmt.Errorf("github API request failed: %s", resp.Status)
	}

	paging := new(PaginationInfo)
	c.populatePaging(resp, paging)
	response.Paging = paging

	return response, nil
}

func (c *HttpApi) fullUrl(path string) string {
	return GITHUB_API_ROOT + path
}

func (c *HttpApi) populatePaging(response *http.Response, paging *PaginationInfo) {
	if links, ok := response.Header["Link"]; ok && len(links) > 0 {
		for _, link := range strings.Split(links[0], ", ") {
			r := regexp.MustCompile(`^<(?P<link>.*?)>; rel="(?P<rel>.*?)"$`)
			matches := r.FindStringSubmatch(link)
			nextLink := matches[1]
			rel := matches[2]

			if rel == "next" {
				paging.NextPagePath = strings.Replace(nextLink, GITHUB_API_ROOT, "", 1)
				break
			}
		}
	}
}
