package controllers

import "github.com/gorilla/mux"

func Handler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/{name:[a-z]+}", nameHandler)
	return r
}

