package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
)

func HandleRefreshAccessToken(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			RefreshToken string `json:"refresh-token"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			badRequest(w, "Invalid refresh token format")
			return
		}

		claim, err := auth.VerifyToken(input.RefreshToken)
		if err != nil {
			Forbidden(w)
			return
		}

		user_id := claim["id"].(float64)

		user, err := app.Users.GetById(int(user_id))
		if err != nil {
			if err == sql.ErrNoRows {
				Forbidden(w)
				return
			}
			serverError(w, app, err)
			return
		}

		bearerToken, err := auth.CreateToken(auth.UserClaim{
			Id:            user.Id,
			Admin:         user.Admin,
			FileSizeLimit: user.FileSizeLimit,
		})
		if err != nil {
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, map[string]any{"bearer": "Bearer " + bearerToken})
	}
}
