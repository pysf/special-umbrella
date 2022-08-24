package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Configuration struct {
	JWTToken      string `env:"JWT_TOKEN" env-required:"true"`
	JWTTokenKey   string `env:"JWT_TOKEN_KEY" env-required:"true"`
	APIBaseURL    string `env:"BASE_URL" env-required:"true"`
	MongoDatabase string `env:"MONGODB_DATABASE" env-required:"true"`
	MongoURI      string `env:"MONGODB_URI" env-required:"true"`

	SimulatorBotCounts     int     `env:"SIMULATOR_BOT_COUNTS" env-default:"3"`
	SimulatorStartDelay    int     `env:"SIMULATOR_START_DELAY" env-default:"10"`
	SimulatorBottomLeftLat float64 `env:"SIMULATOR_BOTTOM_LEFT_LAT" env-default:"52.415548"`
	SimulatorBottomLeftLng float64 `env:"SIMULATOR_BOTTOM_LEFT_LNG" env-default:"13.216032"`
	SimulatorTopRightLat   float64 `env:"SIMULATOR_TOP_RIGHT_LAT" env-default:"52.621397"`
	SimulatorTopRightLng   float64 `env:"SIMULATOR_TOP_RIGHT_LNG" env-default:"13.597643"`

	SeederInitialScooters int     `env:"SEEDER_INITIAL_SCOOTERS" env-default:"100"`
	SeederDistanceShift   int     `env:"SEEDER_DISTANCE_SHIFT" env-default:"1"`
	SeederStartDelay      int     `env:"SEEDER_START_DELAY" env-default:"10"`
	SeederStartLat        float64 `env:"SEEDER_START_LAT" env-default:"52.520008"`
	SeederStartLng        float64 `env:"SEEDER_START_LNG" env-default:"13.404954"`
}

var once sync.Once

var cfg Configuration

func AppConf() Configuration {

	once.Do(func() {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			fmt.Println(err)
			panic(fmt.Errorf("AppConf: Failed to load config err= %w", err))
		}
	})

	return cfg
}
