package jobs

import (
	"log"

	"github.com/bairrya/sjapi/db"
	"github.com/bairrya/sjapi/web"
)

func RefreshAllManga() error {
	// Drop all manga from db
	err := db.DropMangaCollection()
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Scrape manga from viz
	manga, err := web.ScrapeAllSeries()
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Insert manga into db
	err = db.InsertManyManga(manga)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
