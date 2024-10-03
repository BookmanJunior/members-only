package handlers

import (
	"net/http"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
	"github.com/bookmanjunior/members-only/internal/hub"
	"github.com/gorilla/websocket"
)

func HandleWs(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		user, err := app.Users.GetById(currentUser.Id)
		if err != nil {
			serverError(w, app, err)
			return
		}

		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			serverError(w, app, err)
			return
		}

		client := hub.CreateNewClient(user, conn, app.Hub, app.Messages)

		client.Hub.RegisterCh <- client

		go client.Read()
		go client.Write()
	}
}
