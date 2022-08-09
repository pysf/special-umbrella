package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/pysf/special-umbrella/internal/db"
	"github.com/pysf/special-umbrella/internal/model"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

const SCOOTER_COLLECTION string = "scooter"

func NewScooterStatusRepository() (model.ScooterStatusRepository, error) {

	client, err := db.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("NewScooterStatusRepository: create connection err=%w", err)
	}

	dbName, exists := os.LookupEnv("MONGODB_DATABASE")
	if !exists {
		return nil, fmt.Errorf("NewMongoDBScooterRepository: MONGODB_DATABASE is empty")
	}

	collection := client.Database(dbName).Collection(SCOOTER_COLLECTION)

	return &ScooterStatusRepository{
		ScooterCollection: collection,
		MongoDBClient:     client,
	}, nil
}

type ScooterStatusRepository struct {
	ScooterCollection *mongo.Collection
	MongoDBClient     *mongo.Client
}

func (sr *ScooterStatusRepository) AddStatus(ctx context.Context, scooterStatusEvent model.ScooteStatusEvent) (*string, error) {
	coll := sr.ScooterCollection.Database().Collection(SCOOTER_COLLECTION)

	event := bson.D{
		{Key: "_id", Value: uuid.New().String()},
		{Key: "scooterID", Value: scooterStatusEvent.ScooterID},
		{Key: "latitude", Value: scooterStatusEvent.Location.Latitude},
		{Key: "longitude", Value: scooterStatusEvent.Location.Longitude},
		{Key: "timestamp", Value: scooterStatusEvent.Timestamp},
	}

	result, err := coll.InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("AddEvent: insert err=%w", err)
	}

	id := fmt.Sprint(result.InsertedID)
	return &id, nil

}
