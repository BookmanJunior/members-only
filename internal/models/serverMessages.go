package models

import (
	"database/sql"
	"time"

	"github.com/bookmanjunior/members-only/internal/filter"
)

type ServerMessage struct {
	Id        int       `json:"message_id"`
	ServerId  int       `json:"server_id"`
	ChannelId int       `json:"channel_id"`
	Message   string    `json:"message"`
	SentDate  time.Time `json:"sent_date"`
	Blocked   bool      `json:"blocked"`
	Edited    bool      `json:"edited"`
	User      User      `json:"user"`
}

type ServerMessageModel struct {
	DB *sql.DB
}

func (sm *ServerMessageModel) GetMessagesByChannelIdAndUserId(channelId, userId int, filters filter.Filter) ([]ServerMessage, error) {
	var messages []ServerMessage
	queryString := `
	select *
	from (
    select sm.message_id,
    sm.message_body,
    sm.sent_date,
    sm.user_id,
    username,
    admin,
    avatar_url,
    avatar_color,
    server_id,
    channel_id,
    case
    when b.blocked_user_id is not null then true
    else false
    end as blocked,
    case
    when modified_at is not null then true
    else false
    end as edited
    from server_messages as sm
    left join blocked_users as b on sm.user_id = b.blocked_user_id
    and b.user_id = $1
    join users on sm.user_id = users.user_id and channel_id = $2
    join avatars on users.avatar = avatars.avatar_id
    where channel_id = $2
    order by message_id desc
    limit $3 offset $4
    ) as sub
    order by message_id asc;
	`

	res, err := sm.DB.Query(queryString, userId, channelId, filters.Page_Size, filters.CurrentPage())
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var serverMessage ServerMessage
		err := res.Scan(
			&serverMessage.Id,
			&serverMessage.Message,
			&serverMessage.SentDate,
			&serverMessage.User.Id,
			&serverMessage.User.Username,
			&serverMessage.User.Admin,
			&serverMessage.User.Avatar.Url,
			&serverMessage.User.Avatar.Color,
			&serverMessage.ServerId,
			&serverMessage.ChannelId,
			&serverMessage.Blocked,
			&serverMessage.Edited,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, serverMessage)
	}

	return messages, nil
}

func (sm *ServerMessageModel) Insert(m ServerMessage) (ServerMessage, error) {
	var newMessage ServerMessage
	queryString :=
		`
	with new_message as (
	insert into server_messages (server_id, channel_id, user_id, message_body)
	values ($1, $2, $3, $4)
	returning *
	)
	select
	message_id,
	server_id,
	channel_id,
	u.user_id,
	message_body,
	sent_date,
	u.username,
	a.avatar_color,
	a.avatar_url
	from new_message nm
	join users u on u.user_id = nm.user_id
	join avatars a on a.avatar_id = u.avatar
	`

	err := sm.DB.QueryRow(queryString, m.ServerId, m.ChannelId, m.User.Id, m.Message).
		Scan(
			&newMessage.Id,
			&newMessage.ServerId,
			&newMessage.ChannelId,
			&newMessage.User.Id,
			&newMessage.Message,
			&newMessage.SentDate,
			&newMessage.User.Username,
			&newMessage.User.Avatar.Color,
			&newMessage.User.Avatar.Url,
		)
	if err != nil {
		return newMessage, err
	}

	return newMessage, nil
}

func (sm *ServerMessageModel) Update(newMessage string, messageId, userId int) (ServerMessage, error) {
	var updatedMessage ServerMessage
	queryString :=
		`
	with updated_message as (
	update server_messages sm
	set message_body = $1,
	modified_at = now()
	where message_id = $2 and user_id = $3
	returning *
	)
	select
	server_id,
	channel_id,
	um.message_id,
	um.message_body,
	um.sent_date,
	um.user_id,
	username,
	admin,
	avatar_color,
	avatar_url
	from updated_message as um
	join users on um.user_id = users.user_id
	join avatars on users.avatar = avatars.avatar_id
	`

	err := sm.DB.QueryRow(queryString, newMessage, messageId, userId).Scan(
		&updatedMessage.ServerId,
		&updatedMessage.ChannelId,
		&updatedMessage.Id,
		&updatedMessage.Message,
		&updatedMessage.SentDate,
		&updatedMessage.User.Id,
		&updatedMessage.User.Username,
		&updatedMessage.User.Admin,
		&updatedMessage.User.Avatar.Color,
		&updatedMessage.User.Avatar.Url,
	)
	updatedMessage.Edited = true
	if err != nil {
		return updatedMessage, err
	}

	return updatedMessage, nil
}

func (sm *ServerMessageModel) Delete(messageId, userId int) (int, error) {
	var msgId int
	queryString :=
		`
	delete from server_messages sm
	where sm.message_id = $1 and user_id = $2
	returning message_id
	`
	err := sm.DB.QueryRow(queryString, messageId, userId).Scan(&msgId)
	if err != nil {
		return msgId, err
	}

	return msgId, nil
}
