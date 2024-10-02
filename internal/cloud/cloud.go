package cloud

type Cloud interface {
	Open(conString string) error
	UploadFile(userId int, fileName string) (string, error)
}
