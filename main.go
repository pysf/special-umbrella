package main

import (
	"log"

	"github.com/pysf/special-umbrella/internal/server"
)

func main() {

	server, err := server.NewServer()
	if err != nil {
		log.Fatalf("Server failde to start %s", err)
	}
	server.Start()
}
