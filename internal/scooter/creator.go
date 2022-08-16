package scooter

import (
	"context"
	"fmt"

	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/db"
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

func NewScooterCreator() (scooteriface.ScooterCreator, error) {

	client, err := db.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("NewScooterCreator: create connection err=%w", err)
	}

	DB := client.Database(config.GetConfig("MONGODB_DATABASE"))

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
