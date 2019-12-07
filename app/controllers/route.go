package controllers

import (
	"github.com/kawa-yoiko/botany/app/globals"

	"net/http"

	"github.com/gorilla/mux"
)

var router *mux.Router
var apiRouter *mux.Router

func registerRouterFunc(path string, fn func(http.ResponseWriter, *http.Request), methods ...string) {
	if router == nil {
		router = mux.NewRouter()
		apiRouter = router.PathPrefix(globals.Config().ApiPrefix).Subrouter()
	}
	if len(methods) == 0 {
		methods = []string{"GET"}
	}
	apiRouter.HandleFunc(path, fn).Methods(methods...)
}

func GetRootRouterFunc() func(w http.ResponseWriter, req *http.Request) {
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
