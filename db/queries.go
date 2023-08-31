package db

import (
	"context"
	"fmt"

	"github.com/bairrya/sjapi/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllManga() (*[]web.Manga, error) {

	client, err := Init()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())
	coll := client.Database("shonen-jump").Collection("manga")

	manga := []web.Manga{}
	// TODO: add a limit to the number of results
	filter := bson.D{{}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	cursor.All(context.Background(), &manga)
	return &manga, nil
}

func FindRecentManga() (*[]web.Manga, error) {
	client, err := Init()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())
	coll := client.Database("shonen-jump").Collection("manga")

	manga := []web.Manga{}
	filter := bson.D{{Key: "latest_release", Value: bson.D{{Key: "$ne", Value: ""}}}}
	opts := options.Find().SetLimit(30)
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	cursor.All(context.Background(), &manga)
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
		fmt.Printf("Inserted %v documents\n", len(insert.InsertedIDs))
	}
	return nil
}
