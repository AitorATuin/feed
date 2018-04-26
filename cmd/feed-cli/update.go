package main

import (
	"errors"

	"github.com/AitorATuin/feed/cache"
)

func UpdateCommand(args []string) error {
	db, err := cache.GetCacheDriver(FEEDS_CACHE)
	switch len(args) {
	case 0:
		// Update all the feeds?
		feedsIndex, err := cache.NewFeedsIndex(db)
		for feed := range feedsIndex.FeedsIndexList
		return nil
	case 1:
		// Update feed
	}
	return errors.New("Wrong update command")
}
