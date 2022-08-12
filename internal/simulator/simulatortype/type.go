package simulatortype

type UpdateScooterStatusRequestBody struct {
	ScooterID string `json:"scooterID"`
	Timestamp string `json:"timestamp"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	EventType string `json:"eventType"`
}
