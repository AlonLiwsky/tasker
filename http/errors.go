package http

import (
	"errors"
	"log"
	"net/http"
)

// DebugMode writes the response to the requests with the full error message, if set to false returns user-friendly messages
const DebugMode = true

type Error interface {
	StatusAndMsg() (int, string)
}

type wrappedApiError struct {
	error
	apiErr *apiError
}

func (a wrappedApiError) Is(err error) bool {
	return a.apiErr == err
}

func (a wrappedApiError) StatusAndMsg() (int, string) {
	return a.apiErr.StatusAndMsg()
}

func WrapError(err error, apiErr apiError) error {
	return wrappedApiError{error: err, apiErr: &apiErr}
}

var (
	ErrNotFound   = apiError{msg: "not found", status: http.StatusNotFound}
	ErrBadRequest = apiError{msg: "bad request", status: http.StatusBadRequest}
)

type apiError struct {
	msg    string
	status int
}

func (a apiError) Error() string {
	return a.msg
}

func (a apiError) WithMessage(msg string) apiError {
	a.msg = msg
	return a
}

func (a apiError) StatusAndMsg() (int, string) {
	return a.status, a.msg
}

// JSONHandleError logs the error message and replies to the request with the right status code and a user-friendly msg
func JSONHandleError(w http.ResponseWriter, err error) {
	log.Printf("ERROR: %s\n", err.Error())

	defaultMsg := "internal error"
	if DebugMode {
		defaultMsg = err.Error()
	}

	var httpError Error
	if errors.As(err, &httpError) {
		status, msg := httpError.StatusAndMsg()
		if msg == "" || DebugMode {
			msg = defaultMsg
		}
		http.Error(w, msg, status)
	} else {
		http.Error(w, defaultMsg, http.StatusInternalServerError)
	}
}

// Im not proud of this
func IsNotFoundErr(err error) bool {
	var httpError Error
	if !errors.As(err, &httpError) {
		return false
	}
	status, _ := httpError.StatusAndMsg()
	return status == http.StatusNotFound
}
