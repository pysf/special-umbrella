package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pysf/special-umbrella/internal/seeder"
	"github.com/pysf/special-umbrella/internal/server"
	"github.com/pysf/special-umbrella/internal/simulator"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	seeder.Start(ctx,
		seeder.WithCount(2),
		seeder.WithDistanceShift(1),
		seeder.WithStartDelay(3*time.Second),
	)

	simulator.Start(ctx)

	server, err := server.NewServer()
	if err != nil {
		fmt.Printf("Failde to initiate server! err=%v \n", err)
		panic(err)
	}

	if err := server.Start(); err != nil {
		fmt.Printf("Failed to start server! err=%v \n", err)
		panic(err)
	}
}
