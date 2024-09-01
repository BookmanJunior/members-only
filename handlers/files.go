package handlers

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
	"github.com/bookmanjunior/members-only/internal/filter"
	"github.com/bookmanjunior/members-only/internal/pdf"
)

func HandleGetMessagesAsPdf(a *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		currentUser := r.Context().Value("current_user").(auth.UserClaim)

		if err != nil {
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

		messages, _, err := a.Messages.Get(filters, currentUser.Id)

		fileName := "messages" + strconv.Itoa(rand.Intn(10000)) + ".pdf"
		pdf.Generate(messages, fileName)

		defer os.RemoveAll(fileName)

		w.Header().Set("Content-Disposition", "attachment; filename=messages.pdf")
		w.Header().Set("Content-Type", "application/pdf")
		http.ServeFile(w, r, fileName)
	}
}
