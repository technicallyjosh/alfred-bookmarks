package browser

import (
	aw "github.com/deanishe/awgo"
)

type Edge struct {
	Config Config
}

// GetBookmarkItems returns flattened items to search through from bookmarks.
func (b Edge) GetBookmarkItems(wf *aw.Workflow) ([]Item, error) {
	return getChromiumBookmarkItems(b.Config, wf)
}

// GetProfiles returns profiles for the browser.
func (b Edge) GetProfiles() ([]string, error) {
	return nil, nil
}
