package main

import (
	"io/ioutil"
	"path"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/fuzzy"
)

const (
	gitDirName = ".git"
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

func scanDirAtPath(basePath string, workspace string, maxDepth uint, requireDotGit bool) []Project {
	projects := []Project{}
	fullPath := path.Join(basePath, workspace)
	files, _ := ioutil.ReadDir(fullPath)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		p := path.Join(fullPath, file.Name())
		if requireDotGit && IsGitRepo(p) {
			// include current and stop
			projects = append(projects, NewProjectFromPath(p, workspace))
			continue

		} else if maxDepth > 0 {
			// look at children and maybe add current
			projects = append(projects, scanDirAtPath(basePath, path.Join(workspace, file.Name()), maxDepth-1, requireDotGit)...)
		}

		if !requireDotGit {
			// include current
			projects = append(projects, NewProjectFromPath(p, workspace))
		}
	}
	return projects
}

func scanProjects(params *Params) []Project {
	return scanDirAtPath(params.ProjectsPath(), "", params.MaxProjectDepth, params.RequireDotGit)
}

func projects(params *Params) []Project {
	if projects := TryCache(params); len(projects) > 0 {
		return projects
	}

	projects := scanProjects(params)

	SaveCache(params, projects)
	return projects
}

func run() {
	query := wf.Args()[1]
	params, err := NewParamsFromEnv()
	if err != nil {
		wf.Fatalf(err.Error())
	}

	for _, project := range projects(params) {
		wf.NewFileItem(project.Path).
			Title(project.Name()).
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
