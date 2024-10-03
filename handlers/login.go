package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
		err := json.NewDecoder(r.Body).Decode(loginForm)
		if err != nil {
			badCredentials(w)
			return
		}

		user, err := app.Users.GetByUsername(loginForm.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				badCredentials(w)
				return
			}
			serverError(w, app, err)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
			badCredentials(w)
			return
		}

		user.Servers, err = app.Servers.GetUsersServers(user.Id)
		if err != nil {
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

		refreshToken, err := auth.CreateRefreshToken(user.Id)
		if err != nil {
			fmt.Println("refresh error:", refreshToken)
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, map[string]any{"bearer": "Bearer " + bearerToken, "refresh": refreshToken, "user": user})
	}
}
