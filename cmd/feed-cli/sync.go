package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AitorATuin/feed/cache"
	"github.com/mmcdole/gofeed"
)

func readFeedFiles() (cache.Feeds, error) {
	feedUrls := make([]string, 0, 50)
	fn := func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".feed") {
			fmt.Printf("Reading feed %v\n", path)
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			defer func() {
				file.Close()
			}()
			if err != nil {
				fmt.Printf("Error reading feeds file %v\n", path)
			} else {
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					feedUrls = append(feedUrls, scanner.Text())
				}
			}
		}
		return nil
	}

	if _, err := os.Stat(FEEDS_DIR); err != nil {
		return nil, fmt.Errorf("Feeds directory %v does not exist", FEEDS_DIR)
	}

	filepath.Walk(FEEDS_DIR, fn)

	parser := gofeed.NewParser()
	feeds := make(cache.Feeds, len(feedUrls))
	for i, feedUrl := range feedUrls {
		feed, err := parser.ParseURL(feedUrl)
		if err != nil {
			fmt.Printf("Error parting feed url for %v: %v\n", feedUrl, err.Error())
			continue
		}
		feeds[i] = feed
		fmt.Printf("Read feed %v\n", feed.Title)
	}
	return feeds, nil
}

func SyncCommand(args []string) error {
	switch len(args) {
	case 0:
		readFeeds, err := readFeedFiles()
		if err != nil {
			return err
		}
		feedsIndex, err := readFeeds.Cache(FEEDS_CACHE)
		if err != nil {
			return err
		}
		fmt.Printf("Defined %d feeds, saved %d feeds\n", len(readFeeds), len(feedsIndex.FeedsIndexList))
		return nil
	}

	return errors.New("Wrong sync command")
}
