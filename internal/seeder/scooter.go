package seeder

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/scooter/scooteriface"
	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
)

type ScooterDataSeederOptions func(s *ScooterDataSeeder)

type ScooterDataSeeder struct {
	baseLat           float64
	baseLng           float64
	scootersPerCircle int
	numberOfScooters  int
	startDelay        time.Duration
	distanceShift     float64
	scooterCreator    scooteriface.ScooterCreator
	statusUpdater     scooteriface.StatusUpdater
	ctx               context.Context
}

func Start(ctx context.Context, scooterCreator scooteriface.ScooterCreator, statusUpdater scooteriface.StatusUpdater, options ...ScooterDataSeederOptions) {

	seeder := &ScooterDataSeeder{
		baseLat:           52.519511, //Berlin center
		baseLng:           13.403725,
		ctx:               ctx,
		scootersPerCircle: 10,
		numberOfScooters:  100,
		distanceShift:     0.1,
		startDelay:        3,
		scooterCreator:    scooterCreator,
		statusUpdater:     statusUpdater,
	}

	for _, o := range options {
		o(seeder)
	}

	go func() {
		time.Sleep(seeder.startDelay)
		if err := seeder.addRandomScooters(); err != nil {
			fmt.Printf("Start: failed to add scooters err=%v", err)
		}
	}()

}

func WithLat(lat float64) ScooterDataSeederOptions {
	return func(s *ScooterDataSeeder) {
		s.baseLat = lat
	}
}

func WithLng(lng float64) ScooterDataSeederOptions {
	return func(s *ScooterDataSeeder) {
		s.baseLng = lng
	}
}

func WithDistanceShift(distance float64) ScooterDataSeederOptions {
	return func(s *ScooterDataSeeder) {
		s.distanceShift = distance
	}
}

func WithCount(c int) ScooterDataSeederOptions {
	return func(s *ScooterDataSeeder) {
		s.numberOfScooters = c
	}
}

func WithScooterPerCircle(spc int) ScooterDataSeederOptions {
	return func(s *ScooterDataSeeder) {
		s.scootersPerCircle = spc
	}
}

func WithStartDelay(delay time.Duration) ScooterDataSeederOptions {
	return func(s *ScooterDataSeeder) {
		s.startDelay = delay
	}
}

func (s *ScooterDataSeeder) addRandomScooters() error {

	distance := s.distanceShift

	var counter int
	for {

		for i := 0; i < s.scootersPerCircle; i++ {

			// Add a new scooter
			scooterID, err := CreateRandomScooter(s)
			if err != nil {
				return err
			}

			// Add a new scooter init event
			lat, lng := ShiftLocation(s.baseLat, s.baseLng, distance, (360/rand.Float64() + 1))
			s.statusUpdater.UpdateStatus(s.ctx, struct {
				ScooterID string
				Timestamp string
				Latitude  string
				Longitude string
				EventType string
			}{
				ScooterID: *scooterID,
				EventType: scooter.EVENT_PERIODIC_UPDATE,
				Timestamp: time.Now().Format(time.RFC3339),
				Latitude:  fmt.Sprintf("%f", lat),
				Longitude: fmt.Sprintf("%f", lng),
			})

			counter = counter + 1
			if counter >= s.numberOfScooters {
				return nil
			}
		}
		distance = distance + .2
	}

}

func CreateRandomScooter(s *ScooterDataSeeder) (*string, error) {
	scooterID := uuid.New().String()
	if err := s.scooterCreator.Create(s.ctx, scootertype.Scooter{
		ID:    scooterID,
		InUse: false,
	}); err != nil {
		return nil, err
	}
	return &scooterID, nil
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
