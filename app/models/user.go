package models

import "strconv"

const (
	UserPrivilegeSuperuser = iota
	UserPrivilegeOrganizer
	UserPrivilegeNormal
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
	)
}

func (u *User) Create() error {
	// TODO
	return nil
}

func (u *User) Read() error {
	// TODO
	return nil
}

func (u *User) Update() error {
	// TODO
	return nil
}
