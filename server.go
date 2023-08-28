package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bairrya/sjapi/config"
	"github.com/bairrya/sjapi/jobs"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// load config
	config, err := config.ENV()

	if err != nil {
		log.Fatal(err)
	}

	// setu cron jobs
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Hour().Do(func() { jobs.Heartbeat() })

	// setu fiber server
	app := fiber.New()

	// define routes
	// app routes
	// GET / (home)
	// GET /manga/:handle
	// GET /rss
	// GET lite (lite version of site)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"title": "Welcome to Shonen Jump API!",
			"port":  config.PORT,
		})
	})

	// api routes
	// api := app.Group("/api")
	// GET /api/manga
	// GET /api/manga/:handle

	// Start jobs and server
	s.StartAsync()
	log.Fatal(app.Listen(fmt.Sprintf(":%s", config.PORT)))
	// DONE: add cron jobs
	// DONE: add fiber server
	// TODO: add routes
	// TODO: add middleware
	// TODO: add controllers
	// TODO: add views
	// TODO: add xml engine https://github.com/gofiber/recipes/blob/master/rss-feed/main.go
	// TODO: add unocss
	// TODO: finish db queries https://www.youtube.com/watch?v=E6NmSKSUj9g, https://www.mongodb.com/docs/drivers/go/current/usage-examples/findOne/
	// TODO: add swaggerdocs https://github.com/gofiber/swagger
	// TODO: add api tests https://www.youtube.com/watch?v=XQzTUa9LPU8, https://www.youtube.com/watch?v=Ztk9d78HgC0
}
