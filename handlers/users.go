package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/validator"
)

type userPostForm struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	AvatarId        int    `json:"avatar_id"`
	validator.Validator
}

func HandleUserGet(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := parseIdParam(r.URL.Query().Get("id"))
		if err != nil {
			badRequest(w, err.Error())
			return
		}

		res, err := app.Users.GetById(userId)
		if err != nil {
			if err == sql.ErrNoRows {
				notFound(w)
				return
			}
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, res)

	}
}

func HandleUserPost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form userPostForm
		err := json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			badRequest(w, "Bad request")
			return
		}

		form.CheckField(form.NotBlank(form.Username), "username", "Username can't be blank")
		form.CheckField(form.MinChars(form.Username, 2), "username", "Username must be at least 2 characters long")
		form.CheckField(form.MaxChars(form.Username, 20), "username", "Username can't be longer than 20 characters")
		form.CheckField(form.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")
		form.CheckField(form.AreFieldsEqual(form.Password, form.ConfirmPassword), "confirmPassword", "Passwords don't match")

		// if form valid check if username and avatar already exists in db
		if form.Valid() {
			usernameExists := app.Users.Exists(form.Username)
			if usernameExists {
				errorMsg := fmt.Sprintf("%v already exists. Please pick a different username", form.Username)
				form.AddFieldError("username", errorMsg)
			}

			avatarExists := app.Avatar.Exists(form.AvatarId)
			if !avatarExists {
				form.AddFieldError("avatar", "Please pick a valid avatar")
			}
		}

		if !form.Valid() {
			badRequest(w, form.Validator.FieldErrors)
			return
		}

		res, err := app.Users.Insert(form.Username, form.Password, form.AvatarId)
		if err != nil {
			serverError(w, app, err)
			return
		}

		userM, err := json.Marshal(res)
		if err != nil {
			serverError(w, app, err)
			return
		}

		w.Write([]byte(userM))
	}
}
