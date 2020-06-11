package main

import (
	"log"
	"path"
)

type Project struct {
	Path      string `json:"path"`
	Workspace string `json:"workspace"`
	URL       string `json:"url"`
}

func NewProjectFromPath(path string, workspace string) Project {
	url, err := RepoURL(path)
	if err != nil {
		log.Println(err)
		url = ""
	}
	return Project{
		Path:      path,
		Workspace: workspace,
		URL:       url,
	}
}

func (p *Project) Name() string {
	return path.Join(p.Workspace, path.Base(p.Path))
}
