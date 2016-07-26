package github

import (
	"os"
	"path/filepath"
)

func alfredGithubDir(subdir string) string {
	home := os.Getenv("ALFRED_GITHUB_HOME")
	if home == "" {
		home = filepath.Join(os.Getenv("HOME"), ".gh-prj")
	}
	return filepath.Join(home, subdir)
}
