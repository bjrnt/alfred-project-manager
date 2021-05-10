package main

import (
	"log"
	"time"
)

const (
	cacheName   = "projects.json"
	maxCacheAge = 30 * time.Minute
)

type cache struct {
	Params   Params
	Projects []Project
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
	return cache.Projects
}

// SaveCache saves the given projects to the cache file.
func SaveCache(params *Params, projects []Project) {
	err := wf.Cache.StoreJSON(cacheName, cache{*params, projects})
	log.Printf("Saved cache with result: %s", err.Error())
}
