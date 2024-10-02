package utils

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bookmanjunior/members-only/config"
)

var allowedFilesTypes = [3]string{"image/png", "image/jpg", "image/jpeg"}

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

func RemoveCopiedFile(fileName string) error {
	return os.Remove(filepath.Join("./attachments", filepath.Base(fileName)))
}

func CheckFileType(file multipart.File) bool {
	peek := make([]byte, 512)
	file.Read(peek)
	defer file.Seek(0, io.SeekStart) // reset read to line 0 after initial read
	fileType := http.DetectContentType(peek)
	return isAllowedFileType(fileType)
}

func isAllowedFileType(fileType string) bool {
	for _, f := range allowedFilesTypes {
		if fileType == f {
			return true
		}
	}
	return false
}
