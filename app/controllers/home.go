package controllers

import (
	"fmt"
	"net/http"
)

// Routed to /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)
	fmt.Fprintf(w, "This is the home page!")
}

func init() {
	registerRouterFunc("/", rootHandler)
}
