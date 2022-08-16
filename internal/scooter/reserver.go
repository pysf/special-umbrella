package scooter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pysf/special-umbrella/internal/apperror"
	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScooterReserver struct {
	DB *mongo.Database
}

func NewScooterReserver() (*ScooterReserver, error) {

	client, err := db.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("NewScooterReserver: err=%w", err)
	}

	DB := client.Database(config.GetConfig("MONGODB_DATABASE"))

	return &ScooterReserver{
		DB: DB,
	}, nil
}

func (s *ScooterReserver) ReserveScooter(ctx context.Context, scooterID string) (bool, error) {

	query := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{
					{Key: "_id", Value: scooterID},
				},
				bson.D{
					{Key: "inuse", Value: false},
				},
			},
		},
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "inuse", Value: true}}}}
	result, err := s.DB.Collection(SCOOTER_COLLECTION).UpdateOne(ctx, query, update)
	if err != nil {
		return false, fmt.Errorf("ReserveScooter: err= %w", err)
	}

	if result.ModifiedCount == 1 {
		return true, nil
	} else if result.ModifiedCount == 0 {
		return false, nil
	} else {
		return false, apperror.NewAppError(
			apperror.WithStatusCode(http.StatusInternalServerError),
			apperror.WithError(fmt.Errorf(fmt.Sprintf("ReserverScooter: %d invaliad number of scooter are resevred!", result.ModifiedCount))),
		)
	}

}

func (s *ScooterReserver) ReleaseScooter(ctx context.Context, scooterID string) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "inuse", Value: false}}}}
	_, err := s.DB.Collection(SCOOTER_COLLECTION).UpdateByID(ctx, scooterID, update)
	if err != nil {
		return fmt.Errorf("ReleaseScooter: err= %w", err)
	}

	return nil
}
