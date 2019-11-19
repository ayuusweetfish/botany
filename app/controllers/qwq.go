package controllers

import (
	"../globals"

	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var uniqId = 0

// A middleware-like function which retrieves the unique ID from cookies,
// assigning a new one if none exist, and write it to the response.
func middlewareProcessSession(w http.ResponseWriter, r *http.Request) {
	session, err := globals.SessionStore.Get(r, "QAQ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, ok := session.Values["uniq-id"].(int)
	if !ok || session.IsNew {
		uniqId++
		id = uniqId + 1
		session.Values["uniq-id"] = id
		// Save session
		// Note that the cookie store writes to the HTTP header
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	// All headers go before actual content
	//fmt.Fprintf(w, "Your unique ID is %d\n", id)
}

// Routed to /{name:[a-z]+}
func nameHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)

	vars := mux.Vars(r)
	name := vars["name"]
	var count int
	row := globals.DB.QueryRow("SELECT count FROM b_user WHERE username = $1", name)
	err := row.Scan(&count)
	if err == sql.ErrNoRows {
		_, err = globals.DB.Exec("INSERT INTO b_user(username, count) VALUES ($1, 1)", name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		count = 1
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else {
		count += 1
		_, err = globals.DB.Exec("UPDATE b_user SET count = $1 WHERE username = $2", count, name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	fmt.Fprintf(w, "Hi %s, your #%d visit!", name, count)
}
