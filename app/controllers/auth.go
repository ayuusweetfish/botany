package controllers

import (
	"github.com/kawa-yoiko/botany/app/globals"
	"github.com/kawa-yoiko/botany/app/models"

	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/gorilla/mux"
)

// Returns a user ID
func middlewareAuthRetrieve(w http.ResponseWriter, r *http.Request) int32 {
	session, err := globals.SessionStore.Get(r, "auth")
	if err != nil {
		panic(err)
	}

	id, ok := session.Values["uid"].(int32)
	if ok {
		return id
	} else {
		return -1
	}
}

func middlewareAuthGrant(w http.ResponseWriter, r *http.Request, uid int32) {
	session, err := globals.SessionStore.Get(r, "auth")
	if err != nil {
		panic(err)
	}

	session.Values["uid"] = uid
	err = session.Save(r, w)
	if err != nil {
		panic(err)
	}
}

// curl http://localhost:3434/signup --data "handle=abc&password=qwq"
// TODO: Add captcha validation and more fields
func signupHandler(w http.ResponseWriter, r *http.Request) {
	s := r.PostFormValue("handle")
	p := r.PostFormValue("password")

	u := models.User{}
	u.Handle = s
	u.Email = "beta@example.com"

	// TODO: Look for an existing user with the same handle or email

	u.Password = p
	u.Nickname = "Uvuvwevwevwe"
	if err := u.Create(); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Welcome on board, %s!\n", u.Handle)
}

// curl http://localhost:3434/login -i --data "handle=abc&password=qwq"
func loginHandler(w http.ResponseWriter, r *http.Request) {
	s := r.PostFormValue("handle")
	p := r.PostFormValue("password")

	// TODO: Support logging in with e-mail
	u := models.User{}
	u.Handle = s

	ok := true

	if err := u.Read(true); err != nil {
		if err == sql.ErrNoRows {
			ok = false
		} else {
			panic(err)
		}
	}

	if ok && !u.VerifyPassword(p) {
		ok = false
	}

	if ok {
		middlewareAuthGrant(w, r, u.Id)
		fmt.Fprintf(w, "Hi %s, your ID is %d\n", s, u.Id)
	} else {
		fmt.Fprintf(w, "Incorrect credentials\n")
	}
}

// XXX: For debug use
// curl http://localhost:3434/check_auth -i -H "Cookie: auth=..."
func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	uid := middlewareAuthRetrieve(w, r)
	fmt.Fprintf(w, "Your ID is %d\n", uid)
}

func init() {
	registerRouterFunc("/signup", signupHandler, "POST")
	registerRouterFunc("/login", loginHandler, "POST")
	registerRouterFunc("/check_auth", checkAuthHandler, "GET")
}
