package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
	"github.com/bookmanjunior/members-only/internal/filter"
	"github.com/bookmanjunior/members-only/internal/validator"
)

type messagePostRequest struct {
	Message string `json:"message"`
	validator.Validator
}

func HandleMessagesGet(a *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))

		if err != nil || page < 1 || page > 1000 {
			clientError(w, a, err, http.StatusBadRequest, map[string]string{"error": "Wrong page number"})
			return
		}

		filters := filter.Filter{
			Page:      page,
			Page_Size: 10,
			Keyword:   r.URL.Query().Get("keyword"),
			Username:  r.URL.Query().Get("username"),
			Order:     r.URL.Query().Get("order"),
		}

		messages, err := a.Messages.Get(filters)

		if err != nil {
			serverError(w, a, err)
			return
		}

		WriteJSON(w, http.StatusOK, messages)
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
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		WriteJSON(w, http.StatusOK, map[string]string{"success": "Posted message"})
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
			WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
			return
		}

		err = app.Messages.Delete(message_id)

		if err != nil {
			clientError(w, app, err, http.StatusBadRequest, map[string]string{"error": http.StatusText(http.StatusBadRequest)})
			return
		}

		WriteJSON(w, http.StatusOK, map[string]string{"success": "Deleted message"})
	}
}
