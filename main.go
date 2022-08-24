package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/db"
	"github.com/pysf/special-umbrella/internal/seeder"
	"github.com/pysf/special-umbrella/internal/server"
	"github.com/pysf/special-umbrella/internal/simulator"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := db.CreateConnection(ctx)
	if err != nil {
		panic(err)
	}

	DB := client.Database(config.AppConf().MongoDatabase)

	seeder.NewScooterDataSeeder(ctx, DB,
		seeder.WithNumberOfInitialScooters(config.AppConf().SeederInitialScooters),
		seeder.WithDistanceShift(float64(config.AppConf().SeederDistanceShift)),
		seeder.WithStartDelay(time.Duration(config.AppConf().SeederStartDelay)),
		seeder.WithLat(config.AppConf().SeederStartLat),
		seeder.WithLng(config.AppConf().SeederStartLng),
	).Start()

	simulator.NewSimulator(ctx).Start()

	server, err := server.NewServer(DB)
	if err != nil {
		fmt.Printf("Failde to initiate server! err=%v \n", err)
		panic(err)
	}

	if err := server.Start(); err != nil {
		fmt.Printf("Failed to start server! err=%v \n", err)
		panic(err)
	}
}
