package controllers

import (
	"../globals"
	"../models"

	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

	u := models.QwQUser{Name: name}
	err := u.Read()

	if err == sql.ErrNoRows {
		u.Count = 1
		err = u.Create()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else {
		u.Count += 1
		err = u.Update()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	fmt.Fprintf(w, "Hi %s, your #%d visit!", u.Name, u.Count)
}

func init() {
	registerRouterFunc("/{name:[a-z]+}", nameHandler)
}
