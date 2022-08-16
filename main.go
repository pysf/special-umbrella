package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/server"
	"github.com/pysf/special-umbrella/internal/simulator"
)

func main() {

	scooterCreator, err := scooter.NewScooterCreator()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	simulator.NewScooterSimulator(
		simulator.WithCount(5),
		simulator.WithDistanceShift(1),
		simulator.WithStartDelay(3),
		simulator.WithJWTToken(os.Getenv("JWT_TOKEN")),
		simulator.WithScooterCreator(scooterCreator),
	).Start(ctx)
	//todo: fix cancel
	defer cancel()

	server, err := server.NewServer()
	if err != nil {
		fmt.Printf("Failde to initiate server! err=%v", err)
		panic(err)
	}

	if err := server.Start(); err != nil {
		fmt.Printf("Failed to start server! err=%v", err)
		panic(err)
	}
}
