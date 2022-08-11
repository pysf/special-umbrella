package main

import (
	"context"
	"fmt"

	"github.com/pysf/special-umbrella/internal/server"
	simulator "github.com/pysf/special-umbrella/internal/simulator/scooter"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	simulator.NewScooterSimulator(
		simulator.WithCount(100),
		simulator.WithDistanceShift(1),
		simulator.WithStartDelay(3),
	).Start(ctx)
	defer cancel()

	server, err := server.NewServer()
	if err != nil {
		fmt.Println("Server failde to start")
		panic(err)
	}
	server.Start()
}
