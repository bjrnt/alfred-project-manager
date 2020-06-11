package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

const (
	projectDirEnvVar   = "PROJECT_DIRECTORY"
	reqDotGitEnvVar    = "REQUIRE_DOTGIT"
	maxProjDepthEnvVar = "MAX_PROJECT_DEPTH"
	maxResults         = 5
)

type Params struct {
	ProjectsDir     string
	RequireDotGit   bool
	MaxProjectDepth uint
	MaxResults      uint
}

func NewParamsFromEnv() (*Params, error) {
	projDir := os.Getenv(projectDirEnvVar)
	if len(projDir) == 0 {
		return nil, fmt.Errorf("Please set %s before using the workflow", projectDirEnvVar)
	}

	reqDotGit, err := strconv.ParseBool(os.Getenv(reqDotGitEnvVar))
	if err != nil {
		return nil, fmt.Errorf("could not parse %s: %w", reqDotGitEnvVar, err)
	}

	maxProjDepth, err := strconv.Atoi(os.Getenv(maxProjDepthEnvVar))
	if err != nil {
		return nil, fmt.Errorf("could not parse %s: %w", maxProjDepthEnvVar, err)
	}

	return &Params{
		ProjectsDir:     projDir,
		RequireDotGit:   reqDotGit,
		MaxProjectDepth: uint(maxProjDepth),
		MaxResults:      maxResults,
	}, nil
}

func (p *Params) Equal(p2 Params) bool {
	return p.RequireDotGit == p2.RequireDotGit &&
		p.MaxProjectDepth == p2.MaxProjectDepth &&
		p.MaxResults == p2.MaxResults &&
		p.ProjectsDir == p2.ProjectsDir
}

func (p *Params) ProjectsPath() string {
	// absolute paths with be honored
	if path.IsAbs(p.ProjectsDir) {
		return p.ProjectsDir
	}
	// Relative paths are in relation to the user's home directory
	return path.Join(os.Getenv("HOME"), p.ProjectsDir)
}
