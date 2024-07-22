package types

type Message struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
	User
}
