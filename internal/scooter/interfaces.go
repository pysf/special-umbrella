package scooter

import (
	"context"
)

type StatusUpdater interface {
	UpdateStatus(ctx context.Context, scooterStatusEvent ScooteStatusEvent) (*string, error)
}
