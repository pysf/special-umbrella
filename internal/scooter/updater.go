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
	EVENT_TRIP_STARTED    string = "trip-started"
	EVENT_TRIP_ENDED      string = "trip-ended"
	EVENT_TRIP_UPDATED    string = "trip-update"
	EVENT_PERIODIC_UPDATE string = "periodic-update"
)

const (
	SCOOTER_STATUS_AVAILABLE  string = "available"
	SCOOTER_STATUS_INUSE      string = "inuse"
	SCOOTER_STATUS_COLLECTION string = "scooter-status"
)

type StatusUpdater struct {
	DB              *mongo.Database
	ScooterReserver scooteriface.ScooterReserver
}

func NewStatusUpdater(ScooterReserver scooteriface.ScooterReserver) (scooteriface.ScooterStatusUpdater, error) {

	client, err := db.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("NewStatusUpdater: create connection err=%w", err)
	}

	DB := client.Database(config.GetConfig("MONGODB_DATABASE"))

	return &StatusUpdater{
		DB:              DB,
		ScooterReserver: ScooterReserver,
	}, nil
}

func (s *StatusUpdater) UpdateStatus(ctx context.Context, updateStatusInput struct {
	ScooterID string
	Timestamp string
	Latitude  string
	Longitude string
	EventType string
}) (*string, error) {

	timestamp, err := time.Parse(time.RFC3339, updateStatusInput.Timestamp)
	if err != nil {
		return nil, apperror.NewAppError(
			apperror.WithError(fmt.Errorf("UpdateStatus: parse timestamp err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	var location scootertype.GeoJSON
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

	result, err := s.DB.Collection(SCOOTER_STATUS_COLLECTION).InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("UpdateStatus: insert err=%w", err)
	}

	id := fmt.Sprint(result.InsertedID)
	// we need to release the scooter after the end-event received
	if updateStatusInput.EventType == EVENT_TRIP_ENDED {
		if err = s.ScooterReserver.ReleaseScooter(ctx, updateStatusInput.ScooterID); err != nil {
			return nil, fmt.Errorf("UpdateStatus: release scooter err= %w", err)
		}
	}

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
	case EVENT_TRIP_STARTED:
		return SCOOTER_STATUS_INUSE, nil
	case EVENT_TRIP_UPDATED:
		return SCOOTER_STATUS_INUSE, nil
	case EVENT_TRIP_ENDED:
		return SCOOTER_STATUS_AVAILABLE, nil
	case EVENT_PERIODIC_UPDATE:
		return SCOOTER_STATUS_AVAILABLE, nil
	default:
		return "", fmt.Errorf("eventStatus: err: %v is invalid type", eventType)
	}
}
