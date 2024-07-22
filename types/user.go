package types

type User struct {
	Id       int
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}
