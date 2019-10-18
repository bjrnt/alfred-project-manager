package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// RepoURL gets a URL for the git repo at the given path that can be opened in a browser. Supports
// github and http/https origins.
func RepoURL(path string) (string, error) {
	origin, err := gitOriginAt(path)
	if err != nil {
		return "", errors.Wrapf(err, "could not find origin for project at path %s", path)
	}
	url, err := urlForOrigin(origin)
	return url, errors.Wrapf(err, "could not format URL for project at path %s", path)
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

func urlForOrigin(origin string) (string, error) {
	if strings.HasPrefix(origin, "git@github.com") {
		// format: git@github.com:user/repo.git
		repo := origin[strings.Index(origin, ":")+1 : strings.LastIndex(origin, ".")]
		url := fmt.Sprintf("https://github.com/%s", repo)
		return url, nil
	} else if strings.HasPrefix(origin, "http") {
		// assume its some kind of URL we can open directly
		return origin, nil
	}

	return "", errors.Errorf("unrecognized origin format: %s", origin)
}
