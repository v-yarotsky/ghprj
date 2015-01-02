package github

import (
	"fmt"
	"io/ioutil"
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
	return c.request("GET", requestPath)
}

func (c *HttpApi) request(requestType string, requestPath string) (*Response, error) {
	req, err := http.NewRequest(requestType, c.fullUrl(requestPath), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("User-Agent", "gh-prj")
	req.Header.Add("Authorization", "token "+c.accessToken)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Println("request", requestPath, "response status code", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	paging := new(PaginationInfo)
	c.populatePaging(resp, paging)

	return &Response{Body: body, Paging: paging}, nil
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
