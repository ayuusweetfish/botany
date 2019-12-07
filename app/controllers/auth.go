package controllers

import (
	"github.com/kawa-yoiko/botany/app/globals"
	"github.com/kawa-yoiko/botany/app/models"

	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

// Returns a User struct
func middlewareAuthRetrieve(w http.ResponseWriter, r *http.Request) models.User {
	session, err := globals.SessionStore.Get(r, "auth")
	if err != nil {
		panic(err)
	}

	id, ok := session.Values["uid"].(int32)
	if !ok {
		return models.User{Id: -1}
	}

	u := models.User{Id: id}
	if err := u.ReadById(); err != nil {
		if err == sql.ErrNoRows {
			return models.User{Id: -1}
		} else {
			panic(err)
		}
	}
	return u
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

	// TODO: Validate email format.
	// Now it is not complete because there are some situations this one cannot handle.
	// For example the email .list@gmail.com or list..list@gmail.com is not correct according to RFC 5322.
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(u.Email) {
		w.WriteHeader(400)
		fmt.Fprintf(w, "{\"err\": [-1]}")
		return
	}

	// check whether a user with the same handle or email is existing
	userExist := true
	if err := u.ReadByHandle(); err != nil {
		if err == sql.ErrNoRows {
			userExist = false
		} else {
			panic(err)
		}
	}
	if userExist {
		w.WriteHeader(400)
		fmt.Fprintf(w, "{\"err\": [1]}")
		return
	}
	userExist = true
	if err := u.ReadByEmail(); err != nil {
		if err == sql.ErrNoRows {
			userExist = false
		} else {
			panic(err)
		}
	}
	if userExist {
		w.WriteHeader(400)
		fmt.Fprintf(w, "{\"err\": [2]}")
		return
	}

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

	if err := u.ReadByHandle(); err != nil {
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
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.Encode(u.Representation())
	} else {
		w.WriteHeader(400)
		fmt.Fprintf(w, "{}")
	}
}

// curl http://localhost:3434/logout -i -H "Cookie: auth=..." -d ""
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := globals.SessionStore.Get(r, "auth")
	if err != nil {
		panic(err)
	}

	// Delete session
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		panic(err)
	}
}

// curl http://localhost:3434/whoami -i -H "Cookie: auth=..."
func whoAmIHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.Encode(u.ShortRepresentation())
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

func profileHandler(w http.ResponseWriter, r *http.Request) {
	u := models.User{Handle: mux.Vars(r)["handle"]}

	if err := u.ReadByHandle(); err != nil {
		if err == sql.ErrNoRows {
			// No such user
			w.WriteHeader(404)
			return
		} else {
			panic(err)
		}
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]interface{}{
		"user":     u.Representation(),
		"contests": []int{}, // TODO
		"matches":  []int{}, // TODO
	})
}

func profileEditHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}

	user := models.User{Handle: mux.Vars(r)["handle"]}
	if u.Handle != user.Handle && u.Privilege != models.UserPrivilegeSuperuser {
		// this user is not superuser
		// and the user's handle is not the same with the user whose info is being updated
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	if err := user.ReadByHandle(); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		} else {
			panic(err)
		}
	}

	user.Email = r.PostFormValue("email")
	user.Nickname = r.PostFormValue("nickname")
	user.Bio = r.PostFormValue("bio")

	if err := user.Update(); err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "{}")
}

func passwordEditHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}

	user := models.User{Handle: mux.Vars(r)["handle"]}
	if u.Handle != user.Handle && u.Privilege != models.UserPrivilegeSuperuser {
		// this user is not superuser
		// and the user's handle is not the same with the user whose password is being updated
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	if err := user.ReadByHandle(); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		} else {
			panic(err)
		}
	}

	old := r.PostFormValue("old")
	new := r.PostFormValue("new")

	// the old password is wrong
	if !user.VerifyPassword(old) {
		w.WriteHeader(400)
		fmt.Fprintf(w, "{}")
		return
	}

	user.Password = new
	if err := user.UpdatePassword(); err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "{}")
}

// 赋予或撤回主办权限
func promoteHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}

	if u.Privilege != models.UserPrivilegeSuperuser {
		// this user is not superuser
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	user := models.User{Handle: mux.Vars(r)["handle"]}
	if err := user.ReadByHandle(); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		} else {
			panic(err)
		}
	}

	set := r.PostFormValue("set")
	if set == "true" {
		user.Privilege = models.UserPrivilegeOrganizer
	} else if set == "false" && user.Privilege == models.UserPrivilegeOrganizer {
		user.Privilege = models.UserPrivilegeNormal
	}

	if err := user.Update(); err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "{}")
}

func init() {
	registerRouterFunc("/signup", signupHandler, "POST")
	registerRouterFunc("/login", loginHandler, "POST")
	registerRouterFunc("/logout", logoutHandler, "POST")
	registerRouterFunc("/whoami", whoAmIHandler, "GET")

	registerRouterFunc("/captcha", captchaGetHandler, "GET")

	registerRouterFunc("/user/{handle:[a-zA-Z0-9-_]+}/profile", profileHandler, "GET")
	registerRouterFunc("/user/{handle:[a-zA-Z0-9-_]+}/profile/edit", profileEditHandler, "POST")
	registerRouterFunc("/user/{handle:[a-zA-Z0-9-_]+}/password", passwordEditHandler, "POST")
	registerRouterFunc("/user/{handle:[a-zA-Z0-9-_]+}/promote", promoteHandler, "POST")
}
