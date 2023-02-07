package database

func (db *appdbimpl) Like(post string, username string) error {
	_, err := db.c.Exec("PRAGMA foreign_keys=ON; INSERT INTO Like (photo, user) VALUES (?, ?)", post, username)
	return err
}

func (db *appdbimpl) Unlike(post string, username string) error {
	_, err := db.c.Exec("PRAGMA foreign_keys=ON; DELETE FROM Like WHERE photo = ? AND user = ?", post, username)
	return err
}

func (db *appdbimpl) GetLikes(post string, requester string, offset int) ([]string, error) {
	var corrisp = make([]string, 0)
	rows, err := db.c.Query("SELECT l1.user FROM Like l1 WHERE photo = ? AND ? NOT IN (SELECT banned FROM Ban where user = l1.user) LIMIT 30 OFFSET ?", post, requester, offset*30)
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

func (db *appdbimpl) GetComments(post string, requester string, offset int) ([]string, []string, []string, []string, error) {
	var ids = make([]string, 0)
	var users = make([]string, 0)
	var comments = make([]string, 0)
	var dates = make([]string, 0)
	rows, err := db.c.Query("SELECT c1.id, c1.user, c1.comment, c1.created FROM Comment c1 WHERE c1.photo = ? AND ? NOT IN (SELECT banned FROM ban where Ban.user = c1.user) LIMIT 30 OFFSET ?", post, requester, offset*30)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var id string
		var user string
		var comment string
		var created string
		err = rows.Scan(&id, &user, &comment, &created)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		ids = append(ids, id)
		users = append(users, user)
		comments = append(comments, comment)
		dates = append(dates, created)
	}

	return ids, users, comments, dates, err
}

func (db *appdbimpl) Comment(post string, username string, comment string) (string, string, error) {
	var data string

	hi, err := db.c.Begin()
	if err != nil {
		return "", "", err
	}

	err = hi.QueryRow("PRAGMA foreign_keys=ON; INSERT INTO Comment (photo, user, comment) VALUES (?, ?, ?) RETURNING created", post, username, comment).Scan(&data)
	if err != nil {
		return "", "", err
	}
	var id string
	err = hi.QueryRow("SELECT LAST_INSERT_ROWID()").Scan(&id)
	if err != nil {
		return "", "", err
	}

	err = hi.Commit()
	if err != nil {
		return "", "", err
	}

	return id, data, err
}

func (db *appdbimpl) Uncomment(comment string) error {
	_, err := db.c.Exec("PRAGMA foreign_keys=ON; DELETE FROM Comment WHERE id = ?", comment)
	return err
}

func (db *appdbimpl) GetComment(comment string) (string, string, string, string, string, error) {
	var owner string
	var photo string
	var post string
	var commento string
	var created string
	err := db.c.QueryRow("SELECT id, photo, user, comment, created FROM Comment WHERE id = ?", comment).Scan(&photo, &post, &owner, &commento, &created)
	return photo, post, owner, commento, created, err
}

func (db *appdbimpl) GetPost(post string, requester string) (string, string, string, string, string, error) {
	var username string
	var likeCount string
	var commentCount string
	var created_at string
	var liked string

	err := db.c.QueryRow(`
	SELECT p1.user, p1.created_at 
	FROM Photo p1
	WHERE p1.id = ? AND ? NOT IN (SELECT banned FROM Ban WHERE Ban.user = p1.user)
	GROUP BY p1.id`, post, requester).Scan(&username, &created_at)

	_ = db.c.QueryRow("SELECT COUNT(l1.user) FROM Like l1 WHERE l1.photo = ? AND ? NOT IN (SELECT banned FROM Ban WHERE Ban.user = l1.user)", post, requester).Scan(&likeCount)
	_ = db.c.QueryRow("SELECT COUNT(l1.user) FROM Comment l1 WHERE l1.photo = ? AND ? NOT IN (SELECT banned FROM Ban WHERE Ban.user = l1.user)", post, requester).Scan(&commentCount)
	_ = db.c.QueryRow("SELECT EXISTS(SELECT l1.user FROM Like l1 WHERE l1.photo = ? AND l1.user = ?)", post, requester).Scan(&liked)

	return username, likeCount, commentCount, created_at, liked, err
}

func (db *appdbimpl) GetFeed(requester string, offset int) ([]string, error) {
	var corrisp = make([]string, 0)
	rows, err := db.c.Query("SELECT DISTINCT p1.id FROM Photo p1 WHERE p1.user=? OR (p1.user IN (SELECT f1.followed FROM Follow f1 WHERE f1.user = ? AND ? NOT IN (SELECT b1.banned FROM Ban b1 WHERE b1.user = f1.followed))) ORDER BY p1.created_at DESC LIMIT 30 OFFSET ?", requester, requester, requester, offset*30)
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
