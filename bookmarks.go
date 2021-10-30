package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

// Bookmarks represents a "Bookmarks" file.
type Bookmarks struct {
	Roots   Roots `json:"roots"`
	Version int   `json:"version"`
}

// Roots represents the "roots" node in Bookmarks.
type Roots struct {
	BookmarkBar Node `json:"bookmark_bar"`
	Other       Node `json:"other"`
	Synced      Node `json:"synced"`
}

// Node represents a folder or url bookmark.
type Node struct {
	GUID     string `json:"guid"`
	Name     string `json:"name"`
	Type     string `json:"type"` // url or folder
	URL      string `json:"url"`
	Children []Node `json:"children"`
}

// Item represents an actionable item to add to the list in alfred to be filtered on.
type Item struct {
	Name string
	URL  string
	GUID string
}

// Flatten flattens all nodes and their children into a filterable list.
func (b Bookmarks) Flatten() []Item {
	topLevel := []Node{
		b.Roots.BookmarkBar,
		b.Roots.Other,
		b.Roots.Synced,
	}

	items := flatten(topLevel)

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items
}

// flatten recursively appends all url type nodes to an item collection. func flatten(nodes []Node) (items []Item) {
func flatten(nodes []Node) (items []Item) {
	for _, node := range nodes {
		if node.Type == "url" {
			items = append(items, Item{
				Name: node.Name,
				URL:  node.URL,
				GUID: node.GUID,
			})
			continue
		}

		items = append(items, flatten(node.Children)...)
	}

	return
}

// getProfiles returns a list of profiles found in the Brave directory.
func getProfiles(baseDir string) ([]string, error) {
	profiles := []string{"Default"}

	fileInfos, err := ioutil.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() && strings.HasPrefix(fileInfo.Name(), "Profile ") {
			profiles = append(profiles, fileInfo.Name())
		}
	}

	return profiles, nil
}

// getBookmarkItems returns items based on bookmarks. Right now it returns all bookmarks for all
// profiles.
func getBookmarkItems(cacheTTL time.Duration) ([]Item, error) {
	var allItems []Item
	// if it's not expired, we'll return it from cache
	if !wf.Cache.Expired(cacheName, cacheTTL) {
		wf.Var("from_cache", "true")

		var items []Item
		if err := wf.Cache.LoadJSON(cacheName, &items); err != nil {
			return nil, err
		}

		return items, nil
	}

	wf.Var("from_cache", "false")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := path.Join(homeDir, braveDir)
	profiles, err := getProfiles(baseDir)
	if err != nil {
		return nil, err
	}

	for _, p := range profiles {
		filePath := path.Join(baseDir, p, "Bookmarks")

		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var bookmarks Bookmarks
		if err := json.Unmarshal(bytes, &bookmarks); err != nil {
			return nil, err
		}

		allItems = append(allItems, bookmarks.Flatten()...)
	}

	if err := wf.Cache.StoreJSON(cacheName, allItems); err != nil {
		return nil, err
	}

	return allItems, nil
}
