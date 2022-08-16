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
		return apperror.NewAppError(
			apperror.WithError(fmt.Errorf("AddScooterStatus: failed to read request body err=%w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	var updateReq servertype.UpdateScooterStatusRequest
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

	_, err = s.StatusUpdater.UpdateStatus(r.Context(), struct {
		ScooterID string
		Timestamp string
		Latitude  string
		Longitude string
		EventType string
	}{
		EventType: updateReq.EventType,
		ScooterID: updateReq.ScooterID,
		Timestamp: updateReq.Timestamp,
		Latitude:  updateReq.Latitude,
		Longitude: updateReq.Longitude,
	})
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	return nil
}

func (s *Server) FindScooter(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {

	var query servertype.FindScootersRequest
	query.Status = r.URL.Query().Get("status")

	if err := json.Unmarshal([]byte(r.URL.Query().Get("bottomLeft")), &query.BottomLeft); err != nil {
		return apperror.NewAppError(
			apperror.WithStatusCode(http.StatusBadRequest),
			apperror.WithError(fmt.Errorf("FinedScooter: invalid bottomLeft coordination, err= %w , expected : [latitude, longitude] ", err)),
		)
	}

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

	result, err := s.ScooterFinder.RectangularQuery(r.Context(), struct {
		Status     string
		BottomLeft scootertype.Location
		TopRight   scootertype.Location
	}{
		Status:     query.Status,
		BottomLeft: query.BottomLeft,
		TopRight:   query.TopRight,
	})
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

func (s *Server) ReserveScooter(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return apperror.NewAppError(
			apperror.WithError(fmt.Errorf("ReserverScooter: failed to read request body err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	var reserveRequest servertype.ReserveScooterRequest
	if err := json.Unmarshal(b, &reserveRequest); err != nil {
		return apperror.NewAppError(
			apperror.WithError(fmt.Errorf("ReserverScooter: failed to parse json body err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	if err := validator.New().Struct(reserveRequest); err != nil {
		return apperror.NewAppError(
			apperror.WithError(err),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	isReserved, err := s.ScooterReserver.ReserveScooter(r.Context(), reserveRequest.ID)
	if err != nil {
		return fmt.Errorf("ReserveScooter: failed to reserve the scooter err= %w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if isReserved {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}

	return nil
}

func (s *Server) ReleaseScooter(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return apperror.NewAppError(
			apperror.WithError(fmt.Errorf("ReleaseScooter: failed to read request body err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	var releaseRequest servertype.ReleaseScooterRequest
	if err := json.Unmarshal(b, &releaseRequest); err != nil {
		return apperror.NewAppError(
			apperror.WithError(fmt.Errorf("ReleaseScooter: failed to parse json body err= %w", err)),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	if err := validator.New().Struct(releaseRequest); err != nil {
		return apperror.NewAppError(
			apperror.WithError(err),
			apperror.WithStatusCode(http.StatusBadRequest),
		)
	}

	if err = s.ScooterReserver.ReleaseScooter(r.Context(), releaseRequest.ID); err != nil {
		return fmt.Errorf("ReleaseScooter: failed to release the scooter err= %w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	return nil
}
