package handlers

import (
	"net/http"

	"github.com/bookmanjunior/members-only/config"
)

func HandleGetAvatars(a *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		avatars, err := a.Avatar.GetAll()

		if err != nil {
			serverError(w, a, err)
			return
		}

		WriteJSON(w, http.StatusOK, avatars)
	}
}
