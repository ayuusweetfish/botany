package controllers

import (
	"github.com/kawa-yoiko/botany/app/models"

	"fmt"
	"net/http"
)

// Routed to /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page!")
}

func fakeDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	models.FakeDatabase()
	fmt.Fprintf(w, "Fake News prepared\n")
}

func fakeMatchesHandler(w http.ResponseWriter, r *http.Request) {
	models.FakeMatches()
	fmt.Fprintf(w, "Matches automagically run\n")
}

func init() {
	registerRouterFunc("/", rootHandler)
	registerRouterFunc("/fake", fakeDatabaseHandler, "POST")
	registerRouterFunc("/fake_matches", fakeMatchesHandler, "POST")
}
