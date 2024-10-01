package models

import (
	"database/sql"
)

type Avatar struct {
	id    int
	Color string `json:"avatar_color"`
	Url   string `json:"avatar"`
}

type AvatarModel struct {
	DB *sql.DB
}

func (a *AvatarModel) Insert(avatar_color, avatar_url string) error {
	queryString := `insert into avatars("name", "url") values($1, $2)`
	_, err := a.DB.Exec(queryString, avatar_color, avatar_url)
	if err != nil {
		return err
	}

	return nil
}

func (a *AvatarModel) Get(id int) (Avatar, error) {
	avatar := Avatar{}
	queryString := `select "avatar_color", "avatar_url", "avatar_id" from "avatars" where "avatar_id" = $1`
	res := a.DB.QueryRow(queryString, id)
	err := res.Scan(&avatar.Color, &avatar.Url, &avatar.id)

	if err != nil {
		return avatar, err
	}

	return avatar, nil
}

func (a *AvatarModel) GetAll() ([]Avatar, error) {
	var avatars []Avatar
	queryString := `select "avatar_color", "avatar_url", "avatar_id" from avatars`
	res, err := a.DB.Query(queryString)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		avatar := Avatar{}
		res.Scan(&avatar.Color, &avatar.Url, &avatar.id)
		avatars = append(avatars, avatar)
	}

	return avatars, nil
}

func (a *AvatarModel) Exists(id int) bool {
	_, err := a.Get(id)

	return err == nil
}
