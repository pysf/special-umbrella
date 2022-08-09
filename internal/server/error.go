package server

import (
	"encoding/json"
	"log"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ClientError interface {
	Error() string
	ResponseBody() ([]byte, error)
	ResponseHeaders() (int, map[string]string)
}

type HttpError struct {
	Cause  error  `json:"-"`
	Detail string `json:"detail"`
	Status int    `json:"-"`
}

func (e HttpError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Cause.Error()
}

func (e HttpError) ResponseBody() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (e HttpError) ResponseHeaders() (status int, headers map[string]string) {

	status = e.Status
	headers = map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
	return status, headers
}

func NewHttpError(err error, detail string, status int) error {
	return &HttpError{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}

func wrapWithErrorHandler(fn func(http.ResponseWriter, *http.Request, httprouter.Params) error) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := fn(w, r, ps)
		if err == nil {
			return
		}

		clientErr, ok := err.(ClientError)
		if !ok {
			log.Printf("An error occured err= %s", err)
			http.Error(w, "Unexpected error!", http.StatusInternalServerError)
			return
		}

		b, err := clientErr.ResponseBody()
		if err != nil {
			w.WriteHeader(500)
			return
		}

		status, headers := clientErr.ResponseHeaders()
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		w.Write(b)
	}

}
