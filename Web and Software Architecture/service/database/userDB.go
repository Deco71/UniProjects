package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckUser(username string) error {
	var value string
	err := db.c.QueryRow("SELECT username FROM User WHERE username=?", username).Scan(&value)
	return err
}

func (db *appdbimpl) GetFollowed(username string, offset int) ([]string, error) {

	err := db.CheckUser(username)
	if errors.Is(err, sql.ErrNoRows) {
		return make([]string, 0), err
	}

	var values = make([]string, 0)
	rows, err := db.c.Query("SELECT followed FROM Follow WHERE user=? LIMIT 30 OFFSET ?", username, offset*30)
	if err != nil {
		return make([]string, 0), err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, err
}

func (db *appdbimpl) GetFollowers(username string, offset int) ([]string, error) {
	var values = make([]string, 0)
	rows, err := db.c.Query("SELECT user FROM Follow WHERE followed=? LIMIT 30 OFFSET ?", username, offset*30)
	if err != nil {
		return make([]string, 0), err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, err
}

func (db *appdbimpl) Follow(username string, followed string) error {
	_, err := db.c.Exec("PRAGMA foreign_keys=ON; INSERT INTO Follow (user, followed) VALUES (?, ?)", username, followed)
	return err
}

func (db *appdbimpl) Unfollow(username string, followed string) error {
	_, err := db.c.Exec("PRAGMA foreign_keys=ON; DELETE FROM Follow WHERE user = ? AND followed = ?", username, followed)
	return err
}

func (db *appdbimpl) Ban(username string, banned string) error {
	_, err := db.c.Exec("PRAGMA foreign_keys=ON; INSERT INTO Ban (user, banned) VALUES (?, ?)", username, banned)
	return err
}

func (db *appdbimpl) Unban(username string, banned string) error {
	_, err := db.c.Exec("PRAGMA foreign_keys=ON; DELETE FROM Ban WHERE user = ? AND banned = ?", username, banned)
	return err
}

func (db *appdbimpl) GetBan(username string, offset int) ([]string, error) {
	var values = make([]string, 0)
	rows, err := db.c.Query("SELECT banned FROM Ban WHERE user=? LIMIT 30 OFFSET ?", username, offset*30)
	if err != nil {
		return make([]string, 0), err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, err
}
