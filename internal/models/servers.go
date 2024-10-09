package models

import (
	"database/sql"
)

type ServerModel struct {
	DB *sql.DB
}

type Server struct {
	Id       int       `json:"server_id"`
	Name     string    `json:"server_name"`
	Icon     string    `json:"server_icon"`
	Type     string    `json:"server_type"`
	Channels []Channel `json:"server_channels,omitempty"`
	Members  []User    `json:"server_members,omitempty"`
}

func (s ServerModel) GetUsersServers(userId int) ([]Server, error) {
	servers := []Server{}
	queryString :=
		`
	select distinct on(server_id) s.server_id,
	s.server_name,
	s.server_img
	from servers s
	where exists (
	select * from server_members sm where user_id = $1
	)
	`
	res, err := s.DB.Query(queryString, userId)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var server Server
		err = res.Scan(
			&server.Id,
			&server.Name,
			&server.Icon,
		)
		if err != nil {
			return nil, err
		}

		servers = append(servers, server)
	}
	return servers, nil
}

func (s ServerModel) GetById(serverId int) (Server, error) {
	var server Server
	channelQueryString :=
		`
	select
	s.server_id,
	s.server_name,
	s.server_img,
	ch.channel_id,
	ch.channel_name
	from servers s
	join channels ch on ch.server_id = $1 and s.server_id = $1
	`
	membersQueryString :=
		`
	select 
	u.user_id,
	u.username
	from server_members sm
	join users u on u.user_id = sm.user_id and sm.server_id = $1
	`
	res, err := s.DB.Query(channelQueryString, serverId)
	if err != nil {
		return server, err
	}

	for res.Next() {
		var channel Channel
		err := res.Scan(
			&server.Id,
			&server.Name,
			&server.Icon,
			&channel.Id,
			&channel.Name,
		)
		if err != nil {
			return server, err
		}
		server.Channels = append(server.Channels, channel)
	}

	usersRes, err := s.DB.Query(membersQueryString, serverId)
	if err != nil {
		return server, err
	}

	for usersRes.Next() {
		var user User
		err = usersRes.Scan(
			&user.Id,
			&user.Username,
		)
		if err != nil {
			return server, err
		}
		server.Members = append(server.Members, user)
	}

	return server, nil
}

func (s ServerModel) CreateServerTx(name, icon string, ownerId int) (Server, error) {
	var server Server
	// null represents output returned by the procedure
	err := s.DB.QueryRow("call createNewServer($1, $2, $3, null)", name, icon, ownerId).Scan(&server.Id)
	if err != nil {
		return server, err
	}

	return server, nil

}

func (s ServerModel) Insert(name, icon string, ownerId int) (Server, error) {
	var server Server
	queryString :=
		`
	with new_server as (
	insert into servers (server_name, server_img, owner_id)
	values ($1, $2, $3)
	returning *
	)
	select
	server_id,
	server_name,
	server_img
	from new_server
	`
	res, err := s.DB.Query(queryString, name, icon, ownerId)
	if err != nil {
		return server, err
	}

	for res.Next() {
		res.Scan(
			&server.Id,
			&server.Name,
			&server.Icon,
		)
	}

	return server, nil
}

func (s ServerModel) Update(ser Server) (Server, error) {
	var server Server
	queryString :=
		`
	with updated_server as (
	update server
	set server_name = $1,
	server_icon = $2
	where server_id = $3
	returning *
	)
	select
	server_id,
	server_name,
	server_img
	`
	res, err := s.DB.Query(queryString, ser.Name, ser.Icon, ser.Id)
	if err != nil {
		return server, err
	}

	for res.Next() {
		res.Scan(
			&server.Id,
			&server.Name,
			&server.Icon,
		)
	}

	return server, nil
}

func (s ServerModel) Delete(serverId int) (Server, error) {
	var server Server
	queryString :=
		`
	delete
	from servers
	where server_id = $1
	returning *
	`

	res, err := s.DB.Query(queryString, serverId)
	if err != nil {
		return server, err
	}

	for res.Next() {
		res.Scan(
			&server.Id,
			&server.Name,
			&server.Icon,
		)
	}

	return server, nil
}
