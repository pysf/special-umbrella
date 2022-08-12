package db

import (
	"context"
	"fmt"

	"github.com/pysf/special-umbrella/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CreateConnection() (*mongo.Client, error) {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.GetConfig("MONGODB_URI")))
	if err != nil {
		return nil, fmt.Errorf("CreateConnection: err= %w", err)
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, fmt.Errorf("CreateConnoction: ping err= %w", err)
	}

	return client, nil
}
