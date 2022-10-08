package main

import (
	"errors"
	"flag"
	"os"
	"path"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/technicallyjosh/alfred-bookmarks/browser"
)

const defaultCacheTTL = time.Hour * 1

var (
	wf               *aw.Workflow
	shouldClearCache bool
	shouldClearData  bool
	browserName      string
)

func init() {
	flag.BoolVar(&shouldClearCache, "clear-cache", false, "Clear bookmark cache")
	flag.BoolVar(&shouldClearData, "clear-data", false, "Clear all plugin data")
	flag.StringVar(&browserName, "set-browser", "", "Set browser")

	flag.Parse()
	wf = aw.New()
}

func main() {
	wf.Run(run)
}

func clearCache() {
	err := wf.ClearCache()
	if err != nil {
		wf.FatalError(err)
	}
}

func clearData() {
	clearCache()

	err := wf.ClearData()
	if err != nil {
		wf.FatalError(err)
	}
}

func run() {
	if shouldClearCache {
		clearCache()
		return
	}

	if shouldClearData {
		clearData()
		return
	}

	if browserName != "" {
		if err := setBrowserCache(); err != nil {
			if err != nil {
				wf.FatalError(err)
			}
		}

		return
	}

	query := wf.Args()[0]

	var name string
	nameBytes, err := wf.Cache.Load("browser_name")
	if err == nil {
		name = string(nameBytes)
	}

	icon := aw.Icon{
		Value: getBrowserIconPath(name),
	}

	if query == "" {
		if name != "" {
			wf.NewItem("Type keywords to search...").Icon(&icon)
		}

		wf.NewItem("Configure").
			Icon(aw.IconSettings).
			Arg("-show-config").
			Valid(true)

		wf.SendFeedback()
		return
	}

	b, err := getBrowserFromCache(name)
	if err != nil {
		wf.FatalError(err)
	}

	items, err := b.GetBookmarkItems(wf)
	if err != nil {
		wf.FatalError(err)
	}

	for i := 0; i < len(items); i++ {
		wf.NewItem(items[i].Name).
			Icon(&icon).
			Subtitle(items[i].URL).
			Arg(items[i].URL).
			Valid(true)
	}

	res := wf.Filter(query)

	if len(res) == 0 {
		wf.
			NewItem("No results found").
			Subtitle("Try another term").
			Icon(aw.IconError)
	}

	wf.SendFeedback()
}

func setBrowserCache() (err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	switch browserName {
	case "brave":
		err = wf.Cache.StoreJSON("browser_config", browser.Brave{
			Config: browser.Config{
				Directory: path.Join(homeDir, "Library/Application Support/BraveSoftware/Brave-Browser"),
				CacheName: "alfred_bookmarks:brave",
				CacheTTL:  defaultCacheTTL,
			},
		})
	case "edge":
		err = wf.Cache.StoreJSON("browser_config", browser.Edge{
			Config: browser.Config{
				Directory: path.Join(homeDir, "Library/Application Support/Microsoft Edge"),
				CacheName: "alfred_bookmarks:edge",
				CacheTTL:  defaultCacheTTL,
			},
		})
	case "firefox":
		err = wf.Cache.StoreJSON("browser_config", browser.Firefox{
			Config: browser.Config{
				Directory: path.Join(homeDir, "Library/Application Support/Firefox/Profiles"),
				CacheName: "alfred_bookmarks:firefox",
				CacheTTL:  defaultCacheTTL,
			},
		})
	default:
		return errors.New("invalid browser")
	}
	if err != nil {
		return
	}

	err = wf.Cache.Store("browser_name", []byte(browserName))
	if err != nil {
		return
	}

	return
}

func getBrowserFromCache(name string) (browser.Browser, error) {
	const cacheKey = "browser_config"
	var err error
	var b browser.Browser

	switch name {
	case "brave":
		var brave browser.Brave
		if wf.Cache.Exists(cacheKey) {
			err = wf.Cache.LoadJSON(cacheKey, &brave)
		}
		b = brave
	case "edge":
		var edge browser.Edge
		if wf.Cache.Exists(cacheKey) {
			err = wf.Cache.LoadJSON(cacheKey, &edge)
		}
		b = edge
	case "firefox":
		var firefox browser.Firefox
		if wf.Cache.Exists(cacheKey) {
			err = wf.Cache.LoadJSON(cacheKey, &firefox)
		}
		b = firefox
	default:
		return nil, errors.New("invalid browser name")
	}

	return b, err
}

func getBrowserIconPath(name string) string {
	dir := path.Join(wf.Dir(), "List Filter Images")

	switch name {
	case "brave":
		return path.Join(dir, "9efae4be72add683de8ab34e1d6f5e40c1543522.png")
	case "edge":
		return path.Join(dir, "d3ec5bd1cd0d0d91fb5650c7f6cc9a6487f4e966.png")
	case "firefox":
		return path.Join(dir, "ecdf9d0a155d7f261b4da277c4b83d9ddf6fb231.png")
	}

	return ""
}
