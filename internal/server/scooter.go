package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/apperror"
	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
	"github.com/pysf/special-umbrella/internal/server/servertype"
)

func (s *Server) UpdateScooterStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("AddScooterStatus: failed to read request body err=%w", err)
	}

	updateReq := &servertype.UpdateScooterStatusRequest{}
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

	scooterStatusEvent := scootertype.ScooterStatusUpdaterInput{
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

func (s *Server) FindScooter(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {

	query := servertype.FindScootersRequest{}

	query.BottomLeft = scootertype.Location{}
	if err := json.Unmarshal([]byte(r.URL.Query().Get("bottomLeft")), &query.BottomLeft); err != nil {
		return apperror.NewAppError(
			apperror.WithStatusCode(http.StatusBadRequest),
			apperror.WithError(fmt.Errorf("FinedScooter: invalid bottomLeft coordination, err= %w , expected : [latitude, longitude] ", err)),
		)
	}

	query.TopRight = scootertype.Location{}
	if err := json.Unmarshal([]byte(r.URL.Query().Get("topRight")), &query.TopRight); err != nil {
		return apperror.NewAppError(
			apperror.WithStatusCode(http.StatusBadRequest),
			apperror.WithError(fmt.Errorf("FinedScooter: invalid topRight coordination, err= %w , expected : [latitude, longitude] ", err)),
		)
	}

	if err := validator.New().Struct(query); err != nil {
		return apperror.NewAppError(
			apperror.WithError(err),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	result, err := s.ScooterFinder.RectangularQuery(r.Context(), query.BottomLeft, query.TopRight)
	if err != nil {
		return err
	}

	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("FindScooter: failed to create json response, err= %w ", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(body)

	return nil
}
