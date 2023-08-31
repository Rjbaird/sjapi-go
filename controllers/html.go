package controllers

import (
	"log"

	"github.com/bairrya/sjapi/db"
	"github.com/gofiber/fiber/v2"
)

func MainPage(c *fiber.Ctx) error {
	manga, err := db.FindRecentManga()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	return c.Render("manga", fiber.Map{
		"Manga": manga,
	})
}

func AllMangaPage(c *fiber.Ctx) error {
	// Render manga template
	manga, err := db.FindAllManga()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	return c.Render("manga", fiber.Map{
		"Manga": manga,
	})
}

func SeriesPage(c *fiber.Ctx) error {
	// Render series template
	handle := c.Params("handle")
	series, err := db.FindOneManga(handle)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			// TODO: update w/ fiber error for all 500s
			"error": "Internal Server Error",
		})
	}
	log.Println(series)
	return c.Render("series", fiber.Map{
		"Series": series,
	})
}

func RssFeed(c *fiber.Ctx) error {
	// Render rss template
	return c.Render("rss", fiber.Map{
		"Title": "Shonen Jump | Read Free Manga Online!",
	})
}
