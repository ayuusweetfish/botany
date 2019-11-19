package main

var schema = `
CREATE TABLE IF NOT EXISTS b_user (
	uid SERIAL CONSTRAINT uid_unique UNIQUE,
	username TEXT UNIQUE,
	password TEXT,
	email TEXT,
  	last_login timestamp NOT NULL DEFAULT NOW(),
	count int,
	PRIMARY KEY(uid)
);
CREATE TABLE IF NOT EXISTS b_game (
	gid SERIAL CONSTRAINT gid_unique UNIQUE,
	owner int NOT NULL REFERENCES b_user(uid) ON DELETE CASCADE,
	gameName TEXT,
	beginTime timestamp NOT NULL,
	endTime timestamp NOT NULL,
	gameInfo TEXT,
	PRIMARY KEY(gid)
);
CREATE TABLE IF NOT EXISTS managers_games (
	user_id INTEGER NOT NULL,
	game_id INTEGER NOT NULL,
	PRIMARY KEY(user_id, game_id),
	FOREIGN KEY (user_id) REFERENCES b_user(uid) ON UPDATE CASCADE,
	FOREIGN KEY (game_id) REFERENCES b_game(gid) ON UPDATE CASCADE
);
CREATE TABLE IF NOT EXISTS players_games (
	user_id INTEGER NOT NULL,
	game_id INTEGER NOT NULL,
	PRIMARY KEY(user_id, game_id),
	FOREIGN KEY (user_id) REFERENCES b_user(uid) ON UPDATE CASCADE,
 	FOREIGN KEY (game_id) REFERENCES b_game(gid) ON UPDATE CASCADE
);
CREATE TABLE IF NOT EXISTS b_match (
	playerA integer REFERENCES b_user(uid),
	playerB integer REFERENCES b_user(uid),
	winner integer REFERENCES b_user(uid),
	game integer REFERENCES b_game(gid)
);
CREATE TABLE IF NOT EXISTS websiteInfo (
	gameNumber integer DEFAULT 0,
	userNumber integer DEFAULT 0
)
`
