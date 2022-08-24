package server_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pysf/special-umbrella/internal/client/clienttype"
	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
	"github.com/pysf/special-umbrella/internal/server"
	"github.com/pysf/special-umbrella/internal/testutils"
)

func TestServer_FindScooter(t *testing.T) {

	DB := testutils.GetDBConnection(t)
	testutils.PrepareTestDatabase(DB, t)

	server, err := server.NewServer(DB)
	if err != nil {
		t.Fatal(err)
	}

	testserver := httptest.NewServer(server.HttpHandler)
	defer testserver.Close()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/api/scooter/search", testserver.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	bl, err := json.Marshal([]float64{testutils.RecBottomLeftLat, testutils.BerlinCenterLng})
	if err != nil {
		t.Fatal(err)
	}

	tr, err := json.Marshal([]float64{testutils.RecTopRightLat, testutils.RecTopRightLng})
	if err != nil {
		t.Fatal(err)
	}

	q.Add("bottomLeft", string(bl))
	q.Add("topRight", string(tr))
	q.Add("status", "available")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", os.Getenv("JWT_TOKEN"))

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Server.FindScooter() error = %v, want %v", err, http.StatusOK)
	}

	var result clienttype.ScooterScearchResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatal(err)
	}

	if len(result.Scooters) == 0 {
		t.Errorf("Server.FindScooter() error = number of returned scooters can not be 0")
	}

	for _, scooter := range result.Scooters {
		if scooter.ScooterID == "" {
			t.Errorf("Server.FindScooter() scooterID can not be empty = %v", scooter.ScooterID)
		}

		if scooter.Status != "available" {
			t.Errorf("Server.FindScooter() incorrect scooter status error = %v, want %v", scooter.Status, "available")
		}

		if !testutils.IsInRectangle(testutils.BottomLeft, testutils.TopRight, scootertype.Location(scooter.Location), t) {
			t.Errorf("Server.FindScooter() error = %v locatin is not in the rectangle bottomLeft= %v topRight= %v", scooter.Location, testutils.BottomLeft, testutils.TopRight)
		}
	}
}
