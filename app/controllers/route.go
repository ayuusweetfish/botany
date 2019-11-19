package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

var Router *mux.Router

func registerRouterFunc(path string, fn func(http.ResponseWriter, *http.Request)) {
	if Router == nil {
		Router = mux.NewRouter()
	}
	Router.HandleFunc(path, fn)
}
