package apperror

import (
	"encoding/json"
	"net/http"
)

type AppError interface {
	Error() string
	ResponseBody() ([]byte, error)
	ResponseHeaders() (int, map[string]string)
}

type HttpError struct {
	Message string `json:"detail"`
	Status  int    `json:"-"`
}

type ErrorOptions func(*HttpError)

func NewAppError(options ...ErrorOptions) error {
	const (
		defaultStatus = http.StatusInternalServerError
	)

	httpError := &HttpError{
		Status: defaultStatus,
	}

	for _, o := range options {
		o(httpError)
	}

	return httpError
}

func WithError(err error) ErrorOptions {
	return func(he *HttpError) {
		he.Message = err.Error()
	}
}

func WithStatusCode(status int) ErrorOptions {
	return func(he *HttpError) {
		he.Status = status
	}
}

func (e HttpError) Error() string {
	return e.Message
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
