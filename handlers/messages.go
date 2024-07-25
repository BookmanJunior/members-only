package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bookmanjunior/members-only/config"
)

func HandleMessagesGet(a *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages, err := a.Messages.GetAll()

		if err != nil {
			serverError(a, err, &w)
			return
		}

		messagesEncoded, err := json.Marshal(&messages)

		if err != nil {
			serverError(a, err, &w)
			return
		}
		w.Write([]byte(messagesEncoded))
	}
}
