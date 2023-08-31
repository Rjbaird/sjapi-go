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
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
	server.Get("/", controllers.MainPage)
	server.Get("/manga", controllers.AllMangaPage)
	server.Get("/manga/:handle", controllers.SeriesPage)
	server.Get("/rss", controllers.RssFeed)

	// api routes
	api := server.Group("/api")
	api.Use(requestid.New())
	api.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:          20,
		Expiration:   30 * time.Second,
		LimitReached: controllers.LimitReachedJSON,
	}))

	api.Get("/", controllers.ApiWelcome)
	api.Get("manga", controllers.GetAllManga)
	api.Get("manga/:handle", controllers.GetSeries)

	// Start jobs and server
	s.StartAsync()
	log.Fatal(server.Listen(fmt.Sprintf(":%s", config.PORT)))

	// NOTE: TODOs
	// TODO: add metadata to /api {version, uptime, docs, last update, etc}
	// TODO: add api tests https://www.youtube.com/watch?v=XQzTUa9LPU8, https://www.youtube.com/watch?v=Ztk9d78HgC0
	// TODO: add daily manga update cron job

	// TODO: add limiter
	// TODO: add xml engine https://github.com/gofiber/recipes/blob/master/rss-feed/main.go
	// TODO: add swaggerdocs https://github.com/gofiber/swagger
	// TODO: create routes package
	// TODO: add /images route for promotional hero images
	// TODO: add jwt auth to api routes
}
