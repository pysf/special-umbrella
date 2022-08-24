package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/db"
	"github.com/pysf/special-umbrella/internal/seeder"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetDBConnection(t *testing.T) *mongo.Database {

	client, err := db.CreateConnection(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	return client.Database(config.AppConf().MongoDatabase)
}

func PrepareTestDatabase(DB *mongo.Database, t *testing.T) {

	seeder.NewScooterDataSeeder(context.TODO(), DB,
		seeder.WithNumberOfInitialScooters(200),
		seeder.WithDistanceShift(1),
		seeder.WithStartDelay(0),
		seeder.WithLat(BerlinCenterLat),
		seeder.WithLng(BerlinCenterLng),
		seeder.WithScooterPerCircle(18),
	).Start()

	time.Sleep(3 * time.Second)

}
