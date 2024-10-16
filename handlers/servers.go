package handlers

import (
	"net/http"
	"unicode/utf8"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
	"github.com/bookmanjunior/members-only/internal/utils"
)

func HandleGetServer(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		serverId, err := parseIdParam(r.PathValue("id"))
		if err != nil {
			badRequest(w, err.Error())
			return
		}

		isAllowed, err := app.ServerMembers.IsAllowed(serverId, currentUser.Id)
		if err != nil {
			serverError(w, app, err)
			return
		}

		if !isAllowed {
			Forbidden(w)
			return
		}

		server, err := app.Servers.GetById(serverId)
		if err != nil {
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, server)
	}
}

func HandlePostServer(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		r.Body = http.MaxBytesReader(w, r.Body, 3<<20)
		err := r.ParseMultipartForm(3 << 20)
		if err != nil {
			statusCode, parsedErr := parseMultipartFormErrors(err)
			responseError(w, statusCode, parsedErr.Error())
			return
		}

		serverName := r.FormValue("server_name")
		if utf8.RuneCountInString(serverName) < 1 {
			WriteJSON(w, http.StatusBadRequest, CustomError{"message": "Server name must be 1 character or longer"})
			return
		}

		var serverIconUrl string
		file, fileHeader, err := r.FormFile("server_icon")
		if err == nil {
			defer file.Close()
			defer utils.RemoveCopiedFile(fileHeader.Filename)

			isCorrectFileType := utils.CheckFileType(file)
			if !isCorrectFileType {
				badRequest(w, "Icon must be of type jpg or png.")
				return
			}

			err = utils.CopyFile(app, fileHeader, file)
			if err != nil {
				serverError(w, app, err)
				return
			}
			serverIconUrl, err = app.Cloud.UploadFile(currentUser.Id, fileHeader.Filename)
			if err != nil {
				serverError(w, app, err)
				return
			}
		}

		newServerId, err := app.Servers.CreateServerTx(serverName, serverIconUrl, currentUser.Id)
		if err != nil {
			serverError(w, app, err)
			return
		}

		serverInfo, err := app.Servers.GetById(newServerId)
		if err != nil {
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, serverInfo)
	}
}
