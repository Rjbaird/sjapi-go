package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecentManga struct {
	// ID            primitive.ObjectID `json:"id" bson:"_id,omitempts"`
	Title         string `json:"title" bson:"title,omitempts"`
	Link          string `json:"link" bson:"link,omitempts"`
	RecentChapter string `json:"recent_chapters" bson:"recent_chapter,omitempts"`
	RecentLink    string `json:"recent_link" bson:"recent_link,omitempts"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	db := client.Database("shonen-jump")
	mangaCollection := db.Collection("manga")
	defer client.Disconnect(context.Background())
	const recentURL = "https://www.viz.com/read/shonenjump/section/free-chapters"

	results := []RecentManga{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set()
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML(".o_sortable", func(e *colly.HTMLElement) {
		recent := strings.Split(e.ChildText("span"), "  ")[0]
		manga := RecentManga{
			Title:         e.ChildText("div.type-center"),
			Link:          e.ChildAttr("a", "href"),
			RecentChapter: strings.Replace(recent, "Latest: ", "", -1),
			RecentLink:    e.ChildAttr("a.o_inner-link", "href"),
		}

		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", manga)
		results = append(results, manga)
	})
	c.Visit(recentURL)

	mangaDocs := []interface{}{}
	for _, manga := range results {
		mangaDocs = append(mangaDocs, manga)
	}

	// fmt.Printf("%+v\n", mangaDocs...)
	result, err := mangaCollection.InsertMany(context.TODO(), mangaDocs)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
