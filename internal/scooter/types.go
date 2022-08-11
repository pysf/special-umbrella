package scooter

type Location struct {
	Latitude  float64
	Longitude float64
}

type ScooteStatusEvent struct {
	ScooterID string
	Timestamp string
	Latitude  string
	Longitude string
	EventType string
}
