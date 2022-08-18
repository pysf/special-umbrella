package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/db"
	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/seeder"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetDBConnection(t *testing.T) *mongo.Database {

	client, err := db.CreateConnection(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	return client.Database(config.GetConfig("MONGODB_DATABASE"))

}

func PrepareTestDatabase(DB *mongo.Database, t *testing.T) {

	ctx := context.TODO()

	scooterReserver, err := scooter.NewScooterReserver(DB)
	if err != nil {
		t.Fatal(err)
	}

	statusUpdater, err := scooter.NewStatusUpdater(scooterReserver, DB)
	if err != nil {
		t.Fatal(err)
	}

	scooterCreator, err := scooter.NewScooterCreator(DB)
	if err != nil {
		t.Fatal(err)
	}

	seeder.Start(ctx, scooterCreator, statusUpdater,
		seeder.WithCount(200),
		seeder.WithDistanceShift(1),
		seeder.WithStartDelay(0),
		seeder.WithLat(BerlinCenterLat),
		seeder.WithLng(BerlinCenterLng),
		seeder.WithScooterPerCircle(18),
	)

	time.Sleep(3 * time.Second)

}
