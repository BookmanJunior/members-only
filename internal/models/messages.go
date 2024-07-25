package models

import (
	"database/sql"
	"time"
)

type MessageModel struct {
	DB *sql.DB
}

type Message struct {
	Id      int       `json:"id"`
	Message string    `json:"message"`
	Time    time.Time `json:"date"`
	User
}

func (m *MessageModel) GetAll() ([]Message, error) {
	var messages []Message
	var queryString = `select "username", "avatar_url", "message_body", date from "messages"
	inner join "users" on messages.user_id = users.id inner join "avatars" on users.avatar = avatars.id`
	res, err := m.DB.Query(queryString)

	if err != nil {
		return messages, err
	}

	for res.Next() {
		message := &Message{}
		err := res.Scan(&message.Username, &message.Avatar, &message.Message, &message.Time)

		if err != nil {
			return messages, err
		}
		messages = append(messages, *message)
	}
	return messages, nil
}
