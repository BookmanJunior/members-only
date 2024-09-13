package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id            int    `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"-"`
	FileSizeLimit int    `json:"file_limit,omitempty"`
	Avatar
	Admin bool `json:"admin,omitempty"`
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(username, password string, avatar int) (int, error) {
	queryString := `insert into users (username, password, avatar) values ($1, $2, $3) returning id`

	var userId int

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)

	if err != nil {
		return userId, err
	}

	res := u.DB.QueryRow(queryString, username, hashedPassword, avatar)

	err = res.Scan(&userId)

	if err != nil {
		return userId, err
	}

	return userId, nil
}

func (u *UserModel) GetById(id int) (User, error) {
	queryString := `
	select
	users.id,
	"username",
	"password",
	"avatar_color",
	"avatar_url",
	"admin",
	"limit_size" from "users"
	inner join "avatars" on users.avatar = avatars.id and users.id = $1
	inner join file_limit on users.file_limit_id = file_limit.limit_id`

	user := User{}

	res := u.DB.QueryRow(queryString, id)

	err := res.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Avatar.Color,
		&user.Avatar.Url,
		&user.Admin,
		&user.FileSizeLimit)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserModel) GetByUsername(username string) (User, error) {
	queryString := `
	select
	users.id,
	"username",
	"password",
	"avatar_color",
	"avatar_url",
	"admin",
	"limit_size" from "users"
	inner join "avatars" on users.avatar = avatars.id and lower(username) = $1
	inner join file_limit on users.file_limit_id = file_limit.limit_id`

	var user User

	res := u.DB.QueryRow(queryString, username)

	err := res.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Avatar.Color,
		&user.Avatar.Url,
		&user.Admin,
		&user.FileSizeLimit)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserModel) Exists(username string) bool {
	_, err := u.GetByUsername(username)

	return err == nil
}
