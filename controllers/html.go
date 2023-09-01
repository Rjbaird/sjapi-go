package controllers

import (
	"fmt"
	"time"

	"github.com/bairrya/sjapi/db"
	"github.com/bairrya/sjapi/jobs"
	"github.com/gofiber/fiber/v2"
)

func MainPage(c *fiber.Ctx) error {
	now := time.Now()
	today := fmt.Sprintf("%s, %s %v, %v", now.Weekday(), now.Month(), now.Day(), now.Year())
	manga, err := db.FindRecentManga()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	d := map[string]interface{}{}

	d["Manga"] = manga
	d["Meta"] = today

	return c.Render("manga", fiber.Map{
		"Data": d,
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
	return c.Render("series", fiber.Map{
		"Series": series,
	})
}

func ContactPage(c *fiber.Ctx) error {
	// Render contact template
	return c.Render("contact", fiber.Map{})
}

func RssFeed(c *fiber.Ctx) error {
	feed, err := jobs.GenerateRssFeed()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	c.Type("xml")
	return c.XML(feed)
}
