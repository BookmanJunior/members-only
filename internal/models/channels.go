package models

import (
	"database/sql"
)

type ChannelModel struct {
	DB *sql.DB
}

type Channel struct {
	Id       int    `json:"channel_id"`
	Name     string `json:"channel_name"`
	ServerId int    `json:"server_id,omitempty"`
}

func (c ChannelModel) Insert(serverId int, channelName string) (Channel, error) {
	var channel Channel
	queryString :=
		`
	insert into channels(server_id, channel_name)
	values($1, $2)
	returning
	channel_id,
	channel_name,
	server_id
	`
	res, err := c.DB.Query(queryString, serverId, channelName)
	if err != nil {
		return channel, err
	}

	for res.Next() {
		res.Scan(
			&channel.Id,
			&channel.Name,
			&channel.ServerId,
		)
	}

	return channel, nil
}

func (c ChannelModel) Update(channelName string, channelId int) (Channel, error) {
	var updatedChannel Channel
	queryString :=
		`
	update channels
	set channel_name = $1
	where channel_id = $2
	returning
	channel_id,
	channel_name,
	server_id
	`

	err := c.DB.QueryRow(queryString, channelName, channelId).Scan(
		&updatedChannel.Id,
		&updatedChannel.Name,
		&updatedChannel.ServerId,
	)
	if err != nil {
		return updatedChannel, err
	}
	return updatedChannel, nil
}

func (c ChannelModel) Delete(channelId int) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	messagesQueryString :=
		`
	delete from server_messages
	where channel_id = $1;
	`

	_, err = tx.Exec(messagesQueryString, channelId)
	if err != nil {
		tx.Rollback()
		return err
	}

	channelQueryString :=
		`
	delete from channels
	where channel_id = $1
	`

	_, err = tx.Exec(channelQueryString, channelId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
