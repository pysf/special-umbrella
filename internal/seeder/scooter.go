package seeder

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/scooter/scooteriface"
	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScooterDataSeederOptions func(s *ScooterDataSeeder)

type ScooterDataSeeder struct {
	scootersPerCircle int
	numberOfScooters  int
	baseLat           float64
	baseLng           float64
	distanceShift     float64
	startDelay        time.Duration
	scooterCreator    scooteriface.ScooterCreator
	statusUpdater     scooteriface.StatusUpdater
	ctx               context.Context
}

func NewScooterDataSeeder(ctx context.Context, DB *mongo.Database, options ...ScooterDataSeederOptions) *ScooterDataSeeder {
	scooterReserver, err := scooter.NewScooterReserver(DB)
	if err != nil {
		panic(err)
	}

	statusUpdater, err := scooter.NewStatusUpdater(scooterReserver, DB)
	if err != nil {
		panic(err)
	}

	scooterCreator, err := scooter.NewScooterCreator(DB)
	if err != nil {
		panic(err)
	}
	seeder := &ScooterDataSeeder{
		ctx:            ctx,
		statusUpdater:  statusUpdater,
		scooterCreator: scooterCreator,
	}

	for _, o := range options {
		o(seeder)
	}

	return seeder
}

func (s *ScooterDataSeeder) Start() {

	go func() {
		time.Sleep(s.startDelay)
		if err := s.addRandomScooters(); err != nil {
			log.Printf("Start: failed to add scooters err=%v", err)
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

func WithNumberOfInitialScooters(c int) ScooterDataSeederOptions {
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

		for i := 1; i < s.scootersPerCircle; i++ {

			// Add a new scooter
			scooterID, err := CreateRandomScooter(s)
			if err != nil {
				return err
			}

			// Add a new scooter init event
			lat, lng := ShiftLocation(s.baseLat, s.baseLng, distance, (20*i)%360)
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

func ShiftLocation(latitude, longitude float64, distance float64, bearing int) (lat, lng float64) {

	R := 6378.1                              // Radius of the Earth
	brng := float64(bearing) * math.Pi / 180 // Convert bearing to radian
	lat = latitude * math.Pi / 180           // Current coords to radians
	lng = longitude * math.Pi / 180

	// Do the math magic
	lat = math.Asin(math.Sin(lat)*math.Cos(distance/R) + math.Cos(lat)*math.Sin(distance/R)*math.Cos(brng))
	lng += math.Atan2(math.Sin(brng)*math.Sin(distance/R)*math.Cos(lat), math.Cos(distance/R)-math.Sin(lat)*math.Sin(lat))

	return (lat * 180 / math.Pi), (lng * 180 / math.Pi)

}
