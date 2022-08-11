package scooter

type Location struct {
	Latitude  float64
	Longitude float64
}

type ScooteStatusEvent struct {
	ScooterID string `json:"scooterID"`
	Timestamp string `json:"timestamp"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	EventType string `json:"eventType"`
}
