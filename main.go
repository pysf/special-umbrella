package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/db"
	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/seeder"
	"github.com/pysf/special-umbrella/internal/server"
	"github.com/pysf/special-umbrella/internal/simulator"
)

func main() {

	client, err := db.CreateConnection()
	if err != nil {
		panic(err)
	}
	DB := client.Database(config.GetConfig("MONGODB_DATABASE"))

	scooterReserver, err := scooter.NewScooterReserver(DB)
	if err != nil {
		panic(err)
	}

	statusUpdater, err := scooter.NewStatusUpdater(scooterReserver, DB)
	if err != nil {
		panic(err)
	}

	scooterFinder, err := scooter.NewScooterFinder(DB)
	if err != nil {
		panic(err)
	}

	scooterCreator, err := scooter.NewScooterCreator(DB)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	seeder.Start(ctx, scooterCreator, statusUpdater,
		seeder.WithCount(2),
		seeder.WithDistanceShift(1),
		seeder.WithStartDelay(3*time.Second),
	)

	simulator.Start(ctx)

	server, err := server.NewServer(statusUpdater, scooterFinder, scooterReserver)
	if err != nil {
		fmt.Printf("Failde to initiate server! err=%v \n", err)
		panic(err)
	}

	if err := server.Start(); err != nil {
		fmt.Printf("Failed to start server! err=%v \n", err)
		panic(err)
	}
}
