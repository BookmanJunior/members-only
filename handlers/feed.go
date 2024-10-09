package handlers

import (
	"net/http"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
)

func HandleGetFeed(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)

		servers, err := app.Servers.GetUsersServers(currentUser.Id)
		if err != nil {
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, map[string]any{
			"servers": servers,
		})
	}
}
