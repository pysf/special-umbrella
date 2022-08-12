package config

import (
	"fmt"
	"os"
)

func GetConfig(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%v is required evronment variable", key))
	}
	return v
}

func LookupConfig(key string) (string, bool) {
	return os.LookupEnv(key)
}
