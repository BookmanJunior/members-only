package models

import (
	"database/sql"
)

type Avatar struct {
	id   int
	Name string
	Url  string
}

type AvatarModel struct {
	DB *sql.DB
}

func (a *Avatar) Insert(db *sql.DB) func() error {
	return func() error {
		_, err := db.Exec(`insert into avatars("name", "url") values($1, $2) returning id`, a.Name, a.Url)

		if err != nil {
			return err
		}

		return nil
	}
}

func (a *AvatarModel) Get(id int) (Avatar, error) {
	avatar := Avatar{}

	res := a.DB.QueryRow(`select * from "avatars" where "id" = $1`, id)

	err := res.Scan(&avatar.Name, &avatar.Url, &avatar.id)

	if err != nil {
		return avatar, err
	}

	return avatar, nil
}

func (a *AvatarModel) Exists(id int) bool {
	_, err := a.Get(id)

	if err != nil {
		return false
	}

	return true
}
