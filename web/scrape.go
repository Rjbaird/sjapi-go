package web

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type RecentManga struct {
	Title         string `json:"title" bson:"title,omitempts"`
	Link          string `json:"link" bson:"link,omitempts"`
	RecentChapter string `json:"recent_chapters" bson:"recent_chapter,omitempts"`
	RecentLink    string `json:"recent_link" bson:"recent_link,omitempts"`
}

func ScrapeRecentChapters() (*[]RecentManga, error) {
	const recentURL = "https://www.viz.com/read/shonenjump/section/free-chapters"

	results := []RecentManga{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Panic("Something went wrong getting recent chapters:", err)
	})

	c.OnHTML(".o_sortable", func(e *colly.HTMLElement) {
		recent := strings.Split(e.ChildText("span"), "  ")[0]
		manga := RecentManga{
			Title:         e.ChildText("div.type-center"),
			Link:          e.ChildAttr("a", "href"),
			RecentChapter: strings.Replace(recent, "Latest: ", "", -1),
			RecentLink:    e.ChildAttr("a.o_inner-link", "href"),
		}
		results = append(results, manga)
	})
	c.Visit(recentURL)
	return &results, nil
}
