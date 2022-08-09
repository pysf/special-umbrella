package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/model"
)

func (s *Server) AddScooterStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("AddScooterStatus: failed to read request body err=%w", err)
	}

	addScooterStatusRequest := &AddScooterStatusRequest{}
	if err = json.Unmarshal(b, &addScooterStatusRequest); err != nil {
		return NewHttpError(nil, err.Error(), http.StatusBadRequest)
	}

	if err = validator.New().Struct(addScooterStatusRequest); err != nil {
		return NewHttpError(nil, err.Error(), http.StatusBadRequest)
	}

	timestamp, err := time.Parse(time.RFC3339, addScooterStatusRequest.Timestamp)
	if err != nil {
		return NewHttpError(nil, err.Error(), http.StatusBadRequest)
	}

	location := model.Location{}
	if err != fillLocation(addScooterStatusRequest.Latitude, addScooterStatusRequest.Longitude, &location) {
		return NewHttpError(nil, err.Error(), http.StatusBadRequest)
	}

	scooterStatusEvent := model.ScooteStatusEvent{
		ScooterID: addScooterStatusRequest.ScooterID,
		Timestamp: timestamp,
		Location:  location,
	}
	s.ScooterStatusRepository.AddStatus(r.Context(), scooterStatusEvent)

	return nil
}

type AddScooterStatusRequest struct {
	ScooterID string `json:"scooterID" validate:"required,uuid4"`
	Latitude  string `json:"latitude"  validate:"required,latitude"`
	Longitude string `json:"longitude" validate:"required,longitude"`
	Timestamp string `json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05Z"`
}

func fillLocation(lat, lng string, l *model.Location) error {
	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return fmt.Errorf("FillLocation: latitude err=%w", err)
	}

	longitude, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		return fmt.Errorf("FillLocation: longitude err=%w", err)
	}

	l.Latitude = latitude
	l.Longitude = longitude

	return nil
}
