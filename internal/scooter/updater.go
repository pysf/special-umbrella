package scooter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pysf/special-umbrella/internal/apperror"
	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/db"
	"github.com/pysf/special-umbrella/internal/scooter/scooteriface"
	"github.com/pysf/special-umbrella/internal/scooter/scootertype"

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

func NewStatusUpdater() (scooteriface.ScooterStatusUpdater, error) {

	client, err := db.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("NewStatusUpdater: create connection err=%w", err)
	}

	DB := client.Database(config.GetConfig("MONGODB_DATABASE"))

	return &StatusUpdater{
		DB: DB,
	}, nil
}

func (s *StatusUpdater) UpdateStatus(ctx context.Context, updateStatusInput scootertype.ScooterStatusUpdaterInput) (*string, error) {

	timestamp, err := time.Parse(time.RFC3339, updateStatusInput.Timestamp)
	if err != nil {
		return nil, apperror.NewAppError(
			apperror.WithError(fmt.Errorf("UpdateStatus: parse timestamp err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	location := scootertype.GeoJSON{}
	if err := parseLocation(updateStatusInput.Latitude, updateStatusInput.Longitude, &location); err != nil {
		return nil, apperror.NewAppError(
			apperror.WithError(fmt.Errorf("UpdateStatus: err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	status, err := decideStatus(updateStatusInput.EventType)
	if err != nil {
		return nil, apperror.NewAppError(
			apperror.WithError(fmt.Errorf("UpdateStatus: err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	event := scootertype.ScooterStatus{
		ID:        uuid.New().String(),
		Status:    status,
		EventType: updateStatusInput.EventType,
		ScooterID: updateStatusInput.ScooterID,
		Timestamp: timestamp,
		Location:  location,
	}

	result, err := s.DB.Collection(SCOOTER_COLLECTION).InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("UpdateStatus: insert err=%w", err)
	}

	id := fmt.Sprint(result.InsertedID)
	return &id, nil

}

func parseLocation(lat, lng string, l *scootertype.GeoJSON) error {
	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return fmt.Errorf("ParseLocation: parse latitude err= %w", err)
	}

	longitude, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		return fmt.Errorf("ParseLocation: parse longitude err=%w", err)
	}

	l.Type = "Point"
	l.Coordinates = []float64{latitude, longitude}

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