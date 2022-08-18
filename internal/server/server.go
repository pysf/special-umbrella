package server

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/scooter/scooteriface"
)

type httpHandlerFunc func(http.ResponseWriter, *http.Request, httprouter.Params) error

type Server struct {
	ScooterReserver scooteriface.ScooterReserver
	StatusUpdater   scooteriface.StatusUpdater
	ScooterFinder   scooteriface.ScooterFinder
	HttpRouter      *httprouter.Router
	jwtTokenKey     string
}

func NewServer(statusUpdater scooteriface.StatusUpdater, scooterFinder scooteriface.ScooterFinder, scooterReserver scooteriface.ScooterReserver) (*Server, error) {

	server := &Server{
		StatusUpdater:   statusUpdater,
		ScooterFinder:   scooterFinder,
		ScooterReserver: scooterReserver,
		jwtTokenKey:     config.GetConfig("JWT_TOKEN_KEY"),
		HttpRouter:      httprouter.New(),
	}
	server.addRoutes()

	return server, nil
}

func (s *Server) addRoutes() {
	s.HttpRouter.POST("/api/scooter/status", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.AddScooterStatus)))
	s.HttpRouter.PUT("/api/scooter/reserve", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.ReserveScooter)))
	s.HttpRouter.PUT("/api/scooter/release", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.ReleaseScooter)))
	s.HttpRouter.GET("/api/scooter/search", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.FindScooter)))
}

func (s Server) Start() error {

	log.Println("Starting...")
	return http.ListenAndServe(":8080", s.HttpRouter)

}
