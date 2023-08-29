package web

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type RecentManga struct {
	Title         string `json:"title" bson:"title,omitempts"`
	Handle        string `json:"handle" bson:"handle,omitempts"`
	RecentChapter string `json:"recent_chapters" bson:"recent_chapter,omitempts"`
	RecentLink    string `json:"recent_link" bson:"recent_link,omitempts"`
}

type Manga struct {
	Title         string   `json:"title" bson:"title,omitempts"`
	Handle        string   `json:"handle" bson:"handle,omitempts"`
	Description   string   `json:"description" bson:"description,omitempts"`
	Author        string   `json:"author" bson:"author,omitempts"`
	HeroImage     string   `json:"hero_image" bson:"image,omitempts"`
	LatestRelease string   `json:"latest_release" bson:"latest_release,omitempts"`
	Recommended   []string `json:"recommended" bson:"recommended,omitempts"`
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
		recent := strings.Split(e.ChildText("span"), "\n")[0]
		manga := RecentManga{
			Title:         e.ChildText("div.type-center"),
			Handle:        strings.Replace(e.ChildAttr("a", "href"), "/shonenjump/chapters/", "", -1),
			RecentChapter: strings.Replace(recent, "Latest: ", "", -1),
			RecentLink:    e.ChildAttr("a.o_inner-link", "href"),
		}
		results = append(results, manga)
	})
	c.Visit(recentURL)
	return &results, nil
}

func ScrapeOneSeries(handle string) (*Manga, error) {
	seriesURL := "https://www.viz.com/shonenjump/chapters/" + handle
	results := Manga{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		s := fmt.Sprintf("Something went wrong getting data on %s: ", handle)
		log.Panic(s, err)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		recommended := []string{}
		title := e.ChildText("h2.type-lg")
		description := e.ChildText("div.line-solid.type-md")
		author := e.ChildText("span.disp-bl--bm")
		image := e.ChildAttr("img.o_hero-media", "src")
		latest_release := strings.Split(e.ChildText("div.flex-width-2.type-bs.type-sm--sm.style-italic > table > tbody > tr > td"), ", ")[0]
		e.ForEach("a.o_property-link", func(_ int, el *colly.HTMLElement) {
			handle := strings.Replace(el.Attr("href"), "/shonenjump/chapters/", "", -1)
			recommended = append(recommended, handle)
		})
		manga := Manga{
			Title:         title,
			Handle:        handle,
			Description:   description,
			Author:        author,
			HeroImage:     image,
			LatestRelease: latest_release,
			Recommended:   recommended,
		}
		results = manga
	})
	c.Visit(seriesURL)
	return &results, nil
}

func ScrapeAllSeries() (*[]Manga, error) {
	allSeries, err := ScrapeRecentChapters()
	if err != nil {
		return nil, err
	}
	estimateTime := secondsToMinutes(len(*allSeries) * 3)
	fmt.Printf("Estimated time to scrape all series: %v\n", estimateTime)
	mangaSlice := []Manga{}
	fmt.Println("Scraping all series...")
	for _, series := range *allSeries {
		manga, err := ScrapeOneSeries(series.Handle)
		if err != nil {
			// NOTE: can I catch this error and continue by appending the series to the slice? Or does that cause an infinite loop?
			return nil, err
		}
		fmt.Printf("Scraped %s\n", manga.Title)
		mangaSlice = append(mangaSlice, *manga)
		time.Sleep(3 * time.Second)
	}
	return &mangaSlice, nil
}

func secondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	str := fmt.Sprintf("%v:%v", minutes, seconds)
	return str
}
