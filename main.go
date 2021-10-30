package main

import (
	"flag"
	"fmt"
	"time"

	aw "github.com/deanishe/awgo"
)

const (
	cacheName = "brave_bookmarks"
	braveDir  = "Library/Application Support/BraveSoftware/Brave-Browser"
)

var (
	wf         *aw.Workflow
	clearCache bool
)

func init() {
	flag.BoolVar(&clearCache, "clear-cache", false, "Clear bookmark cache")

	wf = aw.New()
}

func run() {
	wf.Args()
	flag.Parse()

	query := wf.Args()[0]

	if clearCache {
		err := wf.ClearCache()
		if err != nil {
			wf.FatalError(err)
		}
		fmt.Println("cleared cache...")
		return
	}

	if query == "" {
		wf.NewItem("Clear cache").
			Subtitle("Manually clear the cache of bookmarks. This is done every 1 hour").
			Arg("-clear-cache").
			Icon(aw.IconTrash).
			Valid(true)
	} else {
		items, err := getBookmarkItems(1 * time.Hour)
		if err != nil {
			wf.FatalError(err)
		}

		for _, item := range items {
			wf.NewItem(item.Name).
				Subtitle(item.URL).
				Arg(item.URL).
				UID(item.GUID).
				Valid(true)
		}

		res := wf.Filter(query)

		if len(res) == 0 {
			wf.NewWarningItem("No results found", "Try another term")
		}
	}

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
