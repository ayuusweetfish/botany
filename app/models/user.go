package models

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	Uid      int
	Name     string
	Password string
	Email    string
}

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

func (u *User) Create() error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	_, err = db.Exec("INSERT INTO b_user(username, password, email) VALUES ($1, $2, $3)", u.Name, hashedPwd, u.Email)
	return err
}

func (u *User) Read() error {
	row := db.QueryRow("SELECT uid FROM b_user WHERE username = $1", u.Name)
	err := row.Scan(&u.Uid)
	return err
}

func (u *User) Update() error {
	_, err := db.Exec("UPDATE b_user SET username = $1 , email = $2 WHERE uid = $3", u.Name, u.Email, u.Uid)
	return err
}

func (u *User) UpdatePassword() error {
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	_, err := db.Exec("UPDATE b_user SET password = $1 WHERE uid = $2", hashedPwd, u.Uid)
	return err
}

func (u *User) VerifyPassword() bool {
	var userPassword string
	row := db.QueryRow("SELECT password FROM b_user WHERE username = $1", u.Name)
	row.Scan(&userPassword)
	byteHash := []byte(userPassword)
	bytePwd := []byte(u.Password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func init() {
	registerSchema(schema)
}
