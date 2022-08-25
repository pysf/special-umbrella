package simulator

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/pysf/special-umbrella/internal/client"
	"github.com/pysf/special-umbrella/internal/client/clienttype"
	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/scooter"
)

type Simulator struct {
	apiClient    *client.APIClient
	bottomLeft   clienttype.Location
	topRight     clienttype.Location
	ctx          context.Context
	numberOfBots int
	startDelay   time.Duration
}

func NewSimulator(ctx context.Context) *Simulator {

	apiClient := client.NewAPIClient(config.AppConf().APIBaseURL, config.AppConf().JWTToken)
	bottomLeft := [2]float64{config.AppConf().SimulatorBottomLeftLat, config.AppConf().SimulatorBottomLeftLng}
	topRight := [2]float64{config.AppConf().SimulatorTopRightLat, config.AppConf().SimulatorTopRightLng}
	numberOfBots := config.AppConf().SimulatorBotCounts
	startDelay := time.Duration(config.AppConf().SimulatorStartDelay) * time.Second

	return &Simulator{
		apiClient:    apiClient,
		bottomLeft:   bottomLeft,
		topRight:     topRight,
		ctx:          ctx,
		numberOfBots: numberOfBots,
		startDelay:   startDelay,
	}

}

func (s *Simulator) Start() {

	go func() {

		time.Sleep(s.startDelay)
		log.Println("Simulator started...")

		searchResult, err := s.apiClient.FindScooters(s.bottomLeft, s.topRight, "available")
		if err != nil {
			fmt.Println(fmt.Errorf("StartSimulator: failed to fetch available scooters err= %w", err))
			return
		}
		scooters := searchResult.Scooters
		log.Printf("%v scooters are available \n", len(scooters))

		var wg sync.WaitGroup

		for _, sc := range scooters {
			wg.Add(1)
			go func(sc clienttype.Scooter) {
				if err := s.RunClientBot(sc); err != nil {
					fmt.Println(fmt.Errorf("StartSimulator: err= %w", err))
				}
				wg.Done()
			}(sc)
		}

		wg.Wait()

	}()

}

func (s *Simulator) RunClientBot(sc clienttype.Scooter) error {

	reserved, err := s.apiClient.ReserverScooter(clienttype.ReserveScooterRequestBody{
		ScooterID: sc.ScooterID,
	})

	if err != nil {
		return fmt.Errorf("StartClientBot: failed to reserve the scooter err= %w", err)
	}

	if !reserved {
		return fmt.Errorf(fmt.Sprintf("StartClientBot: %v scooter is already inuse\n", sc.ScooterID))
	}

	nextLocation := sc.Location

	// Send trip start event
	if err := s.apiClient.PublishScooterStatus(clienttype.UpdateScooterRequestBody{
		ScooterID: sc.ScooterID,
		Timestamp: time.Now().Format(time.RFC3339),
		Latitude:  fmt.Sprintf("%f", sc.Location[0]),
		Longitude: fmt.Sprintf("%f", sc.Location[1]),
		EventType: scooter.EVENT_TRIP_STARTED,
	}); err != nil {
		return err
	}
	time.Sleep(3 * time.Second)

	for i := 0; i < 3; i++ {
		nextLocation = ride(nextLocation, 1, float64((60)*i))
		if err := s.apiClient.PublishScooterStatus(clienttype.UpdateScooterRequestBody{
			ScooterID: sc.ScooterID,
			Timestamp: time.Now().Format(time.RFC3339),
			Latitude:  fmt.Sprintf("%f", nextLocation[0]),
			Longitude: fmt.Sprintf("%f", nextLocation[1]),
			EventType: scooter.EVENT_TRIP_UPDATED,
		}); err != nil {
			return err
		}

		time.Sleep(3 * time.Second)
	}

	nextLocation = ride(nextLocation, 1, 60)
	// Send trip end event
	if err := s.apiClient.PublishScooterStatus(clienttype.UpdateScooterRequestBody{
		ScooterID: sc.ScooterID,
		Timestamp: time.Now().Format(time.RFC3339),
		Latitude:  fmt.Sprintf("%f", nextLocation[0]),
		Longitude: fmt.Sprintf("%f", nextLocation[1]),
		EventType: scooter.EVENT_TRIP_ENDED,
	}); err != nil {
		return err
	}
	time.Sleep(3 * time.Second)

	// Send scooter periodic update
	if err := s.apiClient.PublishScooterStatus(clienttype.UpdateScooterRequestBody{
		ScooterID: sc.ScooterID,
		Timestamp: time.Now().Format(time.RFC3339),
		Latitude:  fmt.Sprintf("%f", nextLocation[0]),
		Longitude: fmt.Sprintf("%f", nextLocation[1]),
		EventType: scooter.EVENT_PERIODIC_UPDATE,
	}); err != nil {
		return err
	}

	return nil

}

func ride(location clienttype.Location, distance, bearing float64) clienttype.Location {

	R := 6378.1                        // Radius of the Earth
	brng := bearing * math.Pi / 180    // Convert bearing to radian
	lat := location[0] * math.Pi / 180 // Current coords to radians
	lon := location[1] * math.Pi / 180
	// Do the math magic

	lat = math.Asin(math.Sin(lat)*math.Cos(distance/R) + math.Cos(lat)*math.Sin(distance/R)*math.Cos(brng))
	lon += math.Atan2(math.Sin(brng)*math.Sin(distance/R)*math.Cos(lat), math.Cos(distance/R)-math.Sin(lat)*math.Sin(lat))

	return clienttype.Location{
		(lat * 180 / math.Pi),
		(lon * 180 / math.Pi),
	}

}
