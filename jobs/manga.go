package jobs

import (
	"log"

	"github.com/bairrya/sjapi/db"
	"github.com/bairrya/sjapi/web"
)

func RefreshAllManga() error {
	// Scrape manga from viz
	manga, err := web.ScrapeAllSeries()
	if err != nil {
		log.Fatal(err)
		return err
	}

	// TODO: update instead of drop and replace

	// Drop all manga from db
	dropErr := db.DropMangaCollection()
	if dropErr != nil {
		log.Fatal(dropErr)
		return dropErr
	}


	// Insert manga into db
	err = db.InsertManyManga(manga)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func RefreshRecentManga() error {
	return nil
}
