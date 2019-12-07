package models

import (
	"regexp"
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

func (u *User) Representation() map[string]interface{} {
	return map[string]interface{}{
		"id":        u.Id,
		"handle":    u.Handle,
		"email":     u.Email,
		"privilege": u.Privilege,
		"joined_at": u.JoinedAt,
		"nickname":  u.Nickname,
		"bio":       u.Bio,
	}
}

func (u *User) ShortRepresentation() map[string]interface{} {
	return map[string]interface{}{
		"id":        u.Id,
		"handle":    u.Handle,
		"privilege": u.Privilege,
		"nickname":  u.Nickname,
	}
}

func (u *User) hashPassword() {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost-3)
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
		"users(handle, email, password, privilege, joined_at, nickname) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		u.Handle,
		u.Email,
		u.Password,
		u.Privilege,
		u.JoinedAt,
		u.Nickname,
	).Scan(&u.Id)
	return err
}

func (u *User) read(field string) error {
	var row *sql.Row
	if field == "handle" {
		row = db.QueryRow("SELECT "+
			"id, handle, email, password, privilege, joined_at, nickname, bio"+
			"FROM users WHERE handle = $1",
			u.Handle,
		)
	} else if field == "id" {
		row = db.QueryRow("SELECT "+
			"id, handle, email, password, privilege, joined_at, nickname, bio "+
			"FROM users WHERE id = $1",
			u.Id,
		)
	} else if field == "email" {
		row = db.QueryRow("SELECT "+
			"id, handle, email, password, privilege, joined_at, nickname, bio "+
			"FROM users WHERE email = $1",
			u.Email,
		)
	}
	err := row.Scan(&u.Id, &u.Handle, &u.Email,
		&u.Password, &u.Privilege, &u.JoinedAt, &u.Nickname, &u.Bio)
	return err
}

func (u *User) ReadById() error {
	return u.read("id")
}

func (u *User) ReadByHandle() error {
	return u.read("handle")
}

func (u *User) ReadByEmail() error {
	return u.read("email")
}

func (u *User) Update() error {
	_, err := db.Exec("UPDATE users SET "+
		"handle = $1, email = $2, privilege = $3, nickname = $4, bio = $5"+
		"WHERE id = $6",
		u.Handle,
		u.Email,
		u.Privilege,
		u.Nickname,
		u.Bio,
		u.Id,
	)
	return err
}

func (u *User) UpdatePassword() error {
	u.hashPassword()
	_, err := db.Exec("UPDATE users SET password = $1 WHERE id = $2", u.Password, u.Id)
	return err
}

func (u *User) VerifyPassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	return err == nil
}

func (u *User) EmailCheck() bool {
	// Now it is not complete because there are some situations this one cannot handle.
	// For example the email .list@gmail.com or list..list@gmail.com is not correct according to RFC 5322.

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(u.Email)
}

func UserSearchByHandle(handle string) ([]User, error) {
	rows, err := db.Query("SELECT "+
		"id, handle, privilege, nickname "+
		"FROM users WHERE handle LIKE '%' || $1 || '%' LIMIT 5",
		handle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	us := []User{}
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.Id, &u.Handle, &u.Privilege, &u.Nickname)
		if err != nil {
			return nil, err
		}
		us = append(us, u)
	}
	return us, rows.Err()
}
