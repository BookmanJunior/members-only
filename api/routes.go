package api

import (
	"net/http"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/handlers"
	"github.com/bookmanjunior/members-only/middleware"
)

func Router(app *config.Application) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /login", handlers.HandleLogin(app))
	router.HandleFunc("GET /users", middleware.IsAuthorized(app, handlers.HandleUserGet(app)))
	router.HandleFunc("POST /users", handlers.HandleUserPost(app))
	router.HandleFunc("GET /messages", middleware.IsAuthorized(app, handlers.HandleMessagesGet(app)))
	router.HandleFunc("POST /messages", middleware.IsAuthorized(app, handlers.HandleMessagePost(app)))
	router.HandleFunc("DELETE /messages/{id}", middleware.IsAuthorized(app, handlers.HandleMessageDelete(app)))
	router.HandleFunc("GET /avatars", handlers.HandleGetAvatars(app))
	router.HandleFunc("GET /files/messages", middleware.IsAuthorized(app, handlers.HandleGetMessagesAsPdf(app)))
	router.HandleFunc("GET /ws", middleware.IsAuthorized(app, handlers.HandleWs(app)))

	return middleware.RecoverPanic(app, middleware.EnableCors(app, router))
}
