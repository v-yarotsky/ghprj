package github

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type CachingClient struct {
	CacheDir string
	Client   *Client
}

func NewCachingClient() (*CachingClient, error) {
	cacheDir := os.Getenv("HOME") + "/.gh-prj/caches"

	return &CachingClient{
		CacheDir: cacheDir,
		Client:   &Client{&HttpApi{AccessToken()}},
	}, nil
}

func (c *CachingClient) Expire() error {
	return os.RemoveAll(c.CacheDir)
}

func (c *CachingClient) UserAndOrgRepos() ([]Repo, error) {
	result := []Repo{}

	r, err := c.fetchCache("user_and_org_repos", &result, func() (interface{}, error) {
		return c.Client.UserAndOrgRepos()
	})

	if err != nil {
		return nil, err
	}

	return *r.(*[]Repo), nil
}

func (c *CachingClient) Repos() ([]Repo, error) {
	result := []Repo{}

	r, err := c.fetchCache("repos", &result, func() (interface{}, error) {
		return c.Client.Repos()
	})

	if err != nil {
		return nil, err
	}

	return *r.(*[]Repo), nil
}

func (c *CachingClient) OrgRepos(orgLogin string) ([]Repo, error) {
	result := []Repo{}

	r, err := c.fetchCache("org_repos_"+orgLogin, &result, func() (interface{}, error) {
		return c.Client.OrgRepos(orgLogin)
	})

	if err != nil {
		return nil, err
	}

	return *r.(*[]Repo), nil
}

func (c *CachingClient) Orgs() ([]Org, error) {
	result := []Org{}
	r, err := c.fetchCache("orgs", &result, func() (interface{}, error) {
		return c.Client.Orgs()
	})

	if err != nil {
		return nil, err
	}

	return *r.(*[]Org), nil
}

func (c *CachingClient) fetchCache(key string, obj interface{}, fetchFn func() (interface{}, error)) (interface{}, error) {
	err := c.readCache(key, &obj)
	if err != nil {
		log.Printf("Cache miss: %s, or error reading cache", key)
		obj, err := fetchFn()
		if err != nil {
			log.Fatalf("Failed to fetch %s: %s", key, err)
		}
		err = c.writeCache(key, &obj)

		if err != nil {
			log.Printf("Failed to write cache %s: %s", key, err)
		}
	}
	return obj, nil
}

func (c *CachingClient) readCache(key string, obj interface{}) error {
	data, err := ioutil.ReadFile(c.expandKey(key))

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, obj)

	if err != nil {
		return err
	}

	return nil
}

func (c *CachingClient) writeCache(key string, obj interface{}) error {
	err := os.MkdirAll(c.CacheDir, 0700)

	if err != nil {
		return err
	}

	data, err := json.Marshal(obj)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.expandKey(key), data, 0600)

	if err != nil {
		return err
	}

	return nil
}

func (c *CachingClient) expandKey(key string) string {
	return c.CacheDir + "/" + key + ".json"
}
