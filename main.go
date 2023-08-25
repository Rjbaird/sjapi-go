package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type RecentManga struct {
	Title         string
	Link          string
	RecentChapter string
	RecentLink    string
}

const recentURL = "https://www.viz.com/read/shonenjump/section/free-chapters"

func main() {
	results := []RecentManga{}
	c := colly.NewCollector()
    c.OnRequest(func(r *colly.Request) {
        // r.Headers.Set()
        fmt.Println("Visiting", r.URL)
    })

    c.OnError(func(_ *colly.Response, err error) {
        fmt.Println("Something went wrong:", err)
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
	data, err := json.MarshalIndent(results, " ", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}
