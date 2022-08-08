package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pysf/special-umbrella/internal/model"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

const SCOOTER_COLLECTION string = "scooter"

func NewMongoDBScooterRepository(client *mongo.Client) (model.ScooterStatusRepository, error) {
	dbName, exists := os.LookupEnv("MONGODB_DATABASE")
	if !exists {
		return nil, fmt.Errorf("NewMongoDBScooterRepository: MONGODB_DATBASE is empty")
	}

	collection := client.Database(dbName).Collection(SCOOTER_COLLECTION)

	return &MongoDBScooterRepository{
		ScooterCollection: collection,
		MongoDBClient:     client,
	}, nil
}

type MongoDBScooterRepository struct {
	ScooterCollection *mongo.Collection
	MongoDBClient     *mongo.Client
}

func (sr *MongoDBScooterRepository) AddStatus(ctx context.Context, scooterID string, latitude float64, longitude float64, timestamp time.Time) (*string, error) {
	coll := sr.ScooterCollection.Database().Collection(SCOOTER_COLLECTION)

	id := uuid.New().String()
	event := bson.D{
		{Key: "_id", Value: id},
		{Key: "scooterID", Value: scooterID},
		{Key: "latitude", Value: latitude},
		{Key: "longitude", Value: longitude},
		{Key: "timestamp", Value: timestamp},
	}

	result, err := coll.InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("AddEvent: insert err=%w", err)
	}

	ss := fmt.Sprint(result.InsertedID)
	return &ss, nil

}
