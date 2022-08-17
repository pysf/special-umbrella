package config

import (
	"fmt"
	"os"
	"strconv"
)

func GetConfig(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%v is required evronment variable", key))
	}
	return v
}

func GetConfigAsInt(key string) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%v is required evronment variable", key))
	}

	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(fmt.Errorf("GetConfigAsInt: err= %w", err))
	}
	return int(n)
}

func GetConfigAsFloat(key string) float64 {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%v is required evronment variable", key))
	}

	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(fmt.Errorf("GetConfigAsFloat: err= %w", err))
	}
	return n
}

func LookupConfig(key string) (string, bool) {
	return os.LookupEnv(key)
}
