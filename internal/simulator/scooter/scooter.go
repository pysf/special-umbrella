package simulator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/scooter/status"
)

type simulatorOptions func(s *ScooterSimulator)

type ScooterSimulator struct {
	baseLat           float64
	baseLng           float64
	scootersPerCircle int
	numberOfScooters  int
	startDelay        int
	distanceShift     float64
	apiURL            string
	jwtToken          string
}

func NewScooterSimulator(options ...simulatorOptions) *ScooterSimulator {
	statusSimulator := &ScooterSimulator{
		baseLat:           52.519511, //Berlin center
		baseLng:           13.403725,
		scootersPerCircle: 10,
		numberOfScooters:  100,
		distanceShift:     0.1,
		startDelay:        10,
		apiURL:            "http://127.0.0.1:8080/api/scooter/status",
	}

	for _, o := range options {
		o(statusSimulator)
	}

	return statusSimulator
}

func WithLat(lat float64) simulatorOptions {
	return func(s *ScooterSimulator) {
		s.baseLat = lat
	}
}

func WithLng(lng float64) simulatorOptions {
	return func(s *ScooterSimulator) {
		s.baseLng = lng
	}
}

func WithDistanceShift(distance float64) simulatorOptions {
	return func(s *ScooterSimulator) {
		s.distanceShift = distance
	}
}

func WithCount(c int) simulatorOptions {
	return func(s *ScooterSimulator) {
		s.numberOfScooters = c
	}
}

func WithScooterPerCircle(spc int) simulatorOptions {
	return func(s *ScooterSimulator) {
		s.scootersPerCircle = spc
	}
}

func WithStartDelay(delayInSec int) simulatorOptions {
	return func(s *ScooterSimulator) {
		s.startDelay = delayInSec
	}
}

func WithAPIURL(apiURL string) simulatorOptions {
	return func(s *ScooterSimulator) {
		s.apiURL = apiURL
	}
}

func WithJWTToken(jwtToken string) simulatorOptions {
	return func(s *ScooterSimulator) {
		s.jwtToken = jwtToken
	}
}

func (s *ScooterSimulator) Start(ctx context.Context) {
	go func() {
		time.Sleep(time.Second * time.Duration(s.startDelay))
		scooters := s.GenerateRandomScooters()

		for _, sc := range scooters {
			if err := s.Seed(sc); err != nil {
				fmt.Println("Err: ", err)
				// ctx.Err()
			}

		}

	}()
}

func (s *ScooterSimulator) Seed(statusEvent scooter.ScooteStatusEvent) error {

	b, err := json.Marshal(statusEvent)
	if err != nil {
		return fmt.Errorf("Seed: json marshal err= %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.apiURL, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("Seed: NewRequest err= %w", err)
	}
	req.Header.Add("Authorization", s.jwtToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Seed: send request err= %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Seed: err= %v", res.Status)
	}

	return nil
}

func (s *ScooterSimulator) GenerateRandomScooters() []scooter.ScooteStatusEvent {

	events := make([]scooter.ScooteStatusEvent, 0, s.numberOfScooters)
	distance := s.distanceShift

	for {

		for i := 0; i < s.scootersPerCircle; i++ {
			lat, lng := ShiftLocation(s.baseLat, s.baseLng, distance, (360/rand.Float64() + 1))
			events = append(events, scooter.ScooteStatusEvent{
				ScooterID: uuid.New().String(),
				EventType: status.PeriodicUpdate,
				Timestamp: time.Now().Format(time.RFC3339),
				Latitude:  fmt.Sprintf("%f", lat),
				Longitude: fmt.Sprintf("%f", lng),
			})

			if len(events) >= s.numberOfScooters {
				return events
			}
		}
		distance = distance + .2

	}

}

func ShiftLocation(latitude, longitude float64, distance, bearing float64) (lat, lng float64) {

	R := 6378.1                     // Radius of the Earth
	brng := bearing * math.Pi / 180 // Convert bearing to radian
	lat = latitude * math.Pi / 180  // Current coords to radians
	lng = longitude * math.Pi / 180

	// Do the math magic
	lat = math.Asin(math.Sin(lat)*math.Cos(distance/R) + math.Cos(lat)*math.Sin(distance/R)*math.Cos(brng))
	lng += math.Atan2(math.Sin(brng)*math.Sin(distance/R)*math.Cos(lat), math.Cos(distance/R)-math.Sin(lat)*math.Sin(lat))

	return (lat * 180 / math.Pi), (lng * 180 / math.Pi)

}
