package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Avatar   string `json:"avatar"`
	Admin    bool   `json:"admin"`
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(username, password string, avatar int) (int, error) {
	var userId int
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return userId, err
	}

	res := u.DB.QueryRow(`insert into users (username, password, avatar) values ($1, $2, $3) returning id`, username, hashedPassword, avatar)

	err = res.Scan(&userId)

	if err != nil {
		return userId, err
	}

	return userId, nil
}

func (u *UserModel) GetById(id int) (User, error) {
	user := User{}
	res := u.DB.QueryRow(`select users.id, "username", "password", "avatar_url", "admin" from "users" inner join "avatars" on users.avatar = avatars.id and users.id = $1`, id)

	err := res.Scan(&user.Id, &user.Username, &user.Password, &user.Avatar, &user.Admin)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserModel) GetByUsername(username string) (User, error) {
	var user User

	res := u.DB.QueryRow(`select "id", "username", "password", "avatar", "admin" from "users" where "username" = $1`, username)

	err := res.Scan(&user.Id, &user.Username, &user.Password, &user.Avatar, &user.Admin)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserModel) Exists(username string) bool {
	_, err := u.GetByUsername(username)

	if err != nil {
		return false
	}

	return true
}
