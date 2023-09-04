package controllers

import (
	"github.com/bairrya/sjapi/db"
	"github.com/gofiber/fiber/v2"
)

func ApiWelcome(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"title":        "Welcome to the Shonen Jump API!",
		"description":  "The Shonen Jump API is a JSON API for the Shonen Jump manga series.",
		"last_updated": "2023-09-01",
	})
}

func GetAllManga(c *fiber.Ctx) error {
	// TODO: add a query param for limit
	manga, err := db.FindAllManga()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": manga,
	})
}

func GetSeries(c *fiber.Ctx) error {
	handle := c.Params("handle")
	series, err := db.FindOneManga(handle)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	return c.JSON(fiber.Map{
		"data": series,
	})
}

func LimitReachedJSON(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"data": "Rate limit reached. Please try again later.",
	})
}
