package scooter_test

import (
	"context"
	"testing"

	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
	"github.com/pysf/special-umbrella/internal/testutils"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestScooterFinder_RectangularQuery(t *testing.T) {
	DB := testutils.GetDBConnection(t)
	testutils.PrepareTestDatabase(DB, t)

	type fields struct {
		DB *mongo.Database
	}

	// Small Rectangel in center of Berlin

	type args struct {
		ctx context.Context
		q   struct {
			Status     string
			BottomLeft scootertype.Location
			TopRight   scootertype.Location
		}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *scootertype.RectangularQueryResult
		wantErr bool
	}{
		{
			name: "Test if returned scooters are in the rectangle!",
			fields: fields{
				DB: DB,
			},
			args: args{
				ctx: context.TODO(),
				q: struct {
					Status     string
					BottomLeft scootertype.Location
					TopRight   scootertype.Location
				}{
					Status:     "available",
					BottomLeft: scootertype.Location{testutils.RecBottomLeftLat, testutils.RecBottomLeftLng},
					TopRight:   scootertype.Location{testutils.RecTopRightLat, testutils.RecTopRightLng},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scooter.ScooterFinder{
				DB: tt.fields.DB,
			}
			got, err := s.RectangularQuery(tt.args.ctx, tt.args.q)

			if (err != nil) != tt.wantErr {
				t.Errorf("ScooterFinder.RectangularQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, scooter := range got.Scooters {
				if !testutils.IsInRectangle(testutils.BottomLeft, testutils.TopRight, scooter.Location, t) {
					t.Error("ScooterFinder.RectangularQuery() err= Returned location is not in rectangle")
				}
			}

		})
	}
}
