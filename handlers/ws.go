package handlers

import (
	"fmt"
	"net/http"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
	"github.com/bookmanjunior/members-only/internal/hub"
	"github.com/gorilla/websocket"
)

func HandleWs(a *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		user, err := a.Users.GetById(currentUser.Id)

		if err != nil {
			fmt.Println(err)
			WriteJSON(w, http.StatusInternalServerError, CustomError{"message": "Failed to get user"})
			return
		}

		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(w, r, w.Header())

		if err != nil {
			fmt.Println(err)
			serverError(w, a, err)
			return
		}

		client := hub.CreateNewClient(user, conn, a.Hub, a.Messages)

		client.Hub.RegisterCh <- client

		go client.Read()
		go client.Write()

	}
}
