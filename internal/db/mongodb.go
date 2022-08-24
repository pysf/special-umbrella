package db

import (
	"context"
	"fmt"

	"github.com/pysf/special-umbrella/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CreateConnection(ctx context.Context) (*mongo.Client, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConf().MongoURI))
	if err != nil {
		return nil, fmt.Errorf("CreateConnection: err= %w", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("CreateConnoction: ping err= %w", err)
	}

	return client, nil
}
