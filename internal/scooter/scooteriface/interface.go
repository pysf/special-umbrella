package scooteriface

import (
	"context"

	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
)

type ScooterStatusUpdater interface {
	UpdateStatus(ctx context.Context, scooterStatusUpdaterInput scootertype.ScooterStatusUpdaterInput) (*string, error)
}

type ScooterFinder interface {
	RectangularQuery(ctx context.Context, bottomLeft scootertype.Location, topRigth scootertype.Location) (*scootertype.RectangularQueryResult, error)
}
