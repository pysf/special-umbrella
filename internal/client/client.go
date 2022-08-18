package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pysf/special-umbrella/internal/client/clienttype"
)

type APIClient struct {
	apiBaseURL string
	jwtToken   string
}

func NewAPIClient(baseURL string, jwtToken string) *APIClient {
	return &APIClient{
		apiBaseURL: baseURL,
		jwtToken:   jwtToken,
	}
}

func (c *APIClient) FindScooters(bottomLeft, topRight [2]float64, status string) (*clienttype.ScooterScearchResult, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/api/scooter/search", c.apiBaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("FindScooters: create new req failed err=%w", err)
	}

	q := req.URL.Query()
	bl, err := json.Marshal(bottomLeft)
	if err != nil {
		return nil, fmt.Errorf("FindScooters: json marshal err=%w", err)
	}

	tr, err := json.Marshal(topRight)
	if err != nil {
		return nil, fmt.Errorf("FindScooters: json marshal err=%w", err)
	}

	q.Add("bottomLeft", string(bl))
	q.Add("topRight", string(tr))
	q.Add("status", status)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", c.jwtToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("FindScooters: http client err=%w", err)
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("FindScooters: failed to read response body err= %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FindScooters: api err=%v %v %v", response.Status, response.StatusCode, string(b))
	}

	var result clienttype.ScooterScearchResult
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, fmt.Errorf("FinedScooters: failed to parse response json body err= %w", err)
	}

	return &result, nil
}

func (c *APIClient) PublishScooterStatus(statusEvent clienttype.UpdateScooterRequestBody) error {

	b, err := json.Marshal(statusEvent)
	if err != nil {
		return fmt.Errorf("PublishStatusEvents: json marshal err= %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/api/scooter/status", c.apiBaseURL), bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("PublishStatusEvents: NewRequest err= %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", c.jwtToken))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("PublishStatusEvents: send request err= %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("PublishStatusEvents: err= %v", res.Status)
	}

	return nil
}

func (c *APIClient) ReserverScooter(reserveRequest clienttype.ReserveScooterRequestBody) (bool, error) {
	b, err := json.Marshal(reserveRequest)
	if err != nil {
		return false, fmt.Errorf("ReserverScooter: json marshal err= %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/api/scooter/reserve", c.apiBaseURL), bytes.NewBuffer(b))
	if err != nil {
		return false, fmt.Errorf("ReserverScooter: NewRequest err= %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", c.jwtToken))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("ReserverScooter: send request err= %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return true, nil
	} else if res.StatusCode == http.StatusForbidden {
		return false, nil
	} else {
		return false, fmt.Errorf("ReserverScooter: err= %w", err)
	}

}
