package scooter

import (
	"context"
	"fmt"

	"github.com/pysf/special-umbrella/internal/scooter/scooteriface"
	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	SCOOTER_COLLECTION = "scooter"
)

type ScooterCreator struct {
	DB *mongo.Database
}

func NewScooterCreator(DB *mongo.Database) (scooteriface.ScooterCreator, error) {

	return &ScooterCreator{
		DB: DB,
	}, nil
}

func (c *ScooterCreator) Create(ctx context.Context, scooter scootertype.Scooter) error {
	coll := c.DB.Collection(SCOOTER_COLLECTION)
	if _, err := coll.InsertOne(ctx, scootertype.Scooter{
		ID:    scooter.ID,
		InUse: scooter.InUse,
	}); err != nil {
		return fmt.Errorf("Create: failed to create scooter err=%w", err)
	}
	return nil
}
