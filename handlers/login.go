package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type handleLoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleLogin(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginForm := &handleLoginForm{}
		userClaim := auth.UserClaim{}
		err := json.NewDecoder(r.Body).Decode(loginForm)

		if err != nil {
			badCredentials(w, app, err)
			return
		}

		user, err := app.Users.GetByUsername(loginForm.Username)

		if err == sql.ErrNoRows {
			badCredentials(w, app, err)
			return
		}

		if err != nil {
			serverError(w, app, err)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
			badCredentials(w, app, err)
			return
		}

		userClaim.Id = user.Id
		userClaim.Admin = user.Admin

		bearerToken, err := auth.CreateToken(userClaim)

		if err != nil {
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, map[string]any{"bearer": "Bearer " + bearerToken, "user": user})
	}
}
