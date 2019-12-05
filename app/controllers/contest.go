package controllers

import (
	"github.com/kawa-yoiko/botany/app/models"

	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func contestListHandler(w http.ResponseWriter, r *http.Request) {
	cs, err := models.ContestReadAll()
	if err != nil {
		panic(err)
	}

	w.Write([]byte("["))
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for i, c := range cs {
		if i != 0 {
			w.Write([]byte(","))
		}
		enc.Encode(c.ShortRepresentation())
	}
	w.Write([]byte("]\n"))
}

// Retrieves the contest referred to in the URL parameter
// Returns the object without relationships loaded; or
// an empty one with an Id of -1 if none is found
func middlewareReferredContest(w http.ResponseWriter, r *http.Request) models.Contest {
	cid, _ := strconv.Atoi(mux.Vars(r)["cid"])
	c := models.Contest{Id: int32(cid)}
	if err := c.Read(); err != nil {
		if err == sql.ErrNoRows {
			return models.Contest{Id: -1}
		} else {
			panic(err)
		}
	}
	return c
}

func contestInfoHandler(w http.ResponseWriter, r *http.Request) {
	c := middlewareReferredContest(w, r)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}

	c.LoadRel()

	uid := middlewareAuthRetrieve(w, r)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(c.Representation(uid))
}

// curl http://localhost:3434/contest/1/join -i -H "Cookie: auth=..." -d ""
func contestJoinHandler(w http.ResponseWriter, r *http.Request) {
	uid := middlewareAuthRetrieve(w, r)
	if uid == -1 {
		w.WriteHeader(401)
		return
	}

	c := middlewareReferredContest(w, r)
	if c.Id == -1 || !c.IsVisibleTo(uid) {
		w.WriteHeader(404)
		return
	}
	if !c.IsRegOpen {
		w.WriteHeader(400)
		// Registration not open
		fmt.Fprintf(w, "{}")
		return
	}

	p := models.ContestParticipation{
		User:    uid,
		Contest: c.Id,
		Type:    models.ParticipationTypeContestant,
	}
	if err := p.Create(); err != nil {
		panic(err)
	}

	// Success
	fmt.Fprintf(w, "{}")
}

// curl http://localhost:3434/contest/1/submit -i -H "Cookie: auth=..." -d "code=123%20456"
func contestSubmitHandler(w http.ResponseWriter, r *http.Request) {
	uid := middlewareAuthRetrieve(w, r)
	if uid == -1 {
		w.WriteHeader(401)
		return
	}

	c := middlewareReferredContest(w, r)
	if c.Id == -1 || !c.IsVisibleTo(uid) {
		// Nonexistent or invisible contest
		w.WriteHeader(404)
		return
	}

	participation := c.ParticipationOf(uid)
	if participation == -1 ||
		(participation != models.ParticipationTypeModerator && !c.IsRunning()) {
		// Did not participate
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	// TODO: Check submission length and character set

	// Create a new submission
	s := models.Submission{
		User:     uid,
		Contest:  c.Id,
		Contents: r.PostFormValue("code"),
	}
	if err := s.Create(); err != nil {
		panic(err)
	}
	s.LoadRel()

	// Success
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]interface{}{
		"err":        0,
		"submission": s.ShortRepresentation(),
	})
}

func contestSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	uid := middlewareAuthRetrieve(w, r)

	c := middlewareReferredContest(w, r)
	if c.Id == -1 || !c.IsVisibleTo(uid) {
		w.WriteHeader(404)
		return
	}

	sid, _ := strconv.Atoi(mux.Vars(r)["sid"])
	s := models.Submission{Id: int32(sid)}
	if err := s.Read(); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		} else {
			panic(err)
		}
	}

	s.LoadRel()

	// Disallow viewing of others' code during a contest for non-moderators
	if !s.IsVisibleTo(uid) {
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(s.Representation())
}

func contestSubmissionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	uid := middlewareAuthRetrieve(w, r)
	if uid == -1 {
		w.WriteHeader(401)
		return
	}

	c := middlewareReferredContest(w, r)
	if c.Id == -1 || !c.IsVisibleTo(uid) {
		w.WriteHeader(404)
		return
	}

	ss, err := models.SubmissionHistory(uid, c.Id, 5)
	if err != nil {
		panic(err)
	}

	// XXX: Avoid duplication?
	w.Write([]byte("["))
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for i, s := range ss {
		if i != 0 {
			w.Write([]byte(","))
		}
		enc.Encode(s.ShortRepresentation())
	}
	w.Write([]byte("]\n"))
}

// XXX: For debug use
// curl http://localhost:3434/contest/create -i -H "Cookie: auth=..." -d ""
func contestCreateHandler(w http.ResponseWriter, r *http.Request) {
	uid := middlewareAuthRetrieve(w, r)
	if uid == -1 {
		w.WriteHeader(401)
		return
	}

	c := models.Contest{
		Title:     "Grand Contest",
		Banner:    "",
		Owner:     uid,
		StartTime: 0,
		EndTime:   365 * 86400,
		Desc:      "Really big contest",
		Details:   "Lorem ipsum dolor sit amet",
		IsVisible: true,
		IsRegOpen: true,
	}
	if err := c.Create(); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "{\"id\": %d}", c.Id)
}

func init() {
	registerRouterFunc("/contest/list", contestListHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/info", contestInfoHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/join", contestJoinHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/submit", contestSubmitHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/submission/{sid:[0-9]+}", contestSubmissionHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/my", contestSubmissionHistoryHandler, "GET")
	registerRouterFunc("/contest/create", contestCreateHandler, "POST")
}
