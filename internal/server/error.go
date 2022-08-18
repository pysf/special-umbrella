package server

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pysf/special-umbrella/internal/apperror"
)

func (Server) wrapWithErrorHandler(fn httpHandlerFunc) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := fn(w, r, ps)
		if err == nil {
			return
		}

		appErr, ok := err.(apperror.AppError)
		if !ok {
			log.Printf("wrapWithErrorHandler: err= %s", err)
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
