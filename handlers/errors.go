package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/bookmanjunior/members-only/config"
)

type CustomError map[string]any

func parseMultipartFormErrors(err error) (int, error) {
	var maxBytesError *http.MaxBytesError
	switch {
	case errors.Is(err, http.ErrNotMultipart):
		return http.StatusBadRequest, errors.New("request must be of type multipart/form-data")
	case errors.As(err, &maxBytesError):
		formattedErr := fmt.Errorf("request size exceeds the set limit of %v mb", maxBytesError.Limit>>20)
		return http.StatusUnprocessableEntity, formattedErr
	case errors.Is(err, io.EOF):
		return http.StatusBadRequest, errors.New("badly formatted multipart/form-data")
	default:
		return http.StatusInternalServerError, errors.New("server failed to process your request")
	}
}

func responseError(w http.ResponseWriter, status int, message any) {
	customError := CustomError{"message": message}

	if err := WriteJSON(w, status, customError); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func serverError(w http.ResponseWriter, app *config.Application, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)

	responseError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func clientError(w http.ResponseWriter, status int, message any) {
	responseError(w, status, message)
}

func notFound(w http.ResponseWriter) {
	clientError(w, http.StatusNotFound, "Resource not found")

}

func badCredentials(w http.ResponseWriter) {
	clientError(w, http.StatusForbidden, "Wrong username or password")
}

func badRequest(w http.ResponseWriter, message any) {
	clientError(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter) {
	clientError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}

func Forbidden(w http.ResponseWriter) {
	clientError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
}
