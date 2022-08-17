package scooter

import (
	"context"
	"fmt"

	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScooterFinder struct {
	DB *mongo.Database
}

func NewScooterFinder(DB *mongo.Database) (*ScooterFinder, error) {

	return &ScooterFinder{
		DB: DB,
	}, nil
}

func (s *ScooterFinder) RectangularQuery(ctx context.Context, q struct {
	Status     string
	BottomLeft scootertype.Location
	TopRight   scootertype.Location
}) (*scootertype.RectangularQueryResult, error) {

	filterByLocation := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "location", Value: bson.D{
				{Key: "$geoWithin", Value: bson.D{
					{Key: "$box", Value: bson.A{
						bson.A{
							q.BottomLeft[0],
							q.BottomLeft[1],
						},
						bson.A{
							q.TopRight[0],
							q.TopRight[1],
						},
					},
					},
				},
				},
			},
			},
		},
		},
	}

	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "timestamp", Value: -1}}}}
	groupByStage := bson.D{
		{Key: "$group",
			Value: bson.D{
				{Key: "_id",
					Value: "$scooterID"},
				{Key: "status", Value: bson.D{{Key: "$first", Value: "$status"}}},
				{Key: "location", Value: bson.D{{Key: "$first", Value: "$location.coordinates"}}},
				{Key: "timestamp", Value: bson.D{{Key: "$first", Value: "$timestamp"}}},
			},
		},
	}

	filterByStatus := bson.D{
		{
			Key: "$match", Value: bson.D{
				{Key: "status", Value: bson.D{{Key: "$eq", Value: q.Status}}},
			},
		},
	}

	//Todo: add limit and offset to enable pagination
	cursor, err := s.DB.Collection(SCOOTER_STATUS_COLLECTION).Aggregate(ctx, mongo.Pipeline{filterByLocation, sortStage, groupByStage, filterByStatus})
	if err != nil {
		return nil, fmt.Errorf("RectangularQuery: err= %w", err)
	}

	result := scootertype.RectangularQueryResult{
		Scooters: []scootertype.ScooterAggregationItems{},
	}
	if err := cursor.All(ctx, &result.Scooters); err != nil {
		return nil, fmt.Errorf("RectangularQuery: faild to read quecy result, err= %w", err)
	}

	return &result, nil
}
