package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/model"
	repository "github.com/pysf/special-umbrella/internal/repository/mongodb"
)

func NewServer() (*Server, error) {

	scooterStatusRepository, err := repository.NewScooterStatusRepository()
	if err != nil {
		return nil, err
	}

	return &Server{
		ScooterStatusRepository: scooterStatusRepository,
	}, nil

}

type Server struct {
	ScooterStatusRepository model.ScooterStatusRepository
}

func (s *Server) Start() {

	router := httprouter.New()
	router.POST("/api/scooter/status", wrapWithErrorHandler(s.AddScooterStatus))

	fmt.Println("Starting...")
	log.Fatal(http.ListenAndServe(":8080", router))

}
