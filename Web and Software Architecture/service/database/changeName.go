package database

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) ChangeName(name string, newname string) error {
	_, err := db.c.Exec("PRAGMA foreign_keys=ON; UPDATE User SET username = ? WHERE username = ?", newname, name)
	return err
}
