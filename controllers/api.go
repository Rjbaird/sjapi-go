package controllers

import (
	"log"

	"github.com/bairrya/sjapi/config"
	"github.com/bairrya/sjapi/db"
	"github.com/gofiber/fiber/v2"
)

func ApiWelcome(c *fiber.Ctx) error {
	// Load environment variables config
	config, _ := config.ENV()
	// TODO: handle config errors
	return c.JSON(fiber.Map{
		"title":       "Welcome to the Shonen Jump API!",
		"description": "The Shonen Jump API is a JSON API for the Shonen Jump manga series.",
		"port":        config.PORT,
	})
}

func GetAllManga(c *fiber.Ctx) error {
	// TODO: add a query param for limit
	manga, err := db.FindAllManga()

	log.Print(manga)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	// TODO: render manga json
	return c.Status(200).JSON(fiber.Map{
		"data": manga,
	})
}

func GetSeries(c *fiber.Ctx) error {
	handle := c.Params("handle")
	// TODO: get handle from query param
	// TODO: get series from db
	series, err := db.FindOneManga(handle)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	// TODO: render series json
	return c.JSON(fiber.Map{
		"data": series,
	})
}
