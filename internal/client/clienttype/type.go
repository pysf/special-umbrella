package clienttype

type UpdateScooterRequestBody struct {
	ScooterID string `json:"scooterID"`
	Timestamp string `json:"timestamp"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	EventType string `json:"eventType"`
}

type ScooterScearchResult struct {
	Scooters []Scooter `json:"scooters"`
}

type Scooter struct {
	ScooterID  string   `json:"scooterID"`
	Status     string   `json:"status"`
	LastUpdate string   `json:"lastStatus"`
	Location   Location `json:"location"`
}

type ReserveScooterRequestBody struct {
	ScooterID string `json:"id"`
}

type Location [2]float64
