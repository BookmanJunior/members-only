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
	m.message_id,
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
	INNER JOIN users on m.user_id = users.user_id
	INNER JOIN avatars on avatar = avatars.avatar_id
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
	m.message_id,
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
	INNER JOIN users on m.user_id = users.user_id
	INNER JOIN avatars on users.avatar = avatars.avatar_id
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
	m.message_id,
	m.message_body,
	m.date,
	m.user_id,
	"username",
	"avatar_url",
	"avatar_color"
	FROM "messages" m
	INNER JOIN users on m.user_id = users.user_id and m.message.id = $1
	INNER JOIN avatars on avatar = avatars.avatar_id and m.message_id = $1`
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

func (m *MessageModel) Insert(message string, user_id int) (Message, error) {
	var newMessage Message
	queryString :=
		`
	with new_message as (
	insert into messages(message_body, user_id)
	values ($1, $2)
	returning *
	)
	select
	nw.message_id,
	nw.message_body,
	nw.date,
	nw.user_id,
	username,
	admin,
	avatar_url,
	avatar_color,
	case
	when b.blocked_user_id is not null then true
	else false
	end as blocked
	from new_message as nw
	left join blocked_users as b on nw.user_id = b.blocked_user_id and b.user_id = $2
	join users on nw.user_id = users.user_id
	join avatars on users.avatar = avatars.avatar_id
	`

	res, err := m.DB.Query(queryString, message, user_id)

	if err != nil {
		return Message{}, err
	}

	for res.Next() {
		res.Scan(
			&newMessage.Id,
			&newMessage.Message,
			&newMessage.Time,
			&newMessage.User.Id,
			&newMessage.User.Username,
			&newMessage.User.Admin,
			&newMessage.User.Avatar.Url,
			&newMessage.User.Avatar.Color,
			&newMessage.Blocked,
		)
	}

	return newMessage, nil
}

func (m *MessageModel) Delete(message_id int) error {
	queryString := `delete from messages where message_id = $1`
	_, err := m.DB.Exec(queryString, message_id)

	if err != nil {
		return err
	}

	return nil
}

func (m *MessageModel) UpdateMessage(messageId int, newMessage string) (Message, error) {
	var updatedMessage Message
	queryString :=
		`
	with updated_message as (
	update messages
	set message_body = $2,
	modified_at = now()
	where message_id = $1
	returning *
	)
	select
	um.message_id,
	um.message_body,
	um.date,
	um.user_id,
	username,
	admin,
	avatar_color,
	avatar_url
	from updated_message as um
	join users on um.user_id = users.user_id
	join avatars on users.avatar = avatars.avatar_id
	`

	res, err := m.DB.Query(queryString, messageId, newMessage)

	if err != nil {
		return Message{}, err
	}

	for res.Next() {
		err := res.Scan(
			&updatedMessage.Id,
			&updatedMessage.Message,
			&updatedMessage.Time,
			&updatedMessage.User.Id,
			&updatedMessage.User.Username,
			&updatedMessage.User.Admin,
			&updatedMessage.User.Avatar.Url,
			&updatedMessage.User.Avatar.Color,
		)
		if err != nil {
			return Message{}, err
		}
	}
	return updatedMessage, nil
}

func (m *MessageModel) GetLatestMessages(filters filter.Filter, user_id int) ([]Message, error) {
	var messages []Message
	const queryString = `
	select *
	from (
    select m.message_id,
    m.message_body,
    m.date,
    m.user_id,
    username,
    admin,
    avatar_url,
    avatar_color,
    case
    when b.blocked_user_id is not null then true
    else false
    end as blocked
    from messages as m
    left join blocked_users as b on m.user_id = b.blocked_user_id
    and b.user_id = $1
    join users on m.user_id = users.user_id
    join avatars on users.avatar = avatars.avatar_id
    order by message_id desc
    limit $2 offset $3
    ) as sub
    order by message_id asc;
	`
	res, err := m.DB.Query(queryString, user_id, filters.Page_Size, filters.CurrentPage())

	if err != nil {
		return nil, err
	}

	for res.Next() {
		var message Message
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
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
