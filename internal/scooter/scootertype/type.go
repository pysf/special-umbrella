package scootertype

import "time"

// type Location struct {
// 	Latitude  float64
// 	Longitude float64
// }

type Location [2]float64

type ScooterStatusUpdaterInput struct {
	ScooterID string
	Timestamp string
	Latitude  string
	Longitude string
	EventType string
}

type ScooterStatus struct {
	ID        string    `bson:"_id,omitempty"`
	ScooterID string    `bson:"scooterID"`
	Status    string    `bson:"status"`
	Timestamp time.Time `bson:"timestamp"`
	Location  GeoJSON   `bson:"location"`
	EventType string    `bson:"eventType"`
}

type GeoJSON struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type RectangularQueryResult struct {
	Scooters []Scooter `json:"scooters"`
}

type Scooter struct {
	ScooterID  string    `bson:"_id" json:"scooterID"`
	Status     string    `bson:"status" json:"status"`
	LastStatus time.Time `bson:"timestamp" json:"lastStatus"`
	Location   Location  `bson:"location" json:"location"`
}
