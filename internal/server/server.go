package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/scooter/scooteriface"
)

type Server struct {
	ScooterReserver scooteriface.ScooterReserver
	StatusUpdater   scooteriface.StatusUpdater
	ScooterFinder   scooteriface.ScooterFinder
	jwtTokenKey     string
}

func NewServer(statusUpdater scooteriface.StatusUpdater, scooterFinder scooteriface.ScooterFinder, scooterReserver scooteriface.ScooterReserver) (*Server, error) {

	return &Server{
		StatusUpdater:   statusUpdater,
		ScooterFinder:   scooterFinder,
		ScooterReserver: scooterReserver,
		jwtTokenKey:     config.GetConfig("JWT_TOKEN_KEY"),
	}, nil

}

type httpHandlerFunc func(http.ResponseWriter, *http.Request, httprouter.Params) error

func (s Server) Start() error {

	router := httprouter.New()
	router.POST("/api/scooter/status", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.AddScooterStatus)))
	router.PUT("/api/scooter/reserve", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.ReserveScooter)))
	router.PUT("/api/scooter/release", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.ReleaseScooter)))
	router.GET("/api/scooter/search", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.FindScooter)))

	fmt.Println("Starting...")
	return http.ListenAndServe(":8080", router)

}
