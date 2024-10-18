package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
	"github.com/bookmanjunior/members-only/internal/utils"
	"github.com/redis/go-redis/v9"
)

type ServerInviteModel struct {
	ServerId  int       `json:"server_id"`
	TimeLimit time.Time `json:"time_limit"`
	UseLimit  int       `json:"use_limit"`
	Link      string    `json:"link"`
	Uses      int       `json:"uses"`
}

func HandleServerInvitation(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		var input struct {
			ServerId  int       `json:"server_id"`
			TimeLimit time.Time `json:"time_limit"`
			UseLimit  int       `json:"use_limit"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			badRequest(w, err.Error())
			return
		}

		ownerId, err := app.Servers.GetOwner(input.ServerId, currentUser.Id)
		if err != nil {
			if err == sql.ErrNoRows || ownerId != currentUser.Id {
				Forbidden(w)
				return
			}
			serverError(w, app, err)
			return
		}

		inviteLink := utils.GenerateInviteLink()
		newInviteRecord := ServerInviteModel{
			ServerId:  input.ServerId,
			TimeLimit: input.TimeLimit,
			UseLimit:  input.UseLimit,
			Link:      inviteLink,
		}
		d, err := json.Marshal(newInviteRecord)
		if err != nil {
			fmt.Println("Failed to marshal")
			return
		}

		err = app.Redis.Set(context.Background(), newInviteRecord.Link, d, 0).Err()
		if err != nil {
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, Envelope{"invite_link": inviteLink})
	}
}

func HandleAddUserToServer(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)

		inv := ServerInviteModel{}
		val, err := app.Redis.Get(context.Background(), r.PathValue("link")).Result()
		if err != nil {
			if err == redis.Nil {
				notFound(w)
				return
			}
			serverError(w, app, err)
			return
		}

		_ = json.Unmarshal([]byte(val), &inv)

		if inv.Uses >= inv.UseLimit {
			responseError(w, http.StatusForbidden, "This link is no longer valid")
			return
		}

		s, err := app.ServerMembers.Insert(inv.ServerId, currentUser.Id)
		if err != nil {
			serverError(w, app, err)
			return
		}

		for u := range app.Hub.Clients {
			if currentUser.Id == u.User.Id {
				server, err := app.Servers.GetById(s)
				if err != nil {
					serverError(w, app, err)
					return
				}
				u.AddServer(server)
			}
		}

		inv.Uses++
		updatedInvLink, err := json.Marshal(inv)
		if err != nil {
			serverError(w, app, err)
			return
		}

		err = app.Redis.Set(context.Background(), inv.Link, updatedInvLink, 0).Err()
		if err != nil {
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, s)
	}
}

func HandleRemoveUserFromServer(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			UserId int `json:"user_id"`
		}
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		serverId, err := parseIdParam(r.PathValue("id"))
		if err != nil {
			badRequest(w, err.Error())
			return
		}

		serverOwnerId, err := app.Servers.GetOwner(serverId, currentUser.Id)
		if err != nil {
			if err == sql.ErrNoRows || serverOwnerId != currentUser.Id || input.UserId == serverOwnerId {
				Unauthorized(w)
				return
			}
		}

		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil || input.UserId < 1 {
			fmt.Println(err)
			badRequest(w, "IDK")
			return
		}

		err = app.ServerMembers.DeleteByUserId(serverId, input.UserId)
		if err != nil {
			serverError(w, app, err)
			return
		}

		for client, _ := range app.Hub.Clients {
			if client.User.Id == input.UserId {
				client.RemoveServer(serverId)
				break
			}
		}

		WriteJSON(w, http.StatusOK, "Removed user from server")
	}
}
