package main

import (
	"fmt"

	"github.com/pysf/special-umbrella/internal/server"
)

func main() {

	//todo start the simulator

	server, err := server.NewServer()
	if err != nil {
		fmt.Println("Server failde to start")
		panic(err)
	}
	server.Start()
}
