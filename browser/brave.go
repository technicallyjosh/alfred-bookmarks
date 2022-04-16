package browser

import (
	aw "github.com/deanishe/awgo"
)

type Brave struct {
	Config Config
}

// GetBookmarkItems returns flattened items to search through from bookmarks.
func (b Brave) GetBookmarkItems(wf *aw.Workflow) ([]Item, error) {
	return getChromiumBookmarkItems(b.Config, wf)
}

// GetProfiles returns profiles for the browser.
func (b Brave) GetProfiles() ([]string, error) {
	return nil, nil
}
