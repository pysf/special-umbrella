package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/scooter"
	scooterStatus "github.com/pysf/special-umbrella/internal/scooter/status"
)

type Server struct {
	StatusUpdater scooter.StatusUpdater
	jwtTokenKey   string
}

func NewServer() (*Server, error) {

	statusUpdater, err := scooterStatus.NewStatusUpdater()
	if err != nil {
		return nil, err
	}

	jwtTokenKey, exists := os.LookupEnv("JWT_TOKEN_KEY")
	if !exists {
		return nil, fmt.Errorf("NewServer: err= JWT_TOKEN_KEY is required")
	}

	return &Server{
		StatusUpdater: statusUpdater,
		jwtTokenKey:   jwtTokenKey,
	}, nil

}

type httpHandlerFunc func(http.ResponseWriter, *http.Request, httprouter.Params) error

func (s Server) Start() {

	router := httprouter.New()
	router.POST("/api/scooter/status", s.wrapWithErrorHandler(s.wrapWithAuthenticator(s.UpdateScooterStatus)))

	fmt.Println("Starting...")
	log.Fatal(http.ListenAndServe(":8080", router))

}
