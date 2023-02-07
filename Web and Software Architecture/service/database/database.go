/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.
To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.
For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetToken(name string) (int, error)
	SetName(name string) error
	CheckToken(token string) (string, error)
	ChangeName(name string, newname string) error
	NewImage(username string) (string, error)
	DeleteImage(id string) error
	CheckImage(id string) (string, error)
	CheckBan(user string, banned string) (int, error)
	GetImages(user string, offset int) ([]string, error)
	CheckUser(username string) error
	GetFollowed(username string, offset int) ([]string, error)
	GetFollowers(username string, offset int) ([]string, error)
	Follow(username string, followed string) error
	Unfollow(username string, followed string) error
	Ban(username string, banned string) error
	Unban(username string, banned string) error
	GetBan(username string, offset int) ([]string, error)
	Like(post string, username string) error
	Unlike(post string, username string) error
	Comment(post string, username string, comment string) (string, string, error)
	Uncomment(post string) error
	GetComment(comment string) (string, string, string, string, string, error)
	GetComments(post string, requester string, offset int) ([]string, []string, []string, []string, error)
	GetLikes(post string, requester string, offset int) ([]string, error)
	GetPost(post string, requester string) (string, string, string, string, string, error)
	GetFeed(requester string, offset int) ([]string, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='User';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `PRAGMA foreign_keys = ON;
		BEGIN TRANSACTION;
		
		create table User (
			username TEXT UNIQUE,  
			token INTEGER PRIMARY KEY AUTOINCREMENT
		);
		
		create table Photo (
			id INTEGER PRIMARY KEY,
			user TEXT NOT NULL REFERENCES User(username) ON DELETE CASCADE ON UPDATE CASCADE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		create table Like (
			photo INTEGER NOT NULL REFERENCES Photo(id) ON DELETE CASCADE ON UPDATE CASCADE,
			user TEXT NOT NULL REFERENCES User(username) ON DELETE CASCADE ON UPDATE CASCADE,
			PRIMARY KEY (user, photo)
		);
		
		create table Comment (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			photo INTEGER NOT NULL REFERENCES Photo(id) ON DELETE CASCADE ON UPDATE CASCADE,
			user TEXT NOT NULL REFERENCES User(username) ON DELETE CASCADE ON UPDATE CASCADE,
			comment TEXT NOT NULL,
			created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		create table Follow (
			user TEXT NOT NULL REFERENCES User(username) ON DELETE CASCADE ON UPDATE CASCADE,
			followed TEXT NOT NULL REFERENCES User(username) ON DELETE CASCADE ON UPDATE CASCADE,
			PRIMARY KEY (user, followed)
		);
		
		create table Ban (
			user TEXT NOT NULL REFERENCES User(username) ON DELETE CASCADE,
			banned TEXT NOT NULL REFERENCES User(username) ON DELETE CASCADE,
			PRIMARY KEY (user, banned)
		);
		COMMIT;
		`

		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	_, err = db.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		return nil, fmt.Errorf("error enabling foreign keys: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
