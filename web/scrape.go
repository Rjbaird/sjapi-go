package web

import (
	"fmt"
	"log"
	"sort"
	"strconv"
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
	// TODO: download images to exertnal storage and save link to db
	Title         string             `json:"title" bson:"title,omitempts"`
	Handle        string             `json:"handle" bson:"handle,omitempts"`
	Description   string             `json:"description" bson:"description,omitempts"`
	Author        string             `json:"author" bson:"author,omitempts"`
	HeroImage     string             `json:"hero_image" bson:"image,omitempts"`
	LatestChapter float64            `json:"latest_chapter" bson:"latest_chapter,omitempts"`
	NextRelease   string             `json:"next_release" bson:"next_release,omitempts"`
	Mature        bool               `json:"mature" bson:"mature,omitempts"`
	Recommended   []RecommendedManga `json:"recommended" bson:"recommended,omitempts"`
	Chapters      []Chapter          `json:"chapters" bson:"chapters,omitempts"`
	Volumes       []Volumn           `json:"volumes" bson:"volumes,omitempts"`
}

type Chapter struct {
	ID     string  `json:"id" bson:"id,omitempts"`
	Number float64 `json:"number" bson:"number,omitempts"`
	Date   string  `json:"date" bson:"date,omitempts"`
	Link   string  `json:"link" bson:"link,omitempts"`
	Free   bool    `json:"free" bson:"free,omitempts"`
	Mature bool    `json:"mature" bson:"mature,omitempts"`
}

type Volumn struct {
	Number   string    `json:"number" bson:"number,omitempts"`
	Link     string    `json:"link" bson:"link,omitempts"`
	Chapters []Chapter `json:"chapters" bson:"chapters,omitempts"`
	// TODO: add extra info from volumn link to volumn struct
	// Image
	// Description
	// Release date
	// ISBN
}

type RecommendedManga struct {
	Title  string `json:"title" bson:"title,omitempts"`
	Handle string `json:"handle" bson:"handle,omitempts"`
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
		recommended := []RecommendedManga{}
		chapters := []Chapter{}
		volumns := []Volumn{}

		// basic info
		title := e.ChildText("h2.type-lg")
		description := e.ChildText("div.line-solid.type-md")
		author := e.ChildText("span.disp-bl--bm")
		image := e.ChildAttr("img.o_hero-media", "src")
		next_release := e.ChildText("div.section_future_chapter")

		// chapters
		e.ForEach("div.o_sortable", func(_ int, el *colly.HTMLElement) {
			free := true
			mature := false
			chapter_link_class := el.ChildAttr("a.o_chapter-container", "class")

			if strings.Contains(chapter_link_class, "o_m-rated") {
				mature = true
			}

			if strings.Contains(chapter_link_class, "o_chapter-archive") {
				free = false
			}

			chapter_id := el.ChildAttr("a.o_chapter-container", "id")
			chapter_string := el.ChildAttr("a.o_chapter-container", "name")
			chapter_date := el.ChildText("td.pad-r-0")

			chapter_number := 0.0

			chapter_link := el.ChildAttr("a.o_chapter-container", "href")
			chapter_float, err := strconv.ParseFloat(chapter_string, 64)

			if err == nil {
				chapter_number = chapter_float
			}

			chapter := Chapter{
				ID:     chapter_id,
				Number: chapter_number,
				Date:   chapter_date,
				Link:   chapter_link,
				Free:   free,
				Mature: mature,
			}
			chapters = append(chapters, chapter)
		})

		// volumes
		e.ForEach("div.o_chapter-vol-container", func(_ int, el *colly.HTMLElement) {
			volumn_chapters := []Chapter{}
			volumn_number := strings.Split(el.ChildAttr("a.o_manga-buy-now", "aria-label"), "Vol. ")[0]
			volumn_link := el.ChildAttr("a.o_manga-buy-now", "href")
			el.ForEach("tr.o_chapter", func(_ int, el *colly.HTMLElement) {
				free := true
				mature := false
				chapter_link_class := el.ChildAttr("a.o_chapter-container", "class")

				if strings.Contains(chapter_link_class, "o_m-rated") {
					mature = true
				}

				if strings.Contains(chapter_link_class, "o_chapter-archive") {
					free = false
				}

				chapter_id := el.ChildAttr("a.o_chapter-container", "id")
				chapter_number := el.ChildAttr("a.o_chapter-container", "name")
				chapter_date := el.ChildText("td.pad-r-0")
				chapter_link := el.ChildAttr("a.o_inner-link", "href")
				chapter_float, _ := strconv.ParseFloat(chapter_number, 64)

				chapter := Chapter{
					ID:     chapter_id,
					Number: chapter_float,
					Date:   chapter_date,
					Link:   chapter_link,
					Free:   free,
					Mature: mature,
				}
				volumn_chapters = append(volumn_chapters, chapter)
			})

			volumn := Volumn{
				Number:   volumn_number,
				Link:     "https://www.viz.com" + volumn_link,
				Chapters: volumn_chapters,
			}
			volumns = append(volumns, volumn)
		})

		// recommended manga
		e.ForEach("a.o_property-link", func(_ int, el *colly.HTMLElement) {
			// TODO: add recommended manga titles to struct
			title := el.Attr("rel")
			handle := strings.Replace(el.Attr("href"), "/shonenjump/chapters/", "", -1)

			recommended = append(recommended, RecommendedManga{
				Title:  title,
				Handle: handle,
			})
		})

		uniqueChapters := removeDuplicates(chapters)
		uniqueChapters = removeEmptyChapters(uniqueChapters)

		sort.Slice(uniqueChapters, func(i, j int) bool {
			return uniqueChapters[i].Number > uniqueChapters[j].Number
		})

		manga := Manga{
			Title:         title,
			Handle:        handle,
			Description:   description,
			Author:        author,
			HeroImage:     image,
			LatestChapter: uniqueChapters[0].Number,
			NextRelease:   next_release,
			Recommended:   recommended,
			Chapters:      uniqueChapters,
			Volumes:       volumns,
			Mature:        containsMature(chapters),
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

func containsMature(chapters []Chapter) bool {
	for _, chapter := range chapters {
		if chapter.Mature {
			return true
		}
	}
	return false
}

func removeDuplicates(chapter []Chapter) []Chapter {
	var unique []Chapter
sampleLoop:
	for _, v := range chapter {
		for i, u := range unique {
			if v.Number == u.Number {
				unique[i] = v
				continue sampleLoop
			}
		}
		unique = append(unique, v)
	}
	return unique
}

func removeEmptyChapters(chapters []Chapter) []Chapter {
	var nonEmpty []Chapter
	for _, v := range chapters {
		if v.Number != 0 {
			nonEmpty = append(nonEmpty, v)
		}
	}
	return nonEmpty
}
