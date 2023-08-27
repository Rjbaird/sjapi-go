package main

import (
	"fmt"

	"github.com/bairrya/sjapi/web"
)

func main() {

	results, err := web.ScrapeSeries("naruto")
	if err != nil {
		panic(err)
	}

	fmt.Println(results)

	// Connect to MongoDB
	// client, err := db.Init()
	// if err != nil {
	// 	panic(err)
	// }
	// defer client.Disconnect(context.Background())

	// db := client.Database("shonen-jump")
	// mangaCollection := db.Collection("manga")

	// mangaDocs := []interface{}{}
	// for _, manga := range *results {
	// 	mangaDocs = append(mangaDocs, manga)
	// }

	// insert, err := mangaCollection.InsertMany(context.TODO(), mangaDocs)
	// if err != nil {
	// 	panic(err)
	// }

	// if insert != nil {
	// 	fmt.Println("Inserted multiple documents")
	// }
}
