package browser

import (
	"sort"
	"time"

	aw "github.com/deanishe/awgo"
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

type Browser interface {
	GetBookmarkItems(*aw.Workflow) ([]Item, error)
	GetProfiles() ([]string, error)
}

// Config represents properties for a browser.
type Config struct {
	Directory string
	CacheName string
	CacheTTL  time.Duration
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
