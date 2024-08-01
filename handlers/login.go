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
		err := json.NewDecoder(r.Body).Decode(&loginForm)

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

		token, err := auth.CreateToken(user.Id)

		if err != nil {
			serverError(w, app, err)
			return
		}

		w.Write([]byte(token))
	}
}
