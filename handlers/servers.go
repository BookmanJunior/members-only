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
	}
}

func HandlePostServer(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		r.Body = http.MaxBytesReader(w, r.Body, 3<<20)
		err := r.ParseMultipartForm(3 << 20)
		if err != nil {
			clientError(w, http.StatusUnprocessableEntity, CustomError{"message": "Server image can't be greater than 3 mb"})
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
				WriteJSON(w, http.StatusBadRequest, CustomError{"message": "Icon must be of type jpg or png."})
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

		createServerRes, err := app.Server.CreateServerTx(serverName, serverIconUrl, currentUser.Id)
		if err != nil {
			serverError(w, app, err)
			return
		}

		app.InfoLog.Printf("Created server %v by user %v", createServerRes.Id, currentUser.Id)
		WriteJSON(w, http.StatusOK, CustomError{"message": "Created Server"})
	}
}
