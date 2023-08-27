package main

import (
	"context"
	"fmt"

	"github.com/bairrya/sjapi/db"
	"github.com/bairrya/sjapi/web"
)

func main() {

	// Connect to MongoDB
	client, err := db.Init()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("shonen-jump")
	mangaCollection := db.Collection("manga")

	results, err := web.ScrapeRecentChapters()
	if err != nil {
		panic(err)
	}

	mangaDocs := []interface{}{}
	for _, manga := range *results {
		mangaDocs = append(mangaDocs, manga)
	}

	insert, err := mangaCollection.InsertMany(context.TODO(), mangaDocs)
	if err != nil {
		panic(err)
	}

	fmt.Println(insert)
}
