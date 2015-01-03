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

func (c *CachingClient) fetchCache(key string, objPtr interface{}, fetchFn func() (interface{}, error)) (interface{}, error) {
	err := c.readCache(key, objPtr)
	if err != nil {
		log.Printf("Cache miss: %s, or error reading cache", key)
		objPtr, err = fetchFn()
		if err != nil {
			log.Fatalf("Failed to fetch %s: %s", key, err)
		}
		err = c.writeCache(key, objPtr)

		if err != nil {
			log.Printf("Failed to write cache %s: %s", key, err)
		}
	}
	return objPtr, nil
}

func (c *CachingClient) readCache(key string, objPtr interface{}) error {
	data, err := ioutil.ReadFile(c.expandKey(key))

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, objPtr)

	if err != nil {
		return err
	}

	return nil
}

func (c *CachingClient) writeCache(key string, objPtr interface{}) error {
	err := os.MkdirAll(c.CacheDir, 0700)

	if err != nil {
		return err
	}

	data, err := json.Marshal(objPtr)

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
