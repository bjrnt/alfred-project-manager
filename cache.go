package main

import (
	"log"
	"os"
	"syscall"
	"time"
)

const (
	cacheName   = "projects.json"
	maxCacheAge = 6 * time.Hour
)

type cache struct {
	Params            Params
	ProjDirModifiedAt time.Time
	Projects          []Project
}

// TryCache reads in any previously cached projects. The cache is invalidated if the params given
// params do not matched the saved ones.
func TryCache(params *Params) []Project {
	if wf.Cache.Expired(cacheName, maxCacheAge) {
		log.Printf("Cache does not exist or has expired -- skipping cache")
		return nil
	}
	cache := cache{}
	// swallow errors for invalid caches to overwrite later
	_ = wf.Cache.LoadJSON(cacheName, &cache)
	// ensure that the cache and the current search have the same parameters
	if !params.Equal(cache.Params) {
		log.Printf("Search params do not match cached -- skipping cache")
		return nil
	}
	// skip the cache if there have been changes to the project dir
	if getModifiedAt(params.ProjectsPath()).After(cache.ProjDirModifiedAt) {
		log.Printf("Projects have been added/removed since last caching -- skipping cache")
		return nil
	}
	return cache.Projects
}

// SaveCache saves the given projects to the cache file.
func SaveCache(params *Params, projects []Project) {
	err := wf.Cache.StoreJSON(cacheName, cache{*params, getModifiedAt(params.ProjectsPath()), projects})
	if err != nil {
		log.Printf("could not save cache: %s", err.Error())
	}
}

func getModifiedAt(path string) time.Time {
	dir, err := os.Stat(path)
	if err != nil {
		log.Printf("could not verify modified time of projects dir: %s", err.Error())
		return time.Time{}
	}
	timespec := dir.Sys().(*syscall.Stat_t).Ctimespec
	return time.Unix(timespec.Unix())
}
