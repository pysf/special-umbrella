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
	"github.com/pysf/special-umbrella/internal/simulator/simulatortype"
)

type Simulator struct {
	apiClient    *client.APIClient
	bottomLeft   clienttype.Location
	topRight     clienttype.Location
	ctx          context.Context
	numberOfBots int
	startDelay   time.Duration
}

func Start(ctx context.Context) {

	apiClient := client.NewAPIClient(config.GetConfig("BASE_URL"), config.GetConfig("JWT_TOKEN"))
	bottomLeft := [2]float64{config.GetConfigAsFloat("SIMULATOR_BOTTOM_LEFT_LAT"), config.GetConfigAsFloat("SIMULATOR_BOTTOM_LEFT_LNG")}
	topRight := [2]float64{config.GetConfigAsFloat("SIMULATOR_TOP_RIGHT_LAT"), config.GetConfigAsFloat("SIMULATOR_TOP_RIGHT_LNG")}
	numberOfBots := config.GetConfigAsInt("SIMULATOR_BOT_COUNTS")
	startDelay := time.Duration(config.GetConfigAsInt("SIMULATOR_START_DELAY")) * time.Second

	simulator := &Simulator{
		apiClient:    apiClient,
		bottomLeft:   bottomLeft,
		topRight:     topRight,
		ctx:          ctx,
		numberOfBots: numberOfBots,
		startDelay:   startDelay,
	}

	go func() {
		botResultsCh := simulator.RunBots()
		for {
			select {
			case <-ctx.Done():
				return
			case sr, ok := <-botResultsCh:
				if !ok {
					return
				}
				if sr.Err != nil {
					log.Printf("Simulator Err= %v \n", sr.Err)
					return
				}
				log.Printf("Simulator: %v \n", sr.Message)
			}
		}

	}()

}

func (s *Simulator) RunBots() chan simulatortype.SimulatorResponse {

	resultChan := make(chan simulatortype.SimulatorResponse)

	go func() {
		defer close(resultChan)
		time.Sleep(s.startDelay)
		log.Println("Simulator started...")

		searchResult, err := s.apiClient.FindScooters(s.bottomLeft, s.topRight, "available")
		if err != nil {
			s.SendResult(resultChan, simulatortype.SimulatorResponse{
				Err: fmt.Errorf("StartSimulator: failed to fetch available scooters err= %w", err),
			})
			return
		}
		scooters := searchResult.Scooters
		log.Printf("%v scooters are available \n", len(scooters))

		scooterCh := make(chan clienttype.Scooter, len(scooters))
		defer close(scooterCh)

		botResultCh := make([]<-chan simulatortype.SimulatorResponse, s.numberOfBots)
		for i := 0; i < s.numberOfBots; i++ {
			botResultCh[i] = s.ClientBot(scooterCh)
		}

		for _, scooter := range scooters {
			select {
			case <-s.ctx.Done():
			case scooterCh <- scooter:
			}
		}

		for r := range s.Fanin(botResultCh...) {
			select {
			case <-s.ctx.Done():
			case resultChan <- r:
			}
		}

	}()

	return resultChan
}

func (s *Simulator) ClientBot(scooterCh <-chan clienttype.Scooter) <-chan simulatortype.SimulatorResponse {

	resultCh := make(chan simulatortype.SimulatorResponse)

	go func() {
		defer close(resultCh)

		for {

			select {
			case <-s.ctx.Done():
				return
			case sc, ok := <-scooterCh:
				if !ok {
					return
				}

				reserved, err := s.apiClient.ReserverScooter(clienttype.ReserveScooterRequestBody{
					ScooterID: sc.ScooterID,
				})
				if err != nil {
					s.SendResult(resultCh, simulatortype.SimulatorResponse{
						Err: fmt.Errorf("StartClientBot: failed to reserve the scooter err= %w", err),
					})
					return
				}

				if !reserved {
					s.SendResult(resultCh, simulatortype.SimulatorResponse{
						Timestamp: time.Now(),
						Message:   fmt.Sprintf("%v scooter is already inuse\n", sc.ScooterID),
					})
					return
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
					s.SendResult(resultCh, simulatortype.SimulatorResponse{
						Err: err,
					})
					return
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
						s.SendResult(resultCh, simulatortype.SimulatorResponse{
							Err: err,
						})
						return
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
					s.SendResult(resultCh, simulatortype.SimulatorResponse{
						Err: err,
					})
					return
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
					s.SendResult(resultCh, simulatortype.SimulatorResponse{
						Err: err,
					})
					return
				}

				s.SendResult(resultCh, simulatortype.SimulatorResponse{
					Timestamp: time.Now(),
					Message:   fmt.Sprintf("%v scooter trip ended\n", sc.ScooterID),
				})

			}
		}

	}()

	return resultCh
}

func (s *Simulator) Fanin(chans ...<-chan simulatortype.SimulatorResponse) <-chan simulatortype.SimulatorResponse {

	var wg sync.WaitGroup
	multiplexedCh := make(chan simulatortype.SimulatorResponse)

	multiplexer := func(ch <-chan simulatortype.SimulatorResponse) {
		defer wg.Done()
		for d := range ch {
			select {
			case <-s.ctx.Done():
				return
			case multiplexedCh <- d:
			}
		}
	}

	wg.Add(len(chans))

	for _, ch := range chans {
		go multiplexer(ch)
	}

	go func() {
		defer close(multiplexedCh)
		wg.Wait()
	}()

	return multiplexedCh
}

func (s *Simulator) SendResult(ch chan<- simulatortype.SimulatorResponse, r simulatortype.SimulatorResponse) {
	select {
	case <-s.ctx.Done():
	case ch <- r:
	}
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
