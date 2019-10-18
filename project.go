package main

import (
	"log"
)

type Projects = []Project

func NewProjectsFromPaths(paths []string) []Project {
	projects := []Project{}
	for _, path := range paths {
		projects = append(projects, NewProjectFromPath(path))
	}
	return projects
}

type Project struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func NewProjectFromPath(path string) Project {
	url, err := RepoURL(path)
	if err != nil {
		log.Println(err)
		url = ""
	}
	return Project{
		Path: path,
		URL:  url,
	}
}
