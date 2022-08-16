package scootertype

import "time"

// type Location struct {
// 	Latitude  float64
// 	Longitude float64
// }

type Location [2]float64

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
	Scooters []ScooterAggregationItems `json:"scooters"`
}

type ScooterAggregationItems struct {
	ScooterID  string    `bson:"_id" json:"scooterID"`
	Status     string    `bson:"status" json:"status"`
	LastUpdate time.Time `bson:"timestamp" json:"lastStatus"`
	Location   Location  `bson:"location" json:"location"`
}

type Scooter struct {
	ID    string `bson:"_id" json:"ID"`
	InUse bool   `bson:"inuse" json:"inuse"`
}
