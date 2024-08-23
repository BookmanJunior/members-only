package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/filter"
)

func HandleGetMessagesAsFile(a *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		error := map[string]string{"error": "Error creating file"}

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

		file, err := os.CreateTemp(".", "messages*.json")

		if err != nil {
			clientError(w, a, err, http.StatusInternalServerError, error)
			return
		}

		defer file.Close()
		defer os.Remove(file.Name())

		d, err := json.MarshalIndent(messages, "", "   ")

		if err != nil {
			clientError(w, a, err, http.StatusInternalServerError, error)
			return
		}

		_, err = file.Write(d)

		if err != nil {
			clientError(w, a, err, http.StatusInternalServerError, error)
			return
		}

		fileName := fmt.Sprintf("filename=%v", file.Name())
		w.Header().Set("Content-Disposition", fileName)
		http.ServeFile(w, r, file.Name())

	}
}
