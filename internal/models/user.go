package models

import (
	"database/sql"
)

type User struct {
	Username string
	password string
	Avatar   string
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(username, password string, avatar uint) {
	u.DB.Exec(`insert into users (username, password, avatar) values ($1, $2, $3)`, username, password, avatar)
}

func (u *UserModel) Get(id int) (User, error) {
	user := User{}
	res := u.DB.QueryRow(`select "username", "password", "avatar_url" from "users" inner join "avatars" on users.avatar = avatars.id and users.id = $1`, id)

	if res.Err() != nil {
		return user, res.Err()
	}

	err := res.Scan(&user.Username, &user.password, &user.Avatar)

	if err != nil {
		return user, err
	}

	return user, nil
}
