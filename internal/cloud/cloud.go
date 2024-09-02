package cloud

type Cloud interface {
	Open(conString string) error
	UploadFile(userId int, file, fileName string) (string, error)
}
