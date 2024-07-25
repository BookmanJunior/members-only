package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Avatar struct {
	id   int
	Name string
	Url  string
}

func (a *Avatar) Insert(db *sql.DB) func() error {
	return func() error {
		if a.Url == "" || a.Name == "" {
			return errors.New("Name or url can't be empty")
		}

		_, err := db.Exec(`insert into avatars("name", "url") values($1, $2) returning id`, a.Name, a.Url)

		if err != nil {
			fmt.Println(err)
			return errors.New("Couldn't complete operation")
		}

		fmt.Printf("Added %v to avatars table\n", a.Name)
		return nil
	}
}
