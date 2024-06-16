package feed

import (
	"log/slog"
	"regexp"
	"sort"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
)

var linkRegex = regexp.MustCompile("<a href=\"(.*)\">(.*)</a>")

func NewFetcher(logger *slog.Logger) *Fetcher {
	return &Fetcher{
		parser: gofeed.NewParser(),
		logger: logger,
	}
}

type Fetcher struct {
	parser *gofeed.Parser
	logger *slog.Logger
}

func (f *Fetcher) Fetch(feed *Feed) ([]Item, error) {
	res, err := f.parser.ParseURL(feed.URL)
	if err != nil {
		return nil, err
	}
	sort.Sort(res)

	var items []Item
	if len(res.Items) > 0 {
		p := bluemonday.StrictPolicy()
		items = make([]Item, 0, len(res.Items))
		for _, it := range res.Items {
			f.logger.Debug("original item", "item", it)

			content := it.Description
			if content == "" {
				content = it.Content
			}

			var extraLinks []Link

			// Note: feed item with just hyperlink
			if match := linkRegex.FindStringSubmatch(content); len(match) == 3 {
				content = linkRegex.ReplaceAllString(content, "")
				extraLinks = append(extraLinks, Link{
					URL:   match[1],
					Value: match[2],
				})
			}

			items = append(items, Item{
				Title:       it.Title,
				Link:        it.Link,
				Content:     p.Sanitize(content),
				ExtraLinks:  extraLinks,
				PublishedAt: *it.PublishedParsed,
			})
		}
	}

	return items, nil
}

type Feed struct {
	ID   int64
	Name string
	URL  string
}

type Item struct {
	Title       string
	Link        string
	Content     string
	ExtraLinks  []Link
	PublishedAt time.Time
}

type Link struct {
	Value string
	URL   string
}

type Subscription struct {
	ID          int64
	ChannelID   string
	GroupID     string
	FeedID      int64
	BotName     string
	PublishedAt time.Time
}

type SubscriptionDetail struct {
	ID       int64
	FeedName string
	FeedURL  string
}
