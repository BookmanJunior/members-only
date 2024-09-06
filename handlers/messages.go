package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/auth"
	"github.com/bookmanjunior/members-only/internal/filter"
	"github.com/bookmanjunior/members-only/internal/utils"
	"github.com/bookmanjunior/members-only/internal/validator"
)

type messagePostRequest struct {
	Message string `json:"message"`
	validator.Validator
}

func HandleMessagesGet(a *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		currentUser := r.Context().Value("current_user").(auth.UserClaim)

		if err != nil || page < 1 || page > 1000 {
			clientError(w, a, err, http.StatusBadRequest, CustomError{"message": "Wrong page number"})
			return
		}

		filters := filter.Filter{
			Page:      page,
			Page_Size: 10,
			Keyword:   r.URL.Query().Get("keyword"),
			Username:  r.URL.Query().Get("username"),
			Order:     r.URL.Query().Get("order"),
		}

		messages, metadata, err := a.Messages.Get(filters, currentUser.Id)

		if len(messages) <= 0 {
			WriteJSON(w, http.StatusNotFound, CustomError{"message": "No records"})
			return
		}

		if err != nil {
			serverError(w, a, err)
			return
		}

		WriteJSON(w, http.StatusOK, map[string]any{"metadata": metadata, "messages": messages})
	}
}

func HandleMessagePost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := &messagePostRequest{}
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		r.Body = http.MaxBytesReader(w, r.Body, int64(currentUser.FileSizeLimit)<<20)
		err := r.ParseMultipartForm(int64(currentUser.FileSizeLimit) << 20)

		if err != nil {
			errMsg := fmt.Sprintf("File size can't exceed %v mb", currentUser.FileSizeLimit)
			clientError(w, app, err, http.StatusRequestEntityTooLarge, map[string]string{"error": errMsg})
			return
		}

		message.Message = r.Form.Get("message")
		message.CheckField(message.NotBlank(message.Message), "message", "Message can't be empty")
		if !message.Valid() {
			WriteJSON(w, http.StatusBadRequest, message.FieldErrors)
			return
		}

		files := r.MultipartForm.File
		var messageId int
		var uploadedFilesUrl []string

		if len(files) > 0 {
			for _, file := range files {
				for _, fileHeader := range file {
					data, err := fileHeader.Open()

					if err != nil {
						serverError(w, app, err)
						return
					}

					defer data.Close()
					defer os.Remove(filepath.Join("./attachments", fileHeader.Filename))
					isCorrectFileType := utils.CheckFileType(data)
					if isCorrectFileType {
						err := utils.CopyFile(app, fileHeader, data)
						if err != nil {
							WriteJSON(w, http.StatusInternalServerError, CustomError{"message": "Failed to copy file"})
							return
						}
						uploadedFileUrl, err := app.Cloud.UploadFile(currentUser.Id, filepath.Join("./attachments", filepath.Base(fileHeader.Filename)), fileHeader.Filename)
						if err != nil {
							serverError(w, app, err)
							return
						}
						uploadedFilesUrl = append(uploadedFilesUrl, uploadedFileUrl)
					} else {
						WriteJSON(w, http.StatusUnprocessableEntity, CustomError{"message": "File must be of type image"})
						return
					}
				}
			}
		}

		if len(uploadedFilesUrl) <= 0 {
			messageId, err = app.Messages.Insert(message.Message, int(currentUser.Id))
		} else {
			fileUrls := strings.Join(uploadedFilesUrl, " ")
			messageId, err = app.Messages.Insert(fileUrls+" "+message.Message, currentUser.Id)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		messageRes, err := app.Messages.GetById(messageId)

		if err != nil {
			serverError(w, app, err)
			return
		}

		WriteJSON(w, http.StatusOK, messageRes)
	}
}

func HandleMessageDelete(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("current_user").(auth.UserClaim)
		message_id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			WriteJSON(w, http.StatusNotFound, CustomError{"message": "Message not found"})
			return
		}

		if !currentUser.Admin {
			WriteJSON(w, http.StatusUnauthorized, CustomError{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		err = app.Messages.Delete(message_id)

		if err != nil {
			clientError(w, app, err, http.StatusBadRequest, CustomError{"message": http.StatusText(http.StatusBadRequest)})
			return
		}

		WriteJSON(w, http.StatusOK, CustomError{"message": "Deleted message"})
	}
}
