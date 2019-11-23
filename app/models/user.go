package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uid      int
	Name     string
	Password string
	Email    string
	Count    int
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
`

func (u *User) Create() error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	_, err = db.Exec("INSERT INTO b_user(username, password, email, count) VALUES ($1, $2, $3, $4)", u.Name, hashedPwd, u.Email, u.Count)
	return err
}

func (u *User) Read() error {
	row := db.QueryRow("SELECT uid , username, email, count FROM b_user WHERE username = $1", u.Name)
	err := row.Scan(&u.Uid, &u.Name, &u.Email, &u.Count)
	return err
}

//updata the other fields except the password, before update you need to Read() first
func (u *User) Update() error {
	_, err := db.Exec("UPDATE b_user SET username = $1, email = $2, count = $3 WHERE uid = $4", u.Name, u.Email, u.Count, u.Uid)
	return err
}

func (u *User) UpdatePassword() error {
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	_, err := db.Exec("UPDATE b_user SET password = $1 WHERE uid = $2", hashedPwd, u.Uid)
	return err
}

//you need to make sure the user exists before verify the password
func (u *User) VerifyPassword() bool {
	var userPassword string
	row := db.QueryRow("SELECT password FROM b_user WHERE username = $1", u.Name)
	err := row.Scan(&userPassword)
	if err != nil {
		panic(err)
	}
	byteHash := []byte(userPassword)
	bytePwd := []byte(u.Password)
	err = bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		panic(err)
	}
	return true
}

func init() {
	registerSchema(schema)
}
