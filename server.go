package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bairrya/sjapi/config"
	"github.com/bairrya/sjapi/controllers"
	"github.com/bairrya/sjapi/jobs"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// load config
	config, err := config.ENV()

	if err != nil {
		log.Fatal(err)
	}

	// setup cron jobs
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Hour().Do(func() { jobs.Heartbeat() })

	// setup fiber server
	engine := html.New("./views", ".html")
	server := fiber.New(fiber.Config{Views: engine, ViewsLayout: "layouts/main"})

	// define middleware
	server.Use(logger.New())
	server.Use(helmet.New())
	server.Use(cors.New())
	server.Use(recover.New())

	// define routes
	// app routes
	// GET /rss (rss feed)
	// GET lite (lite version of all manga)
	server.Get("/", controllers.MainPage)
	// server.Get("/manga", controllers.AllMangaPage)
	// server.Get("/manga/:handle", controllers.SeriesPage)
	server.Get("/lite", controllers.LitePage)
	// server.Get("/rss", controllers.RssFeed)

	// api routes
	api := server.Group("/api")
	api.Use(requestid.New())
	// TODO: add limiter
	api.Get("/", controllers.ApiWelcome)
	api.Get("manga", controllers.GetAllManga)
	api.Get("manga/:handle", controllers.GetSeries)

	// Start jobs and server
	s.StartAsync()
	log.Fatal(server.Listen(fmt.Sprintf(":%s", config.PORT)))

	// TODO: add xml engine https://github.com/gofiber/recipes/blob/master/rss-feed/main.go
	// TODO: add swaggerdocs https://github.com/gofiber/swagger
	// TODO: add api tests https://www.youtube.com/watch?v=XQzTUa9LPU8, https://www.youtube.com/watch?v=Ztk9d78HgC0
}
