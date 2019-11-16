package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func createUser(name string, password string) bool {
	var count int
	var err error
	row := db.QueryRow("SELECT username = $1 from b_user", name)
	err = row.Scan(&count)
	if err == sql.ErrNoRows {
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		_, err = db.Exec("INSERT INTO b_user(username, password) VALUES ($1, $2)", name, hashedPwd)
		if err != nil {
			log.Println(err)
			return false
		}
	}
	return true
}

func verifyPassword(password string, hashedPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func updataUser(uid int, item string, args ...string) bool {
	var err error
	err = nil
	if item == "username" {
		_, err = db.Exec("UPDATE b_user SET username = $1 WHERE uid = $2", args[0], uid)
		if err != nil {
			log.Println(err)
			return false
		}
	} else if item == "password" {
		hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(args[0]), bcrypt.MinCost)
		_, err = db.Exec("UPDATE b_user SET password = $1 WHERE uid = $2", hashedPwd, uid)
		if err != nil {
			log.Println(err)
			return false
		}
	}
	return true
}
