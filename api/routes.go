package api

import (
	"net/http"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/handlers"
	"github.com/bookmanjunior/members-only/middleware"
)

func Router(app *config.Application) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Caught"))
	})
	router.HandleFunc("GET /users", handlers.HandleUserGet(app))
	router.HandleFunc("POST /users", handlers.HandleUserPost(app))
	router.HandleFunc("GET /messages", handlers.HandleMessagesGet(app))

	return middleware.RecoverPanic(app, middleware.Logger(app, router))
}
