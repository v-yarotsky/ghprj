package github

import (
	"os"
)

func AccessToken() string {
	return os.Getenv("GH_PRJ_GITHUB_TOKEN")
}

func ExpireAccessToken() {
}

func obtainAccessToken() {
}
