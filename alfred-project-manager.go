package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/fuzzy"
	"github.com/pkg/errors"
)

const (
	cacheName         = "projects.json"
	projectPathEnvVar = "PROJECT_DIRECTORY"
	maxCacheAge       = 30 * time.Minute
	maxResults        = 3
)

var (
	wf           *aw.Workflow
	fuzzyOptions []fuzzy.Option
)

func init() {
	fuzzyOptions = []fuzzy.Option{
		fuzzy.AdjacencyBonus(10.0),
		fuzzy.UnmatchedLetterPenalty(-0.5),
	}
	wf = aw.New(aw.MaxResults(maxResults), aw.SortOptions(fuzzyOptions...))
}

type Project struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func getRemote(path string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "could not retrieve git command output")
	}
	if len(out) == 0 {
		return "", errors.New("could not find repo's origin")
	}
	res := strings.TrimSpace(string(out))
	repo := out[strings.Index(res, ":")+1 : strings.LastIndex(res, ".")]
	url := fmt.Sprintf("https://github.com/%s", repo)
	return url, nil
}

func projectsPath() string {
	dir := os.Getenv(projectPathEnvVar)
	if len(dir) == 0 {
		wf.Fatalf("Please set %s before using the workflow", projectPathEnvVar)
	}
	// Absolute paths will be honored
	if path.IsAbs(dir) {
		return dir
	}
	// Relatives path are in relation to the user's home directory
	return path.Join(os.Getenv("HOME"), dir)
}

func getProjects() []*Project {
	projects := []*Project{}

	if wf.Cache.Exists(cacheName) && !wf.Cache.Expired(cacheName, maxCacheAge) {
		wf.Cache.LoadJSON(cacheName, &projects)
		return projects
	}

	projectsPath := projectsPath()
	files, _ := ioutil.ReadDir(projectsPath)
	for _, file := range files {
		if file.IsDir() {
			path := path.Join(projectsPath, file.Name())
			url, _ := getRemote(path)
			projects = append(projects, &Project{Path: path, URL: string(url)})
		}
	}

	wf.Cache.StoreJSON(cacheName, projects)
	return projects
}

func run() {
	query := wf.Args()[1]
	for _, project := range getProjects() {
		wf.NewFileItem(project.Path).
			Subtitle("Open in editor").
			Arg(project.Path).
			Var("url", project.URL).
			UID(project.Path).
			Valid(true)
	}
	wf.Filter(query)
	wf.WarnEmpty("No matching projects found", "Try a different query")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
