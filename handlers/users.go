package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/validator"
)

type userPostForm struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Avatar          int    `json:"avatar"`
	validator.Validator
}

func HandleUserGet(a *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.URL.Query().Get("id"))

		if err != nil || userId < 0 {
			notFound(&w, a, err)
			return
		}

		res, err := a.Users.Get(userId)

		if err != nil {
			notFound(&w, a, err)
			return
		}

		responseData, err := json.Marshal(res)

		if err != nil {
			notFound(&w, a, err)
			return
		}

		w.Write([]byte(responseData))

	}
}

func HandleUserPost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			clientError(&w, app, err, 400)
			return
		}

		avatar, err := strconv.Atoi(r.Form.Get("avatar"))

		form := userPostForm{
			Username:        r.Form.Get("username"),
			Password:        r.Form.Get("password"),
			ConfirmPassword: r.Form.Get("confirmPassword"),
			Avatar:          avatar,
		}

		form.CheckField(form.NotBlank(form.Username), "username", "Username can't be blank")
		form.CheckField(form.MinChars(form.Username, 2), "username", "Username must be at least 2 characters long")
		form.CheckField(form.MaxChars(form.Username, 20), "username", "Username can't be longer than 20 characters")
		form.CheckField(form.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")
		form.CheckField(form.AreFieldsEqual(form.Password, form.ConfirmPassword), "confirmPassword", "Passwords don't match")
		form.CheckField(err == nil, "avatar", "Please pick a valid avatar")

		// if form valid check if username and avatar already exists in db
		if form.Valid() {
			usernameExists := app.Users.Exists(form.Username)
			if usernameExists {
				errorMsg := fmt.Sprintf("%v already exists. Please pick a different username\n", usernameExists)
				form.AddFieldError("username", errorMsg)
				encoded, _ := json.Marshal(form.Validator.FieldErros)
				w.Write(encoded)
				return
			}
			avatarExists := app.Avatar.Exists(form.Avatar)

			if !avatarExists {
				form.AddFieldError("avatar", "Please pick a valid avatar")
				encoded, _ := json.Marshal(form.Validator.FieldErros)
				w.Write(encoded)
				return
			}
		}

		if !form.Valid() {
			encoded, _ := json.Marshal(form.Validator.FieldErros)
			w.Write(encoded)
			return
		}

		res, err := app.Users.Insert(form.Username, form.Password, form.Avatar)

		if err != nil {
			clientError(&w, app, err, 400)
			return
		}

		userM, err := json.Marshal(res)

		if err != nil {
			serverError(app, err, &w)
			return
		}

		w.Write([]byte(userM))
	}
}
