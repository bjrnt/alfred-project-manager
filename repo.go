package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func IsGitRepo(p string) bool {
	if _, err := os.Stat(path.Join(p, ".git")); os.IsNotExist(err) {
		return false
	}
	return true
}

// RepoURL gets a URL for the git repo at the given path that can be opened in a browser. Supports
// github and http/https origins.
func RepoURL(path string) (string, error) {
	origin, err := gitOriginAt(path)
	if err != nil {
		return "", errors.Wrapf(err, "could not find origin for project at path %s", path)
	}
	return urlForOrigin(origin), nil
}

func gitOriginAt(path string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "git command failed")
	}
	if len(out) == 0 {
		return "", errors.New("no output from git command")
	}
	return strings.TrimSpace(string(out)), nil
}

func urlForOrigin(origin string) string {
	if strings.HasPrefix(origin, "git@github.com") {
		// format: git@github.com:user/repo.git
		end := len(origin) - 1
		// sometimes the .git can be missing
		indexOfLastPeriod := strings.LastIndex(origin, ".")
		if indexOfLastPeriod != -1 {
			end = indexOfLastPeriod
		}
		repo := origin[strings.Index(origin, ":")+1 : end]
		url := fmt.Sprintf("https://github.com/%s", repo)
		return url
	} else if strings.HasPrefix(origin, "http") {
		// assume its some kind of URL we can open directly
		return origin
	}
	return ""
}
