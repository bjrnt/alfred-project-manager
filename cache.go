package main

import "time"

const (
	cacheName   = "projects.json"
	maxCacheAge = 30 * time.Minute
)

// TryCache reads in any previously cached projects. May give incorrect results right after changing
// the PROJECTS_DIR variable.
func TryCache() Projects {
	if !wf.Cache.Exists(cacheName) {
		return nil
	}
	if wf.Cache.Expired(cacheName, maxCacheAge) {
		return nil
	}
	projects := Projects{}
	// swallow errors for invalid caches to overwrite later
	_ = wf.Cache.LoadJSON(cacheName, &projects)
	return projects
}

// SaveCache saves the given projects to the cache file.
func SaveCache(projects Projects) {
	wf.Cache.StoreJSON(cacheName, projects)
}
