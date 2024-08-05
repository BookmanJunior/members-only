package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/bookmanjunior/members-only/config"
)

type customError map[string]string

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func serverError(w http.ResponseWriter, a *config.Application, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.ErrorLog.Printf(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func clientError(w http.ResponseWriter, a *config.Application, err error, status int, message customError) {
	a.ErrorLog.Println(err)
	WriteJSON(w, status, message)
}

func notFound(w http.ResponseWriter, a *config.Application, err error) {
	clientError(w, a, err, http.StatusNotFound, map[string]string{"error": "Not found"})

}

func badCredentials(w http.ResponseWriter, a *config.Application, err error) {
	clientError(w, a, err, http.StatusUnauthorized, map[string]string{"error": "Wrong username or password"})
}

func Unauthorized(w http.ResponseWriter, a *config.Application, err error) {
	WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
}
