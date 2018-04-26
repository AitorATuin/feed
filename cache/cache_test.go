package cache

import (
	"errors"
	"testing"

	"github.com/franela/goblin"
)

type FakePersistence struct{}

func (*FakePersistence) Write(collection, resource string, v interface{}) error {
	return nil
}

func (*FakePersistence) Read(collection, resource string, v interface{}) error {
	v.FeedsIndexList = []string{"feed one", "feed two"}
	return nil
}

func Test(t *testing.T) {
	g := goblin.Goblin(t)

	var feedsIndex *FeedsIndex

	g.Describe("FeedsIndex", func() {
		g.Before(func() {
			feedsIndex = new(FeedsIndex)
			feedsIndex.FeedsIndexList = []string{"feed one", "feed two"}
		})

		g.It("listFeed out of index", func() {
			err := listFeedItems(feedsIndex, 3)
			g.Assert(err).Eql(errors.New("feed not found"))
		})
		g.It("listItem out of index", func() {
			err := listFeedItems(feedsIndex, 1)
			g.Assert()
		})
	})
}
