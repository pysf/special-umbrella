package model

import (
	"context"
	"time"
)

type ScooterStatus struct {
}

type ScooterStatusRepository interface {
	AddStatus(ctx context.Context, scooterStatusEvent ScooteStatusEvent) (*string, error)
}

type Location struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
}

type ScooteStatusEvent struct {
	ScooterID string    `json:"scooterId"`
	Timestamp time.Time `json:"timestamp"`
	Location  Location  `json:"location"`
}
