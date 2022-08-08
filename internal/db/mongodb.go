package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CreateConnection() (*mongo.Client, error) {

	mongodbURI, exist := os.LookupEnv("MONGODB_URI")
	if !exist {
		return nil, fmt.Errorf("CreateConnection: MONGODB_URI is not set")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodbURI))
	if err != nil {
		return nil, fmt.Errorf("CreateConnection: err= %w", err)
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, fmt.Errorf("CreateConnoction: ping err= %w", err)
	}

	return client, nil
}
