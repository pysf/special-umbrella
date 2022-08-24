package server

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/config"
	"github.com/pysf/special-umbrella/internal/scooter"
	"github.com/pysf/special-umbrella/internal/scooter/scooteriface"
	"go.mongodb.org/mongo-driver/mongo"
)

type httpHandlerFunc func(http.ResponseWriter, *http.Request, httprouter.Params) error

type Server struct {
	ScooterReserver scooteriface.ScooterReserver
	StatusUpdater   scooteriface.StatusUpdater
	ScooterFinder   scooteriface.ScooterFinder
	HttpHandler     http.Handler
	jwtTokenKey     string
}

func NewServer(DB *mongo.Database) (*Server, error) {

	scooterReserver, err := scooter.NewScooterReserver(DB)
	if err != nil {
		panic(err)
	}

	statusUpdater, err := scooter.NewStatusUpdater(scooterReserver, DB)
	if err != nil {
		panic(err)
	}

	scooterFinder, err := scooter.NewScooterFinder(DB)
	if err != nil {
		panic(err)
	}

	server := &Server{
		StatusUpdater:   statusUpdater,
		ScooterFinder:   scooterFinder,
		ScooterReserver: scooterReserver,
		jwtTokenKey:     config.AppConf().JWTTokenKey,
	}
	server.addRoutes()

	return server, nil
}

func (s *Server) addRoutes() {
	httpRouter := httprouter.New()
	httpRouter.POST("/api/scooter/status", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.AddScooterStatus)))
	httpRouter.PUT("/api/scooter/reserve", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.ReserveScooter)))
	httpRouter.PUT("/api/scooter/release", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.ReleaseScooter)))
	httpRouter.GET("/api/scooter/search", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.FindScooter)))
	s.HttpHandler = httpRouter
}

func (s Server) Start() error {

	log.Println("Starting...")
	return http.ListenAndServe(":8080", s.HttpHandler)

}
