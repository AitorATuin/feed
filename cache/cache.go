package cache

import (
	"errors"
	"fmt"
	"sort"

	"github.com/mmcdole/gofeed"
	scribble "github.com/nanobox-io/golang-scribble"
)

var (
	FEEDS_CACHE_COL = "feeds"
	cache_drives    = make(map[string]*scribble.Driver)
)

type Cacheable interface {
	Cache() (int, error)
}

type Persistence interface {
	Read(collection, resource string, v interface{}) error
	Write(collection, resource string, v interface{}) error
}

type Feed gofeed.Feed

type Feeds []*gofeed.Feed

type FeedsIndex struct {
	FeedsIndexMap  map[string]string `json:"indexMap"`
	FeedsIndexList sort.StringSlice  `json:"indexList"`
	cache          Persistence
}

var CACHE_DRIVER *scribble.Driver

func GetCacheDriver(dir string) (*scribble.Driver, error) {
	if cache_drives[dir] == nil {
		db, err := scribble.New(dir, nil)
		if err != nil {
			return nil, err
		}
		cache_drives[dir] = db
	}
	return cache_drives[dir], nil
}

func (feed *Feed) Cache(dir string) (int, error) {
	db, err := GetCacheDriver(dir)
	if err != nil {
		return 0, err
	}
	if err := db.Write(FEEDS_CACHE_COL, feed.Title, nil); err != nil {
		return 0, err
	}
	return 1, nil
}

func (feeds Feeds) Cache(dir string) (*FeedsIndex, error) {
	db, err := GetCacheDriver(dir)
	if err != nil {
		return nil, err
	}
	savedFeeds := 0
	feedsIndex := new(FeedsIndex)
	feedsIndex.FeedsIndexList = make([]string, len(feeds), len(feeds))
	feedsIndex.FeedsIndexMap = make(map[string]string)
	for _, feed := range feeds {
		if err := db.Write(FEEDS_CACHE_COL, feed.Title, feed); err != nil {
			fmt.Printf("Error saving feed in disk: %v\n", err.Error())
			continue
		}
		feedsIndex.FeedsIndexMap[feed.Title] = feed.Title
		feedsIndex.FeedsIndexList[savedFeeds] = feed.Title
		savedFeeds++
	}
	if err := db.Write(FEEDS_CACHE_COL, "feeds_index", feedsIndex); err != nil {
		fmt.Printf("Error saving feed index in disk: %v\n", err.Error())
		return nil, err
	}
	return feedsIndex, nil
}

func loadFeed(cache Persistence, feed_name string) (*gofeed.Feed, error) {
	feed := &gofeed.Feed{}
	if err := cache.Read(FEEDS_CACHE_COL, feed_name, feed); err != nil {
		return nil, err
	}

	return feed, nil
}

func NewFeedsIndex(cache Persistence) (*FeedsIndex, error) {
	feedsIndex := &FeedsIndex{}
	feedsIndex.cache = cache
	if err := feedsIndex.cache.Read(FEEDS_CACHE_COL, "feeds_index", feedsIndex); err != nil {
		return nil, err
	}
	feedsIndex.FeedsIndexList.Sort()
	return feedsIndex, nil
}

func (feedsIndex *FeedsIndex) GetFeed(feedNumber int) (*gofeed.Feed, error) {
	if feedNumber >= len(feedsIndex.FeedsIndexList) {
		return nil, errors.New("feed not found")
	}
	feedName := feedsIndex.FeedsIndexList[feedNumber]

	feed, err := loadFeed(feedsIndex.cache, feedName)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (feedsIndex *FeedsIndex) GetItem(feedNumber int, itemNumber int) (*gofeed.Item, error) {
	feed, err := feedsIndex.GetFeed(feedNumber)
	if err != nil {
		return nil, err
	}

	if feedNumber >= len(feed.Items) {
		return nil, errors.New("item not found")
	}

	return feed.Items[itemNumber], nil
}
