package handlers

import (
	"database/sql"
	"encoding/json"
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

func HandleCreateChannel(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			ChannelName string `json:"channel_name"`
		}
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		serverId, err := parseIdParam(r.PathValue("id"))
		if err != nil {
			badRequest(w, err.Error())
			return
		}

		serverOwnerId, err := app.Servers.GetOwner(serverId, currentUser.Id)
		if err != nil {
			if err != sql.ErrNoRows || serverOwnerId != currentUser.Id {
				Unauthorized(w)
				return
			}
			serverError(w, app, err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			badRequest(w, "Channel name can't be empty")
			return
		}

		newChannel, err := app.Channels.Insert(serverId, input.ChannelName)
		if err != nil {
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, newChannel)
	}
}

func HandleUpdateChannel(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			ChannelName string `json:"channel_name"`
		}
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
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

		serverOwnerId, err := app.Servers.GetOwner(serverId, currentUser.Id)
		if err != nil {
			if err != sql.ErrNoRows || serverOwnerId != currentUser.Id {
				Unauthorized(w)
				return
			}
			serverError(w, app, err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			badRequest(w, "Channel name can't be empty")
			return
		}

		updatedChannel, err := app.Channels.Update(input.ChannelName, channelId)
		if err != nil {
			if err == sql.ErrNoRows {
				badRequest(w, "Bad channel id")
				return
			}
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, updatedChannel)
	}
}

func HandleDeleteChannel(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
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

		serverOwnerId, err := app.Servers.GetOwner(serverId, currentUser.Id)
		if err != nil {
			if err != sql.ErrNoRows || serverOwnerId != currentUser.Id {
				Unauthorized(w)
				return
			}
			serverError(w, app, err)
			return
		}

		err = app.Channels.Delete(channelId)
		if err != nil {
			serverError(w, app, err)
			return
		}

		fmt.Println("Called")

		WriteJSON(w, http.StatusOK, Envelope{"message": "Deleted channel"})
	}
}
