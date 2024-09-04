package models

import (
	"database/sql"
	"time"

	"github.com/bookmanjunior/members-only/internal/filter"
)

type MessageModel struct {
	DB *sql.DB
}

type Message struct {
	Id      int       `json:"message_id"`
	Message string    `json:"message"`
	Blocked bool      `json:"blocked"`
	Time    time.Time `json:"date"`
	User    User      `json:"user"`
}

func (m *MessageModel) GetAll(userId int) ([]Message, error) {
	var messages []Message
	var queryString = `
	select
	m.id,
	m.message_body,
	m.date,
	m.user_id,
	"username",
	"admin",
	"avatar_url",
	"avatar_color",
	CASE
	    WHEN b.blocked_user_id IS NOT NULL THEN TRUE
	    ELSE FALSE
	END AS blocked
	FROM "messages" m
	LEFT JOIN blocked_users b
	ON m.user_id = b.blocked_user_id
	AND b.user_id = $1
	INNER JOIN users on m.user_id = users.id
	INNER JOIN avatars on avatar = avatars.id
	`
	res, err := m.DB.Query(queryString, userId)

	if err != nil {
		return messages, err
	}

	for res.Next() {
		message := &Message{}
		err := res.Scan(
			&message.Id,
			&message.Message,
			&message.Time,
			&message.User.Id,
			&message.User.Username,
			&message.User.Admin,
			&message.User.Avatar.Url,
			&message.User.Avatar.Color,
			&message.Blocked,
		)

		if err != nil {
			return messages, err
		}
		messages = append(messages, *message)
	}
	return messages, nil
}

func (m *MessageModel) Get(filters filter.Filter, userId int) ([]Message, filter.MetaData, error) {
	var messages []Message
	var totalRecrods int
	var queryString = `
	select
	count(*) OVER(),
	m.id,
	m.user_id,
	m.message_body,
	m.date,
	"username",
	"admin",
	"avatar_url",
	"avatar_color",
	CASE
	    WHEN b.blocked_user_id IS NOT NULL THEN TRUE
	    ELSE FALSE
	END AS blocked
	FROM "messages" m
	LEFT JOIN blocked_users b
	ON m.user_id = b.blocked_user_id
	AND b.user_id = $5
	INNER JOIN users on m.user_id = users.id
	INNER JOIN avatars on users.avatar = avatars.id
	WHERE m.message_body ILIKE $1 and "username" ILIKE $2
	limit $3 offset $4;`
	res, err := m.DB.Query(queryString,
		"%"+filters.Keyword+"%", "%"+filters.Username+"%", filters.Page_Size, filters.CurrentPage(), userId)

	if err != nil {
		return nil, filter.MetaData{}, err
	}

	for res.Next() {
		var message Message
		res.Scan(
			&totalRecrods,
			&message.Id,
			&message.User.Id,
			&message.Message,
			&message.Time,
			&message.User.Username,
			&message.User.Admin,
			&message.User.Avatar.Url,
			&message.User.Avatar.Color,
			&message.Blocked,
		)
		messages = append(messages, message)
	}

	metadata := filter.CalculateMetaData(totalRecrods, filters.Page, filters.Page_Size)

	return messages, metadata, nil
}

func (m *MessageModel) GetById(messageId int) (Message, error) {
	var queryString = `
	select
	m.id,
	m.message_body,
	m.date,
	m.user_id,
	"username",
	"avatar_url",
	"avatar_color"
	FROM "messages" m
	INNER JOIN users on m.user_id = users.id and m.id = $1
	INNER JOIN avatars on avatar = avatars.id and m.id = $1`
	var message Message
	res, err := m.DB.Query(queryString, messageId)

	if err != nil {
		return message, err
	}

	for res.Next() {
		res.Scan(
			&message.Id,
			&message.Message,
			&message.Time,
			&message.User.Id,
			&message.User.Username,
			&message.User.Avatar.Url,
			&message.User.Avatar.Color,
		)
	}

	return message, nil
}

func (m *MessageModel) Insert(message string, user_id int) (int, error) {
	queryString := `insert into messages (message_body, user_id) values ($1, $2) returning id`
	var messageId int
	res, err := m.DB.Query(queryString, message, user_id)

	if err != nil {
		return messageId, err
	}

	for res.Next() {
		res.Scan(&messageId)
	}

	return messageId, nil
}

func (m *MessageModel) Delete(message_id int) error {
	queryString := `delete from messages where id = $1`
	_, err := m.DB.Exec(queryString, message_id)

	if err != nil {
		return err
	}

	return nil
}
