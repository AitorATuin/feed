// Package main provides feed
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func getFeedHome() string {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		fmt.Printf("Unable to find home")
		os.Exit(1)

	}

	return filepath.Join(home, ".config", "feed")
}

var (
	FEED_HOME   = getFeedHome()
	FEED_RC     = filepath.Join(FEED_HOME, "feed.rc")
	FEEDS_DIR   = filepath.Join(FEED_HOME, "feeds")
	FEEDS_CACHE = filepath.Join(FEED_HOME, "cache")
)

func Usage() {
	usage := `Usage: feed-cli command
	commands:

	get     : List all the feeds defined
	get n   : List all the items for feed n
	get n m : List item m for feed n
	sync    : Sync local defined feeds. This will delete the ones cached
	update  : Updates all the feeds

`
	fmt.Printf(usage)
}

func doCommand(command string, args []string) error {
	switch command {
	case "get":
		return GetCommand(args)
	case "update":
		return UpdateCommand(args)
	case "sync":
		return SyncCommand(args)
	default:
		Usage()
		os.Exit(0)
	}
	return nil
}

func main() {
	flag.Parse()
	if len(os.Args) == 1 {
		Usage()
		os.Exit(1)
	}
	command, args := os.Args[1], os.Args[2:]
	if err := doCommand(command, args); err != nil {
		fmt.Printf("Error found: %v\n", err.Error())
		os.Exit(2)
	}
}
