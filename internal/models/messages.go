package models

import (
	"database/sql"
	"time"
)

type MessageModel struct {
	DB *sql.DB
}

type Message struct {
	Id      int       `json:"message_id"`
	Message string    `json:"message"`
	Time    time.Time `json:"date"`
	User    User      `json:"user"`
}

func (m *MessageModel) GetAll() ([]Message, error) {
	var messages []Message
	var queryString = `select messages.id , users.id, "username", "avatar_color", "avatar_url", "message_body", date from "messages"
	inner join "users" on messages.user_id = users.id inner join "avatars" on users.avatar = avatars.id`
	res, err := m.DB.Query(queryString)

	if err != nil {
		return messages, err
	}

	for res.Next() {
		message := &Message{}
		err := res.Scan(&message.Id, &message.User.Id, &message.User.Username,
			&message.User.Avatar.Color, &message.User.Avatar.Url, &message.Message, &message.Time)

		if err != nil {
			return messages, err
		}
		messages = append(messages, *message)
	}
	return messages, nil
}

func (m *MessageModel) Insert(message string, user_id int) error {
	queryString := `insert into messages (message_body, user_id) values ($1, $2)`
	_, err := m.DB.Exec(queryString, message, user_id)

	if err != nil {
		return err
	}

	return nil
}

func (m *MessageModel) Delete(message_id int) error {
	queryString := `delete from messages where id = $1`
	_, err := m.DB.Exec(queryString, message_id)

	if err != nil {
		return err
	}

	return nil
}
