package feed

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func PrintFeed() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("http://feeds.twit.tv/twit.xml")
	fmt.Println(feed.Title)
}
