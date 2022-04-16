package browser

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

	aw "github.com/deanishe/awgo"
)

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

func getChromiumBookmarkItems(config Config, wf *aw.Workflow) (allItems []Item, err error) {
	cacheName := config.CacheName
	// if it's not expired, we'll return it from cache
	if !wf.Cache.Expired(cacheName, config.CacheTTL) {
		wf.Var("from_cache", "true")

		var items []Item
		if err := wf.Cache.LoadJSON(cacheName, &items); err != nil {
			return nil, err
		}

		return items, nil
	}

	wf.Var("from_cache", "false")

	profiles, err := getChromiumProfiles(config.Directory)
	if err != nil {
		return
	}

	for _, p := range profiles {
		filePath := path.Join(config.Directory, p, "Bookmarks")

		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}

		var bookmarks Bookmarks
		err = json.Unmarshal(fileBytes, &bookmarks)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, bookmarks.Flatten()...)
	}

	if err := wf.Cache.StoreJSON(cacheName, allItems); err != nil {
		return nil, err
	}

	return allItems, nil
}

func getChromiumProfiles(dir string) ([]string, error) {
	profiles := []string{"Default"}

	fileInfos, err := ioutil.ReadDir(dir)
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
