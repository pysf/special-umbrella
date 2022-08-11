package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/apperror"
	"github.com/pysf/special-umbrella/internal/scooter"
	scooterStatus "github.com/pysf/special-umbrella/internal/scooter/status"
)

type Server struct {
	StatusUpdater scooter.StatusUpdater
}

func NewServer() (*Server, error) {

	statusUpdater, err := scooterStatus.NewStatusUpdater()
	if err != nil {
		return nil, err
	}

	return &Server{
		StatusUpdater: statusUpdater,
	}, nil

}

func (s *Server) Start() {

	router := httprouter.New()
	router.POST("/api/scooter/status", s.wrapWithErrorHandler(s.UpdateScooterStatus))

	fmt.Println("Starting...")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func (Server) wrapWithErrorHandler(fn func(http.ResponseWriter, *http.Request, httprouter.Params) error) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := fn(w, r, ps)
		if err == nil {
			return
		}

		appErr, ok := err.(apperror.AppError)
		if !ok {
			fmt.Printf("An error occured err= %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		b, err := appErr.ResponseBody()
		if err != nil {
			w.WriteHeader(500)
			return
		}

		status, headers := appErr.ResponseHeaders()
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		w.Write(b)
	}

}
