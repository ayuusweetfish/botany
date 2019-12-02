package controllers

import (
	"github.com/kawa-yoiko/botany/app/globals"
	"github.com/kawa-yoiko/botany/app/models"

	"database/sql"
	"encoding/json"
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

func middlewareCaptchaVerify(w http.ResponseWriter, r *http.Request) bool {
	captchaKey := r.PostFormValue("captcha_key")
	captchaValue := r.PostFormValue("captcha_value")
	return globals.CaptchaVerfiy(captchaKey, captchaValue)
}

// curl http://localhost:3434/signup -i -d "handle=abc&password=qwq&email=gamma@example.com&nickname=ABC&captcha_key=...&captcha_value=..."
func signupHandler(w http.ResponseWriter, r *http.Request) {
	if !middlewareCaptchaVerify(w, r) {
		w.WriteHeader(400)
		fmt.Fprintf(w, "{\"err\": [3]}")
		return
	}

	s := r.PostFormValue("handle")
	p := r.PostFormValue("password")
	email := r.PostFormValue("email")
	nickname := r.PostFormValue("nickname")

	u := models.User{}
	u.Handle = s
	u.Email = email

	// TODO: Validate email format
	// TODO: Look for an existing user with the same handle or email

	u.Password = p
	u.Nickname = nickname
	if err := u.Create(); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "{\"err\": []}")
}

// curl http://localhost:3434/login -i -d "handle=abc&password=qwq"
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
		w.WriteHeader(200)
		fmt.Fprintf(w, "{}")
	} else {
		w.WriteHeader(400)
		fmt.Fprintf(w, "{}")
	}
}

// XXX: For debug use
// curl http://localhost:3434/check_auth -i -H "Cookie: auth=..."
func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	uid := middlewareAuthRetrieve(w, r)
	if uid == -1 {
		w.WriteHeader(401)
	} else {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Your ID is %d\n", uid)
	}
}

func captchaGetHandler(w http.ResponseWriter, r *http.Request) {
	key, img := globals.CaptchaCreate()
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]string{
		"key": key,
		"img": img,
	})
}

func init() {
	registerRouterFunc("/signup", signupHandler, "POST")
	registerRouterFunc("/login", loginHandler, "POST")
	registerRouterFunc("/check_auth", checkAuthHandler, "GET")

	registerRouterFunc("/captcha", captchaGetHandler, "GET")
}
