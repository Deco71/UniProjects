-- to launch use sqlite3 wasadb.db -init config.sql
PRAGMA foreign_keys = ON;
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
