package database

func (db *appdbimpl) GetToken(username string) (int, error) {
	var token int
	err := db.c.QueryRow("SELECT token FROM User WHERE username=?", username).Scan(&token)
	return token, err
}
