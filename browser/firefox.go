package browser

import (
	"database/sql"
	"os"
	"path/filepath"

	aw "github.com/deanishe/awgo"
	_ "github.com/mattn/go-sqlite3"
)

type Firefox struct {
	Config Config
}

func (b Firefox) GetBookmarkItems(wf *aw.Workflow) ([]Item, error) {
	if !wf.Cache.Expired(b.Config.CacheName, b.Config.CacheTTL) {
		wf.Var("from_cache", "true")

		var items []Item
		if err := wf.Cache.LoadJSON(b.Config.CacheName, &items); err != nil {
			return nil, err
		}

		return items, nil
	}

	var items []Item
	err := filepath.Walk(b.Config.Directory, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() { // profile dir
			return nil
		}

		// This will run on each profile's found place, ultimately merging them all.
		if !info.IsDir() && info.Name() == "places.sqlite" {
			items, err = b.getItemsFromDatabase(currentPath, items)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := wf.Cache.StoreJSON(b.Config.CacheName, items); err != nil {
		return nil, err
	}

	return items, nil
}

func (b Firefox) getItemsFromDatabase(dbFile string, items []Item) ([]Item, error) {
	oldFile, err := os.ReadFile(dbFile)
	if err != nil {
		return items, err
	}

	err = os.WriteFile("db.sqlite", oldFile, 0644)
	if err != nil {
		return items, err
	}

	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		return items, err
	}

	rows, err := db.Query(`
	SELECT b.guid, b.title, p.url
	FROM moz_bookmarks b 
	JOIN moz_places p ON p.id = b.fk
	WHERE type = 1`)
	if err != nil {
		return items, err
	}

	for rows.Next() {
		var guid, url string
		var title sql.NullString

		err = rows.Scan(&guid, &title, &url)
		if err != nil {
			return items, err
		}

		items = append(items, Item{
			GUID: guid,
			Name: title.String,
			URL:  url,
		})
	}

	_ = rows.Close()
	_ = db.Close()

	return items, nil
}

func (b Firefox) GetProfiles() ([]string, error) {
	return nil, nil
}
