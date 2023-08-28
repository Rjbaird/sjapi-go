package db

import (
	"context"
	"fmt"

	"github.com/bairrya/sjapi/web"
	"go.mongodb.org/mongo-driver/bson"
)

func FindAllManga() (*[]web.Manga, error) {
	client, err := Init()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	manga := make([]web.Manga, 0)
	cursor, err := client.Database("shonen-jump").Collection("manga").Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var series web.Manga
		if err = cursor.Decode(&series); err != nil {
			return nil, err
		}
		manga = append(manga, series)
	}

	return &manga, nil
}

func FindOneManga(handle string) (*web.Manga, error) {
	client, err := Init()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	coll := client.Database("shonen-jump").Collection("manga")
	filter := bson.D{{Key: "handle", Value: handle}}

	var result web.Manga
	err = coll.FindOne(context.TODO(), filter).Decode(&result)

	return &result, err
}

func UpdateOneManga(handle string) error {
	return nil
}

func InsertManyManga(results *[]web.Manga) error {
	// Connect to MongoDB
	client, err := Init()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	db := client.Database("shonen-jump")
	mangaCollection := db.Collection("manga")

	mangaDocs := []interface{}{}
	for _, manga := range *results {
		mangaDocs = append(mangaDocs, manga)
	}

	insert, err := mangaCollection.InsertMany(context.TODO(), mangaDocs)
	if err != nil {
		return err
	}

	if insert != nil {
		fmt.Printf("Inserted %v documents", len(insert.InsertedIDs))
	}
	return nil
}
