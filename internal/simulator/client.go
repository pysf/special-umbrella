package simulator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
)

type ClientSimulator struct {
	updateIntervals time.Duration
	tripDuration    time.Duration
	startDelay      int
	distanceShift   float64
	apiURL          string
	jwtToken        string
}

func (c *ClientSimulator) FindFreeScooters() error {
	req, err := http.NewRequest(http.MethodGet, c.apiURL, nil)
	if err != nil {
		return fmt.Errorf("FindFreeScooters create new req failed err=%w", err)
	}

	var bottomLeft scootertype.Location = [2]float64{}
	bl, err := json.Marshal(bottomLeft)
	if err != nil {
		return fmt.Errorf("FindFreeScooters json marshal err=%w", err)
	}
	req.URL.Query().Add("bottomLeft", string(bl))

	var topRight scootertype.Location = [2]float64{}
	tr, err := json.Marshal(topRight)
	if err != nil {
		return fmt.Errorf("FindFreeScooters json marshal err=%w", err)
	}
	req.URL.Query().Add("bottomLeft", string(tr))
	req.Header.Add("Authorization", c.jwtToken)

	client := &http.Client{}
	result, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("FindFreeScooters http client err=%w", err)
	}

	if result.StatusCode != http.StatusOK {
		return fmt.Errorf("FindFreeScooters api err=%v %v", result.Status, result.StatusCode)
	}
	return nil
}

func StartTrip(scooterID string) {

	go func() {

	}()
}
