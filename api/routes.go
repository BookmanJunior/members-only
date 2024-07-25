package api

import (
	"net/http"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/handlers"
)

func Router(app *config.Application) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Caught"))
	})
	router.HandleFunc("GET /users", handlers.HandleUserGet(app))
	router.HandleFunc("GET /messages", handlers.HandleMessagesGet(app))

	return router
}
