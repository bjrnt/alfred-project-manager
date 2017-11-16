// Copyright (c) Bjorn Tegelund 2017 <b.tegelund@gmail.com>
// MIT Licence. See http://opensource.org/licenses/MIT

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/deanishe/awgo"
	"github.com/deanishe/awgo/fuzzy"
	"github.com/docopt/docopt-go"
)

var usage = `alfred-project-manager

  Usage:
    alfred-project-manager search <query>
    alfred-project-manager -h | --help

  Options:
    -h --help     Show this screen.`

var (
	cacheName    = "projects.json"
	wf           *aw.Workflow
	fuzzyOptions []fuzzy.Option
	maxCacheAge  = 10 * time.Minute
	maxResults   = 3
)

func init() {
	fuzzyOptions = []fuzzy.Option{
		fuzzy.AdjacencyBonus(10.0),
		fuzzy.UnmatchedLetterPenalty(-0.5),
	}
	wf = aw.New(aw.MaxResults(maxResults), aw.SortOptions(fuzzyOptions...))
}

func getProjectDirectories() []string {
	var projectsDirectories []string
	if wf.Cache.Exists(cacheName) && !wf.Cache.Expired(cacheName, maxCacheAge) {
		wf.Cache.LoadJSON(cacheName, &projectsDirectories)
		return projectsDirectories
	}

	projectsPath := filepath.Join(os.Getenv("HOME"), os.Getenv("PROJECT_DIRECTORY"))
	files, _ := ioutil.ReadDir(projectsPath)

	projectsDirectories = make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			projectsDirectories = append(projectsDirectories, filepath.Join(projectsPath, file.Name()))
		}
	}
	wf.Cache.StoreJSON(cacheName, projectsDirectories)
	return projectsDirectories
}

func run() {
	arguments, _ := docopt.Parse(usage, wf.Args(), true, wf.Version(), false, true)

	query := arguments["<query>"].(string)

	for _, project := range getProjectDirectories() {
		wf.NewFileItem(project).
			Subtitle("Open in VSCode").
			Arg(project).
			UID(project).
			Valid(true)
	}

	wf.Filter(query)

	wf.WarnEmpty("No matching projects found", "Try a different query")

	wf.SendFeedback()
}

func main() {
	aw.Run(run)
}
