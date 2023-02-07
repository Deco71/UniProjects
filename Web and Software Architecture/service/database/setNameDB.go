package database

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) SetName(name string) error {
	_, err := db.c.Exec("PRAGMA foreign_keys=ON; INSERT INTO User (username) VALUES (?)", name)
	return err
}
