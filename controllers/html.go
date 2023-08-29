package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func MainPage(c *fiber.Ctx) error {
	// Render index template
	return c.Render("index", fiber.Map{
		"Title": "Shonen Jump | Read Free Manga Online!",
	})
}

func LitePage(c *fiber.Ctx) error {
	// Render lite template
	return c.Render("lite", fiber.Map{
		"Title": "Shonen Jump Lite | Read Free Manga Online!",
	})
}

func AllMangaPage(c *fiber.Ctx) error {
	// Render manga template
	return c.Render("manga", fiber.Map{
		"Title": "Shonen Jump | Read Free Manga Online!",
	})
}

func SeriesPage(c *fiber.Ctx) error {
	// Render series template
	// TODO: get handle from query param
	// TODO: get series from db
	// TODO: render series template
	return c.Render("series", fiber.Map{
		"Title": "Shonen Jump | Read Free Manga Online!",
	})
}

func RssFeed(c *fiber.Ctx) error {
	// Render rss template
	return c.Render("rss", fiber.Map{
		"Title": "Shonen Jump | Read Free Manga Online!",
	})
}
