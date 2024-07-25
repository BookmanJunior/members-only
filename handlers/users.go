package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bookmanjunior/members-only/config"
)

func HandleUserGet(a *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.URL.Query().Get("id"))

		if err != nil || userId < 0 {
			notFound(&w, a, err)
			return
		}

		res, err := a.Users.Get(userId)

		if err != nil {
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
