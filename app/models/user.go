package models

import (
	"strconv"
	"time"

	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserPrivilegeNormal = iota
	UserPrivilegeOrganizer
	UserPrivilegeSuperuser
)

type User struct {
	// Authorization
	Id       int32
	Handle   string
	Email    string
	Password string

	// Privilege
	Privilege int8

	// Miscellaneous statistics
	JoinedAt int64

	// Customized profile
	Nickname string
	Bio      string
}

func init() {
	registerSchema("users",
		"id SERIAL PRIMARY KEY",
		"handle TEXT NOT NULL",
		"email TEXT NOT NULL",
		"password TEXT NOT NULL",
		"privilege SMALLINT NOT NULL DEFAULT "+strconv.Itoa(UserPrivilegeNormal),
		"joined_at BIGINT NOT NULL",
		"nickname TEXT NOT NULL",
		"bio TEXT NOT NULL DEFAULT ''",
	)
}

func (u *User) hashPassword() {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost+1)
	if err != nil {
		panic(err)
	}
	u.Password = string(hashed)
}

func (u *User) Create() error {
	u.hashPassword()

	u.JoinedAt = time.Now().Unix()
	// lib/pq driver does not support LastInsertId()
	// https://github.com/lib/pq/issues/24
	err := db.QueryRow("INSERT INTO "+
		"users(handle, email, password, joined_at, nickname) "+
		"VALUES ($1, $2, $3, $4, $5) RETURNING id",
		u.Handle,
		u.Email,
		u.Password,
		u.JoinedAt,
		u.Nickname,
	).Scan(&u.Id)
	return err
}

func (u *User) Read(byHandle bool) error {
	var row *sql.Row
	if byHandle {
		row = db.QueryRow("SELECT "+
			"id, handle, email, password, privilege, joined_at, nickname "+
			"FROM users WHERE handle = $1",
			u.Handle,
		)
	} else {
		row = db.QueryRow("SELECT "+
			"id, handle, email, password, privilege, joined_at, nickname "+
			"FROM users WHERE id = $1",
			u.Id,
		)
	}
	err := row.Scan(&u.Id, &u.Handle, &u.Email,
		&u.Password, &u.Privilege, &u.JoinedAt, &u.Nickname)
	return err
}

func (u *User) Update() error {
	u.hashPassword()

	_, err := db.Exec("UPDATE users SET "+
		"handle = $1, email = $2, password = $3, privilege = $4, nickname = $5 "+
		"WHERE id = $1",
		u.Handle,
		u.Email,
		u.Password,
		u.Privilege,
		u.Nickname,
	)
	return err
}

func (u *User) VerifyPassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	return err == nil
}
