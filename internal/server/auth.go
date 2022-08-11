package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/apperror"
)

func (s Server) wrapWithAuthenticator(fn httpHandlerFunc) httpHandlerFunc {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {

		token := r.Header.Get("Authorization")
		//todo use regex to replace bearer case-insensitively
		token = strings.Replace(token, "Bearer ", "", 1)

		if token == "" {
			return apperror.NewAppError(
				apperror.WithError(fmt.Errorf(http.StatusText(http.StatusUnauthorized))),
				apperror.WithStatusCode(http.StatusUnauthorized),
			)
		}

		_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(s.jwtTokenKey), nil
		}, jwt.WithValidMethods([]string{"HS256"}))

		if err != nil {
			return apperror.NewAppError(
				apperror.WithError(fmt.Errorf(http.StatusText(http.StatusUnauthorized))),
				apperror.WithStatusCode(http.StatusUnauthorized),
			)
		}

		if err = fn(w, r, p); err != nil {
			return err
		}

		return nil
	}
}
