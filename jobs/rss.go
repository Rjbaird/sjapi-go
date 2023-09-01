package jobs

import (
	"time"

	"github.com/bairrya/sjapi/db"
	"github.com/gorilla/feeds"
)

// Description: This is a job that will generate a new RSS feed for the latest chapters

func GenerateRssFeed() (interface{}, error) {
	results, err := db.FindRecentManga()
	if err != nil {
		return nil, err
	}

	feed := &feeds.Feed{
		Title:       "Title of Your Feeds",
		Link:        &feeds.Link{Href: "your feed link"},
		Description: "Description of your feeds",
		Author:      &feeds.Author{Name: "author name"},
		Created:     time.Now(),
	}
	var feedItems []*feeds.Item
	for _, manga := range *results {
		feedItems = append(feedItems,
			&feeds.Item{
				Id:          manga.Handle,
				Title:       manga.Title,
				Link:        &feeds.Link{Href: "feeds item link"},
				Description: manga.Description,
				Created:     time.Now(),
			})
	}
	feed.Items = feedItems
	rssFeed := (&feeds.Rss{Feed: feed}).RssFeed()
	xmlRssFeeds := rssFeed.FeedXml()
	return xmlRssFeeds, nil
}
