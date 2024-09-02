package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/bookmanjunior/members-only/config"
)

func CopyFile(app *config.Application, fileHeader *multipart.FileHeader, file multipart.File) error {
	copyDst, err := os.Create(filepath.Join("./attachments", filepath.Base(fileHeader.Filename)))

	if err != nil {
		app.ErrorLog.Println(err)
		return err
	}

	defer copyDst.Close()

	_, err = io.Copy(copyDst, file)

	if err != nil {
		app.ErrorLog.Println(err)
		return err
	}

	app.InfoLog.Printf("Successfully copied file: %v\n", fileHeader.Filename)
	return nil
}
