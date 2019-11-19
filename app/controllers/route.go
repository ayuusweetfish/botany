package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

var router *mux.Router

func registerRouterFunc(path string, fn func(http.ResponseWriter, *http.Request)) {
	if router == nil {
		router = mux.NewRouter()
	}
	router.HandleFunc(path, fn)
}

func GetGlobalRouterFunc() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				if e, ok := e.(error); ok {
					http.Error(w, e.Error(), 500)
				}
			}
		}()
		router.ServeHTTP(w, req)
	}
}
