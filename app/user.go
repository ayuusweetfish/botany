package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)
	id, err := r.Cookie("QAQ")
	if err == http.ErrNoCookie {
		//w.WriteHeader()
		//？遗留问题暂待处理
	}
	if r.Method == "GET" {
		_, captcha := captchaCreate(id.Value)
		result := []byte(fmt.Sprintf(`{"pic": %s}`, captcha))
		w.Write(result)
		return
	} else if r.Method == "POST" {
		vars := mux.Vars(r)
		username := vars["username"]
		password := vars["password"]
		captcha := vars["captcha"]
		captchaPass := captchaVerfiy(id.Value, captcha)
		if captchaPass {
			var userPassword string
			row := db.QueryRow("SELECT count FROM b_user WHERE username = $1", username)
			err := row.Scan(&userPassword)
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "该用户不存在"}`))
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal Server Error"}`))
				return
			} else {
				passwordPass := verifyPassword(password, userPassword)
				if passwordPass {
					w.Write([]byte(`{"success": "登录成功"}`))
				} else {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(`{"error": "密码错误"}`))
				}
			}
		}
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)
	//vars := mux.Vars(r)
	if r.Method == "GET" {
		//_, captcha := captchaCreate()
		//w.Write([]byte(captcha))
	}
}
