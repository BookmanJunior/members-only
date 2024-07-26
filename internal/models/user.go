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

func (u *UserModel) Insert(username, password string, avatar int) (User, error) {
	user := User{}
	var userId int
	res := u.DB.QueryRow(`insert into users (username, password, avatar) values ($1, $2, $3) returning id`, username, password, avatar)

	err := res.Scan(&userId)

	if err != nil {
		return user, err
	}

	user, err = u.Get(userId)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserModel) Get(id int) (User, error) {
	user := User{}
	res := u.DB.QueryRow(`select "username", "password", "avatar_url" from "users" inner join "avatars" on users.avatar = avatars.id and users.id = $1`, id)

	err := res.Scan(&user.Username, &user.password, &user.Avatar)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserModel) GetByUsername(username string) (string, error) {
	var usernameExists string
	res := u.DB.QueryRow(`select "username" from "users" where "username" = $1`, username)

	err := res.Scan(&usernameExists)

	if err != nil {
		return usernameExists, err
	}

	return usernameExists, nil

}
