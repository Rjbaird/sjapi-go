package db

import (
	"context"

	"github.com/bairrya/sjapi/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func Init() (*mongo.Client, error) {
	config, err := config.ENV()

	if err != nil {
		return nil, err
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MONGODB_URI))
	if err != nil {
		return nil, err
	}
	return client, nil
}
