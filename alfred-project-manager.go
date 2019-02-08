// Copyright (c) Bjorn Tegelund 2017 <b.tegelund@gmail.com>
// MIT Licence. See http://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/fuzzy"
	docopt "github.com/docopt/docopt-go"
)

var usage = `alfred-project-manager

  Usage:
    alfred-project-manager search <query>
    alfred-project-manager -h | --help

  Options:
    -h --help     Show this screen.`

var (
	cacheName    = "projects1234.json"
	wf           *aw.Workflow
	fuzzyOptions []fuzzy.Option
	maxCacheAge  = 30 * time.Minute
	maxResults   = 3
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

func getProjects() []*Project {
	projects := []*Project{}
	if wf.Cache.Exists(cacheName) && !wf.Cache.Expired(cacheName, maxCacheAge) {
		wf.Cache.LoadJSON(cacheName, &projects)
		return projects
	}

	projectsPath := filepath.Join(os.Getenv("HOME"), os.Getenv("PROJECT_DIRECTORY"))
	files, _ := ioutil.ReadDir(projectsPath)

	for _, file := range files {
		if file.IsDir() {
			path := filepath.Join(projectsPath, file.Name())
			url, _ := getRemote(path)
			projects = append(projects, &Project{Path: path, URL: string(url)})
		}
	}
	wf.Cache.StoreJSON(cacheName, projects)
	return projects
}

func run() {
	arguments, _ := docopt.Parse(usage, wf.Args(), true, wf.Version(), false, true)

	query := arguments["<query>"].(string)

	for _, project := range getProjects() {
		log.Println(project.Path)
		log.Println(project.URL)
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
