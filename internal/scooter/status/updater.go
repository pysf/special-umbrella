package status

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pysf/special-umbrella/internal/apperror"
	"github.com/pysf/special-umbrella/internal/db"
	"github.com/pysf/special-umbrella/internal/scooter"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TripStarted    string = "trip-started"
	TripEnded      string = "trip-ended"
	TripUpdate     string = "trip-update"
	PeriodicUpdate string = "periodic-update"
)

const (
	SCOOTER_STATUS_FREE  string = "free"
	SCOOTER_STATUS_INUSE string = "occupied"
	SCOOTER_COLLECTION   string = "scooter"
)

type StatusUpdater struct {
	DB *mongo.Database
}

func NewStatusUpdater() (scooter.StatusUpdater, error) {

	client, err := db.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("NewStatusUpdater: create connection err=%w", err)
	}

	dbName, exists := os.LookupEnv("MONGODB_DATABASE")
	if !exists {
		return nil, fmt.Errorf("NewStatusUpdater: MONGODB_DATABASE is empty")
	}

	DB := client.Database(dbName)

	return &StatusUpdater{
		DB: DB,
	}, nil
}

func (sr *StatusUpdater) UpdateStatus(ctx context.Context, statusEvent scooter.ScooteStatusEvent) (*string, error) {
	coll := sr.DB.Collection(SCOOTER_COLLECTION)

	timestamp, err := time.Parse(time.RFC3339, statusEvent.Timestamp)
	if err != nil {
		return nil, apperror.NewAppError(
			apperror.WithError(fmt.Errorf("UpdateStatus: parse timestamp err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	location := scooter.Location{}
	if err := ParseLocation(statusEvent.Latitude, statusEvent.Longitude, &location); err != nil {
		return nil, apperror.NewAppError(
			apperror.WithError(fmt.Errorf("UpdateStatus: err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	status, err := decideStatus(statusEvent.EventType)
	if err != nil {
		return nil, apperror.NewAppError(
			apperror.WithError(fmt.Errorf("UpdateStatus: err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	event := bson.D{
		{Key: "_id", Value: uuid.New().String()},
		{Key: "status", Value: status},
		{Key: "type", Value: statusEvent.EventType},
		{Key: "scooterID", Value: statusEvent.ScooterID},
		{Key: "latitude", Value: location.Latitude},
		{Key: "longitude", Value: location.Longitude},
		{Key: "timestamp", Value: timestamp},
	}

	result, err := coll.InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("UpdateStatus: insert err=%w", err)
	}

	id := fmt.Sprint(result.InsertedID)
	return &id, nil

}

func ParseLocation(lat, lng string, l *scooter.Location) error {
	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return fmt.Errorf("ParseLocation: parse latitude err= %w", err)
	}

	longitude, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		return fmt.Errorf("ParseLocation: parse longitude err=%w", err)
	}

	l.Latitude = latitude
	l.Longitude = longitude

	return nil
}

func decideStatus(eventType string) (string, error) {
	switch eventType {
	case TripStarted:
		return SCOOTER_STATUS_INUSE, nil
	case TripUpdate:
		return SCOOTER_STATUS_INUSE, nil
	case TripEnded:
		return SCOOTER_STATUS_FREE, nil
	case PeriodicUpdate:
		return SCOOTER_STATUS_FREE, nil
	default:
		return "", fmt.Errorf("eventStatus: err: %v is invalid type", eventType)
	}
}
