package models

import (
	"database/sql"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return user, err
	}

	res := u.DB.QueryRow(`insert into users (username, password, avatar) values ($1, $2, $3) returning id`, strings.ToLower(username), hashedPassword, avatar)

	err = res.Scan(&userId)

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

func (u *UserModel) GetByUsername(username string) (User, error) {
	var user User

	res := u.DB.QueryRow(`select "username" from "users" where "username" = $1`, strings.ToLower(username))

	err := res.Scan(&user.Username)

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
