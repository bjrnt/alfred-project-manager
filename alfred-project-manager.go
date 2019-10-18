package main

import (
	"io/ioutil"
	"os"
	"path"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/fuzzy"
)

const (
	projectDirEnvVar = "PROJECT_DIRECTORY"
	maxResults       = 3
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

// projectsDir is where we look for projects
func projectsDir() string {
	dir := os.Getenv(projectDirEnvVar)
	if len(dir) == 0 {
		const err = "Please set %s before using the workflow"
		wf.Fatalf(err, projectDirEnvVar)
	}
	// Absolute paths will be honored
	if path.IsAbs(dir) {
		return dir
	}
	// Relatives path are in relation to the user's home directory
	return path.Join(os.Getenv("HOME"), dir)
}

// projectsPaths returns a list of project paths
func projectPaths() []string {
	paths := []string{}
	dir := projectsDir()
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if file.IsDir() {
			path := path.Join(dir, file.Name())
			paths = append(paths, path)
		}
	}
	return paths
}

func projects() []Project {
	if projects := TryCache(); len(projects) > 0 {
		return projects
	}

	projects := NewProjectsFromPaths(projectPaths())

	SaveCache(projects)
	return projects
}

func run() {
	query := wf.Args()[1]
	for _, project := range projects() {
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
