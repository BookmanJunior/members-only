package handlers

import (
	"fmt"
	"net/http"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
	"github.com/bookmanjunior/members-only/internal/filter"
)

func HandleGetChannel(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := filter.Filter{
			Page_Size: 50,
		}
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		page, err := parseIdParam(r.URL.Query().Get("page"))
		if err != nil {
			filters.Page = 1
		} else {
			filters.Page = page
		}
		serverId, err := parseIdParam(r.PathValue("id"))
		if err != nil {
			badRequest(w, err.Error())
			return
		}

		channelId, err := parseIdParam(r.PathValue("channelId"))
		if err != nil {
			badRequest(w, err.Error())
			return
		}

		if isAllowed, err := app.ServerMembers.IsAllowed(serverId, currentUser.Id); !isAllowed || err != nil {
			Unauthorized(w)
			return
		}

		messages, err := app.ServerMessages.GetMessagesByChannelIdAndUserId(channelId, currentUser.Id, filters)
		if err != nil {
			fmt.Println("Error happened here")
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, messages)
	}
}
