package handlers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/bookmanjunior/members-only/config"
)

func serverError(a *config.Application, err error, w *http.ResponseWriter) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.ErrorLog.Printf(trace)
	http.Error(*w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func clientError(w *http.ResponseWriter, a *config.Application, err error, status int) {
	a.ErrorLog.Println(err)
	http.Error(*w, http.StatusText(status), http.StatusBadRequest)
}

func notFound(w *http.ResponseWriter, a *config.Application, err error) {
	clientError(w, a, err, 404)
}
