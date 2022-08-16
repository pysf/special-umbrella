package scooteriface

import (
	"context"

	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
)

type ScooterStatusUpdater interface {
	UpdateStatus(context.Context, struct {
		ScooterID string
		Timestamp string
		Latitude  string
		Longitude string
		EventType string
	}) (*string, error)
}

type ScooterFinder interface {
	RectangularQuery(context.Context, struct {
		Status     string
		BottomLeft scootertype.Location
		TopRight   scootertype.Location
	}) (*scootertype.RectangularQueryResult, error)
}

type ScooterReserver interface {
	ReserveScooter(ctx context.Context, scooterID string) (bool, error)
	ReleaseScooter(ctx context.Context, scooterID string) error
}

type ScooterCreator interface {
	Create(ctx context.Context, scooter scootertype.Scooter) error
}
