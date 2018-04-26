package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/AitorATuin/feed/cache"
)

func listFeeds(feedsIndex *cache.FeedsIndex) error {
	for i, feed := range feedsIndex.FeedsIndexList {
		fmt.Printf("[%d] %v\n", i, feed)
	}
	return nil
}

func listFeedItems(feedsIndex *cache.FeedsIndex, feedNumber int) error {
	feed, err := feedsIndex.GetFeed(feedNumber)
	if err != nil {
		return err
	}
	for n, item := range feed.Items {
		fmt.Printf("[%d] %v\n", n, item.Title)
	}
	return nil
}

func listFeedItem(feedsIndex *cache.FeedsIndex, feedNumber int, itemNumber int) error {
	item, err := feedsIndex.GetItem(feedNumber, itemNumber)
	if err != nil {
		return err
	}

	fmt.Printf("== %v ==\n\n", item.Title)
	fmt.Printf("%v\n", item.Content)

	return nil
}

func GetCommand(args []string) error {
	db, err := cache.GetCacheDriver(FEEDS_CACHE)
	if err != nil {
		return nil, err
	}

	feedsIndex, err := cache.NewFeedsIndex(db)
	if err != nil {
		return err
	}

	switch len(args) {
	case 0:
		// List feeds
		return listFeeds(feedsIndex)
	case 1:
		// List items in feed
		if feedNumber, err := strconv.Atoi(args[0]); err != nil {
			return err
		} else {
			return listFeedItems(feedsIndex, feedNumber)
		}
	case 2:
		// List item in feed
		feedNumber, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		itemNumber, err := strconv.Atoi(args[1])
		if err != nil {
			return nil
		}
		return listFeedItem(feedsIndex, feedNumber, itemNumber)
	}

	return errors.New("Unknown get command")
}
