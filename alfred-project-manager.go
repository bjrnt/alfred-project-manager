// Copyright (c) Bjorn Tegelund 2017 <b.tegelund@gmail.com>
// MIT Licence. See http://opensource.org/licenses/MIT

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
)

const cacheName = "projects.json"
const projectPathEnvVar = "PROJECT_DIRECTORY"
const maxCacheAge = 30 * time.Minute
const maxResults = 3

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

func getRemote(path string) (repo string, err error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = path
	out, err := cmd.Output()
	if len(out) != 0 {
		out := strings.TrimSpace(string(out))
		remote := out[strings.Index(out, ":")+1 : strings.LastIndex(out, ".")]
		repo = fmt.Sprintf("https://github.com/%s", remote)
	}
	return
}

func projectsPath() string {
	dir := os.Getenv(projectPathEnvVar)
	if len(dir) == 0 {
		wf.Fatalf("Please set %s before using the workflow", projectPathEnvVar)
	}
	if path.IsAbs(dir) {
		return dir
	}
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
			Subtitle("Open in VSCode").
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
