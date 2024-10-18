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
	router.HandleFunc("POST /refresh-token", handlers.HandleRefreshAccessToken(app))

	router.HandleFunc("GET /users", middleware.IsAuthorized(app, handlers.HandleUserGet(app)))
	router.HandleFunc("POST /users", handlers.HandleUserPost(app))

	router.HandleFunc("GET /messages", middleware.IsAuthorized(app, handlers.HandleMessagesGet(app)))
	router.HandleFunc("POST /messages", middleware.IsAuthorized(app, handlers.HandleMessagePost(app)))
	router.HandleFunc("DELETE /messages/{id}", middleware.IsAuthorized(app, handlers.HandleMessageDelete(app)))

	router.HandleFunc("GET /avatars", handlers.HandleGetAvatars(app))
	router.HandleFunc("GET /files/messages", middleware.IsAuthorized(app, handlers.HandleGetMessagesAsPdf(app)))

	router.HandleFunc("GET /servers/{id}", middleware.IsAuthorized(app, handlers.HandleGetServer(app)))
	router.HandleFunc("POST /servers", middleware.IsAuthorized(app, handlers.HandlePostServer(app)))
	router.HandleFunc("DELETE /servers/{id}", middleware.IsAuthorized(app, handlers.HandleDeleteServer(app)))
	router.HandleFunc("POST /invite", middleware.IsAuthorized(app, handlers.HandleServerInvitation(app)))

	router.HandleFunc("GET /servers/{id}/{channelId}", middleware.IsAuthorized(app, handlers.HandleGetChannel(app)))
	router.HandleFunc("POST /servers/{id}/channel", middleware.IsAuthorized(app, handlers.HandleCreateChannel(app)))
	router.HandleFunc("PATCH /servers/{id}/{channelId}", middleware.IsAuthorized(app, handlers.HandleUpdateChannel(app)))
	router.HandleFunc("DELETE /servers/{id}/{channelId}", middleware.IsAuthorized(app, handlers.HandleDeleteChannel(app)))
	router.HandleFunc("PATCH /servers/members/{inviteLink}", middleware.IsAuthorized(app, handlers.HandleAddUserToServer(app)))
	router.HandleFunc("DELETE /servers/members/{id}", middleware.IsAuthorized(app, handlers.HandleRemoveUserFromServer(app)))
	router.HandleFunc("GET /feed", middleware.IsAuthorized(app, handlers.HandleGetFeed(app)))
	router.HandleFunc("GET /ws", middleware.IsAuthorized(app, handlers.HandleWs(app)))

	return middleware.RecoverPanic(app, middleware.Logger(app, middleware.EnableCors(app, router)))
}
