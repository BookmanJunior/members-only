package utils

import (
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

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

var lan = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
}

func GenerateInviteLink() string {
	var inviteString string

	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 8; i++ {
		n := rand.Int31n(int32(len(lan)))
		inviteString += lan[n]
	}

	return inviteString
}
