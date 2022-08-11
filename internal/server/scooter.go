package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/apperror"
	"github.com/pysf/special-umbrella/internal/scooter"
)

func (s *Server) UpdateScooterStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("AddScooterStatus: failed to read request body err=%w", err)
	}

	updateReq := &UpdateScooterStatusRequest{}
	if err = json.Unmarshal(b, &updateReq); err != nil {
		return apperror.NewAppError(
			apperror.WithError(err),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	if err = validator.New().Struct(updateReq); err != nil {
		return apperror.NewAppError(
			apperror.WithError(err),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	scooterStatusEvent := scooter.ScooteStatusEvent{
		EventType: updateReq.EventType,
		ScooterID: updateReq.ScooterID,
		Timestamp: updateReq.Timestamp,
		Latitude:  updateReq.Latitude,
		Longitude: updateReq.Longitude,
	}

	_, err = s.StatusUpdater.UpdateStatus(r.Context(), scooterStatusEvent)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	return nil
}

type UpdateScooterStatusRequest struct {
	EventType string `json:"eventType" validate:"required,ascii"`
	ScooterID string `json:"scooterID" validate:"required,uuid4"`
	Latitude  string `json:"latitude"  validate:"required,latitude"`
	Longitude string `json:"longitude" validate:"required,longitude"`
	Timestamp string `json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05Z"`
}
