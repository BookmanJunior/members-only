package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
	"github.com/bookmanjunior/members-only/internal/validator"
)

type messagePostRequest struct {
	Message string `json:"message_body"`
	validator.Validator
}

func HandleMessagesGet(a *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages, err := a.Messages.GetAll()

		if err != nil {
			serverError(w, a, err)
			return
		}

		messagesEncoded, err := json.Marshal(&messages)

		if err != nil {
			serverError(w, a, err)
			return
		}
		w.Write([]byte(messagesEncoded))
	}
}

func HandleMessagePost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := &messagePostRequest{}
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		err := json.NewDecoder(r.Body).Decode(message)

		if err != nil {
			clientError(w, app, err, http.StatusBadRequest, map[string]string{"error": http.StatusText(http.StatusBadRequest)})
			return
		}

		message.CheckField(message.NotBlank(message.Message), "message", "Message can't be empty")

		if !message.Valid() {
			WriteJSON(w, http.StatusBadRequest, message.FieldErrors)
			return
		}

		err = app.Messages.Insert(message.Message, int(currentUser.Id))

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		WriteJSON(w, 200, map[string]string{"success": "Posted message"})
	}
}

func HandleMessageDelete(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		message_id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Message not found"})
			return
		}

		if !currentUser.Admin {
			WriteJSON(w, 401, map[string]string{"error": http.StatusText(401)})
			return
		}

		err = app.Messages.Delete(message_id)

		if err != nil {
			clientError(w, app, err, 400, map[string]string{"error": http.StatusText(400)})
			return
		}

		WriteJSON(w, 200, map[string]string{"success": "Deleted message"})
	}
}
