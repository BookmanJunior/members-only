package models

import (
	"database/sql"
	"time"
)

type ChannelModel struct {
	DB *sql.DB
}

type Channel struct {
	Id        int       `json:"channel_id"`
	Name      string    `json:"channel_name"`
	ServerId  int       `json:"server_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (c ChannelModel) Insert(serverId int, channelName string) (Channel, error) {
	var channel Channel
	queryString :=
		`
	insert into channels(server_id, channel_name)
	values($1, $2)
	returning *
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
			&channel.CreatedAt,
		)
	}

	return channel, nil
}
