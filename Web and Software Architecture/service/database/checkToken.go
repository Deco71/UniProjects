package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) CheckToken(token string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT username FROM User WHERE token=?", token).Scan(&name)
	return name, err
}
