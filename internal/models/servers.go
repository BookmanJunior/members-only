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
	s.server_img,
	po.type
	from servers s
	join privacy_options po on po.id = s.type
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
			&server.Type,
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
	po.type,
	ch.channel_id,
	ch.channel_name
	from servers s
	join privacy_options po on po.id = s.type
	join channels ch on ch.server_id = $1 and s.server_id = $1
	`
	membersQueryString :=
		`
	select
	u.user_id,
	u.username,
	a.avatar_color,
	a.avatar_url
	from server_members sm
	join users u on u.user_id = sm.user_id and sm.server_id = $1
	join avatars a on a.avatar_id = u.avatar
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
			&server.Type,
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
			&user.Avatar.Color,
			&user.Avatar.Url,
		)
		if err != nil {
			return server, err
		}
		server.Members = append(server.Members, user)
	}

	return server, nil
}

func (s ServerModel) CreateServerTx(name, icon string, ownerId int) (int, error) {
	var serverId int
	// null represents output returned by the procedure
	err := s.DB.QueryRow("call createNewServer($1, $2, $3, null)", name, icon, ownerId).Scan(&serverId)
	if err != nil {
		return serverId, err
	}

	return serverId, nil

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
	tx, err := s.DB.Begin()
	if err != nil {
		return server, err
	}

	serverMembersQueryString :=
		`
	delete from server_members
	where server_id = $1
	`

	_, err = tx.Exec(serverMembersQueryString, serverId)
	if err != nil {
		tx.Rollback()
		return server, err
	}

	messagesQueryString :=
		`
	delete from server_messages
	where server_id = $1
	`

	_, err = tx.Exec(messagesQueryString, serverId)
	if err != nil {
		tx.Rollback()
		return server, err
	}

	channelQueryString :=
		`
	delete from channels
	where server_id = $1
	`

	_, err = tx.Exec(channelQueryString, serverId)
	if err != nil {
		tx.Rollback()
		return server, err
	}

	serverQueryString :=
		`
	delete
	from servers
	where server_id = $1
	returning
	server_id,
	server_name,
	server_img
	`

	err = tx.QueryRow(serverQueryString, serverId).Scan(
		&server.Id,
		&server.Name,
		&server.Icon,
	)
	if err != nil {
		tx.Rollback()
		return server, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()

		return server, err
	}

	// for res.Next() {
	// 	res.Scan(
	// 		&server.Id,
	// 		&server.Name,
	// 		&server.Icon,
	// 	)
	// }

	return server, nil
}

func (s ServerModel) GetOwner(serverId, userId int) (int, error) {
	var owner_id int
	queryString :=
		`
	select owner_id
	from servers
	where server_id = $1 and owner_id = $2
	`

	err := s.DB.QueryRow(queryString, serverId, userId).Scan(&owner_id)
	if err != nil {
		return owner_id, err
	}
	return owner_id, nil
}
