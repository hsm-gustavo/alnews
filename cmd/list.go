package cmd

import (
	"fmt"
	"time"

	"github.com/hsm-gustavo/alnews/internal/cache"
	"github.com/hsm-gustavo/alnews/internal/fetch"
	"github.com/hsm-gustavo/alnews/internal/model"
	"github.com/hsm-gustavo/alnews/internal/platform"
	"github.com/hsm-gustavo/alnews/internal/render"
	"github.com/hsm-gustavo/alnews/internal/search"
)

func listCmd(limit uint8, refetch bool, searchParam string, openIndex int8, inspectIndex int8) error {
	if err := platform.IsUnsupportedPlatform(); err != nil {
		return err
	}

	cm, err := cache.New("alnews", "alnews.json")
	if err != nil {
		return err
	}

	var feed model.RSS

	if cm.Exists() && cm.IsFresh(12*time.Hour) && !refetch {
		fmt.Println("cache hit")

		if err := cm.ReadJSON(&feed); err != nil {
			return err
		}
	} else {
		feed, err = fetch.FetchRSS("") // if empty uses default
		if err != nil {
			return err
		}
		cm.WriteJSON(feed)
	}

	items := feed.Channel.Items

	if searchParam != "" {
		items = search.Filter(items, searchParam)
	}

	if len(items) > int(limit) {
		items = items[:limit]
	}

	render.List(items)

	if openIndex >= 0 {
		render.Open(items[openIndex].Link)
	}

	if inspectIndex >= 0 {
		render.Inspect(items[inspectIndex])
	}

	return nil
}