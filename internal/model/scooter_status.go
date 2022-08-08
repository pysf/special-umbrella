package model

import (
	"context"
	"time"
)

type ScooterStatus struct {
}

type ScooterStatusRepository interface {
	AddStatus(ctx context.Context, scooterID string, latitude float64, longitude float64, timestamp time.Time) (*string, error)
}
