package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/bookmanjunior/members-only/config"
)

type CustomError map[string]any

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
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

func Unauthorized(w http.ResponseWriter) {
	clientError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}
