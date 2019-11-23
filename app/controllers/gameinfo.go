package controllers

import (
	"fmt"
	"net/http"
)

// Routed to /gamelist
func gamelistHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)
	fmt.Fprintf(w, `{"total": `))
}

func init () {
	registerFunc("/gamelist", gamelistHandler)
}
