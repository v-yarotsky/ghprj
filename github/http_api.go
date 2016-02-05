package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const GITHUB_API_ROOT = "https://api.github.com"

var err2FAOTPRequired = errors.New("Valid 2FA OTP is required")

type HttpApi struct {
	accessToken string
	credentials Credentials
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

func (c *HttpApi) Put(requestPath string, body interface{}) (*Response, error) {
	return c.request("PUT", requestPath, body)
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

	if c.accessToken == "" {
		req.SetBasicAuth(c.credentials.Username, c.credentials.Password)
		req.Header.Add("X-Github-OTP", c.credentials.TwoFactorToken)
	} else {
		req.Header.Add("Authorization", "token "+c.accessToken)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 && strings.Contains(resp.Header.Get("X-GitHub-OTP"), "required;") {
		return nil, err2FAOTPRequired
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	paging := new(PaginationInfo)
	c.populatePaging(resp, paging)

	return &Response{Body: responseBody, Paging: paging}, nil
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
