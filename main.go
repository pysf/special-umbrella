package main

import (
	"fmt"

	"github.com/pysf/special-umbrella/internal/db"
)

func main() {
	fmt.Println("Hello Nord!")
	_, err := db.CreateConnection()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}
