package utils

import (
	"context"
	"testing"
	"time"

	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/db"
	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/seeder"
)

func PrepareTestDatabase(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := db.CreateConnection()
	if err != nil {
		t.Fatal(err)
	}
	DB := client.Database(config.GetConfig("MONGODB_DATABASE"))

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
		seeder.WithCount(100),
		seeder.WithDistanceShift(1),
		seeder.WithStartDelay(0),
		seeder.WithLat(52.520091),
		seeder.WithLng(13.393079),
	)

	time.Sleep(3 * time.Second)
}
