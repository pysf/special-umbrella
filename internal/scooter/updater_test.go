package scooter_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/scooter/scooteriface"
	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
	"github.com/pysf/special-umbrella/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestStatusUpdater_UpdateStatus(t *testing.T) {
	DB := testutils.GetDBConnection(t)

	lat := 52.519121
	lng := 13.402381
	type fields struct {
		DB              *mongo.Database
		ScooterReserver scooteriface.ScooterReserver
	}

	type args struct {
		ctx               context.Context
		updateStatusInput struct {
			ScooterID string
			Timestamp string
			Latitude  string
			Longitude string
			EventType string
		}
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *string
		wantErr bool
	}{
		{
			name: "Test if scooter status is saved correctly",
			fields: fields{
				DB: DB,
			},
			args: args{
				ctx: context.TODO(),
				updateStatusInput: struct {
					ScooterID string
					Timestamp string
					Latitude  string
					Longitude string
					EventType string
				}{
					ScooterID: uuid.New().String(),
					Timestamp: time.Now().Format(time.RFC3339),
					Latitude:  fmt.Sprintf("%f", lat),
					Longitude: fmt.Sprintf("%f", lng),
					EventType: "trip-started",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scooter.StatusUpdater{
				DB:              tt.fields.DB,
				ScooterReserver: tt.fields.ScooterReserver,
			}

			id, err := s.UpdateStatus(tt.args.ctx, tt.args.updateStatusInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("StatusUpdater.UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var result scootertype.ScooterStatus

			DB.Collection(scooter.SCOOTER_STATUS_COLLECTION).FindOne(tt.args.ctx, bson.D{{Key: "_id", Value: id}}).Decode(&result)

			if !reflect.DeepEqual(result.EventType, tt.args.updateStatusInput.EventType) {
				t.Errorf("StatusUpdater.UpdateStatus() = %v, want %v", result.EventType, tt.args.updateStatusInput.EventType)
			}
			if !reflect.DeepEqual(result.Status, "inuse") {
				t.Errorf("StatusUpdater.UpdateStatus() = %v, want %v", result.Status, "inuse")
			}
			if !reflect.DeepEqual(fmt.Sprintf("%f", result.Location.Coordinates[0]), tt.args.updateStatusInput.Latitude) {
				t.Errorf("StatusUpdater.UpdateStatus() = %v, want %v", result.Location.Coordinates[0], tt.args.updateStatusInput.Latitude)
			}
			if !reflect.DeepEqual(fmt.Sprintf("%f", result.Location.Coordinates[1]), tt.args.updateStatusInput.Longitude) {
				t.Errorf("StatusUpdater.UpdateStatus() = %v, want %v", result.Location.Coordinates[1], tt.args.updateStatusInput.Longitude)
			}
			if !reflect.DeepEqual(result.Timestamp.Format(time.RFC3339), tt.args.updateStatusInput.Timestamp) {
				t.Errorf("StatusUpdater.UpdateStatus() = %v, want %v", result.Timestamp.Format(time.RFC3339), tt.args.updateStatusInput.Timestamp)
			}
		})
	}
}
