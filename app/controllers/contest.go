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
	u := middlewareAuthRetrieve(w, r)

	cs, err := models.ContestReadAll()
	if err != nil {
		panic(err)
	}

	w.Write([]byte("["))
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	first := true
	for _, c := range cs {
		// Skip invisible contests
		if !c.IsVisibleTo(u) {
			continue
		}

		if first {
			first = false
		} else {
			w.Write([]byte(","))
		}
		enc.Encode(c.ShortRepresentation())
	}
	w.Write([]byte("]\n"))
}

// Retrieves the contest referred to in the URL parameter
// Returns the object without relationships loaded; or
// an empty one with an Id of -1 if none is found
func middlewareReferredContest(w http.ResponseWriter, r *http.Request, u models.User) models.Contest {
	cid, _ := strconv.Atoi(mux.Vars(r)["cid"])
	c := models.Contest{Id: int32(cid)}
	if err := c.Read(); err != nil {
		if err == sql.ErrNoRows {
			return models.Contest{Id: -1}
		} else {
			panic(err)
		}
	}
	if c.IsVisibleTo(u) {
		return c
	} else {
		return models.Contest{Id: -1}
	}
}

func contestInfoHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}

	c.LoadRel()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(c.Representation(u))
}

// curl http://localhost:3434/contest/1/publish -i -H "Cookie: auth=..." -d "set=true"
func contestPublishHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}

	if u.Privilege != models.UserPrivilegeSuperuser {
		// No privilege
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		// No such contest
		w.WriteHeader(404)
		return
	}

	c.IsVisible = (r.PostFormValue("set") == "true")
	if err := c.Update(); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "{}")
}

// curl http://localhost:3434/contest/1/join -i -H "Cookie: auth=..." -d ""
func contestJoinHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}

	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 || !c.IsVisibleTo(u) {
		w.WriteHeader(404)
		return
	}
	if !c.IsRegOpen {
		w.WriteHeader(403)
		// Registration not open
		fmt.Fprintf(w, "{}")
		return
	}

	p := models.ContestParticipation{
		User:    u.Id,
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
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}

	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 || !c.IsVisibleTo(u) {
		// Nonexistent or invisible contest
		w.WriteHeader(404)
		return
	}

	participation := c.ParticipationOf(u)
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
		User:     u.Id,
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
	u := middlewareAuthRetrieve(w, r)

	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 || !c.IsVisibleTo(u) {
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
	if !s.IsVisibleTo(u) {
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(s.Representation())
}

func contestSubmissionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}

	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 || !c.IsVisibleTo(u) {
		w.WriteHeader(404)
		return
	}
	if c.ParticipationOf(u) == -1 {
		// Did not participate
		w.WriteHeader(403)
		fmt.Fprintf(w, "[]")
		return
	}

	ss, err := models.SubmissionHistory(u.Id, c.Id, 5)
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

func contestRanklistHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}

	ps, err := c.AllParticipations()
	if err != nil {
		panic(err)
	}

	w.Write([]byte("["))
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for i, p := range ps {
		if i != 0 {
			w.Write([]byte(","))
		}
		enc.Encode(p.Representation())
	}
	w.Write([]byte("]\n"))
}

// XXX: For debug use
// curl http://localhost:3434/contest/create -i -H "Cookie: auth=..." -d ""
func contestCreateHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}

	c := models.Contest{
		Title:     "Grand Contest",
		Banner:    "",
		Owner:     u.Id,
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

func contestMatchesHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}

	matches, err := models.ReadByContest(c.Id)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	leftparam := `[`
	rightparam := `]`
	dot := `,`
	enc.Encode(leftparam)
	for _, item := range matches {
		matchMap := item.ShortRepresentation()
		enc.Encode(matchMap)
		enc.Encode(dot)
	}
	enc.Encode(rightparam)
}

func init() {
	registerRouterFunc("/contest/list", contestListHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/publish", contestPublishHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/info", contestInfoHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/join", contestJoinHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/submit", contestSubmitHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/submission/{sid:[0-9]+}", contestSubmissionHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/my", contestSubmissionHistoryHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/ranklist", contestRanklistHandler, "GET")
	registerRouterFunc("/contest/create", contestCreateHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+/matches}", contestMatchesHandler, "GET")
}
