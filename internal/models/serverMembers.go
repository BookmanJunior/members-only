package models

import (
	"database/sql"
)

type ServerMembersModel struct {
	DB *sql.DB
}

type ServerMembers struct {
	serverId int
	userId   int
}

func (s ServerMembersModel) IsAllowed(serverId, userId int) (bool, error) {
	var isAllowed bool
	queryString :=
		`
	select
	case when po.type = 'public' then true
	when po.type = 'private' and sm.user_id is not null then true
	else false
	end as isAllowed
	from servers s
	join privacy_options po on po.id = s.type
	left join server_members sm on sm.server_id = $1 and sm.user_id = $2
	where s.server_id = $1;
	`

	err := s.DB.QueryRow(queryString, serverId, userId).Scan(&isAllowed)
	if err != nil {
		return isAllowed, err
	}

	return isAllowed, nil
}

func (s ServerMembersModel) IsUserInServer(serverId, userId int) error {
	var ser ServerMembers
	queryString :=
		`
	select server_id, user_id from server_members
	where server_id = $1 and user_id = $2
	`

	err := s.DB.QueryRow(queryString, serverId, userId).Scan(&ser.serverId, &ser.userId)
	if err != nil {
		return err
	}
	return nil
}

func (s ServerMembersModel) Insert(serverId, userId int) (int, error) {
	var sId int
	queryString :=
		`
	insert into server_members (server_id, user_id)
	values($1, $2)
	returning server_id
	`
	err := s.DB.QueryRow(queryString, serverId, userId).Scan(&sId)
	if err != nil {
		return 0, err
	}

	return sId, nil
}

func (s ServerMembersModel) DeleteByUserId(serverId, userId int) error {
	queryString :=
		`
	delete from server_members
	where server_id = $1 and user_id = $2
	`

	_, err := s.DB.Exec(queryString, serverId, userId)
	if err != nil {
		return err
	}
	return nil
}
