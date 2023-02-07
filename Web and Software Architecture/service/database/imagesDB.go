package database

import (
	"strconv"
)

// GetName is an example that shows you how to query data
func (db *appdbimpl) NewImage(username string) (string, error) {

	tr, err := db.c.Begin()
	if err != nil {
		return "", err
	}

	_, err = tr.Exec("PRAGMA foreign_keys=ON; INSERT INTO Photo (user) VALUES (?)", username)
	if err != nil {
		return "", err
	}

	var id string
	err = tr.QueryRow("SELECT LAST_INSERT_ROWID()").Scan(&id)
	if err != nil {
		return "", err
	}
	err = tr.Commit()
	if err != nil {
		return "", err
	}

	return id, err
}

func (db *appdbimpl) DeleteImage(id string) error {
	_, err := db.c.Exec("PRAGMA foreign_keys=ON; DELETE FROM Photo WHERE id = ?", id)
	return err
}

func (db *appdbimpl) CheckImage(id string) (string, error) {
	var corrisp string
	err := db.c.QueryRow("SELECT user FROM Photo WHERE id = ?", id).Scan(&corrisp)
	return corrisp, err
}

func (db *appdbimpl) CheckBan(user string, banned string) (int, error) {
	var corrisp string
	err := db.c.QueryRow("SELECT count(*) FROM Ban WHERE user = ? AND banned = ?", user, banned).Scan(&corrisp)

	var count int
	count, _ = strconv.Atoi(corrisp)

	return count, err
}

func (db *appdbimpl) GetImages(user string, offset int) ([]string, error) {
	var corrisp = make([]string, 0)
	rows, err := db.c.Query("SELECT id FROM Photo WHERE user = ? ORDER BY created_at DESC LIMIT 30 OFFSET ?", user, offset*30)
	if err != nil {
		return nil, err
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

		corrisp = append(corrisp, value)
	}

	return corrisp, err
}
