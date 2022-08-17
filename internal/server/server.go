package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/scooter/scooteriface"
)

type Server struct {
	ScooterReserver scooteriface.ScooterReserver
	StatusUpdater   scooteriface.ScooterStatusUpdater
	ScooterFinder   scooteriface.ScooterFinder
	jwtTokenKey     string
}

func NewServer() (*Server, error) {

	scooterReserver, err := scooter.NewScooterReserver()
	if err != nil {
		return nil, err
	}

	statusUpdater, err := scooter.NewStatusUpdater(scooterReserver)
	if err != nil {
		return nil, err
	}

	scooterFinder, err := scooter.NewScooterFinder()
	if err != nil {
		return nil, err
	}

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
	router.POST("/api/scooter/status", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.UpdateScooterStatus)))
	router.PUT("/api/scooter/reserve", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.ReserveScooter)))
	router.PUT("/api/scooter/release", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.ReleaseScooter)))
	router.GET("/api/scooter/search", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.FindScooter)))

	fmt.Println("Starting...")
	return http.ListenAndServe(":8080", router)

}
