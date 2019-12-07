package controllers

import (
	"github.com/kawa-yoiko/botany/app/models"

	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func parseRequestContest(r *http.Request) (models.Contest, bool) {
	title := r.PostFormValue("title")
	banner := r.PostFormValue("banner")
	startTime, err1 := strconv.ParseInt(r.PostFormValue("start_time"), 10, 64)
	endTime, err2 := strconv.ParseInt(r.PostFormValue("end_time"), 10, 64)
	desc := r.PostFormValue("desc")
	details := r.PostFormValue("details")
	isVisible := (r.PostFormValue("is_visible") == "true")
	isRegOpen := (r.PostFormValue("is_reg_open") == "true")

	if err1 != nil || err2 != nil || startTime >= endTime {
		return models.Contest{}, false
	}
	// TODO: Check validity of other parameters

	c := models.Contest{
		Title:     title,
		Banner:    banner,
		StartTime: startTime,
		EndTime:   endTime,
		Desc:      desc,
		Details:   details,
		IsVisible: isVisible,
		IsRegOpen: isRegOpen,
	}
	return c, true
}

// curl http://localhost:3434/contest/create -i -H "Cookie: auth=..." -d
// "title=Grand+Contest&banner=1.png&start_time=0&end_time=1576000000&desc=Really+big+contest&details=Lorem+ipsum+dolor+sit+amet&is_visible=true&is_reg_open=true"
func contestCreateHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}
	if u.Privilege < models.UserPrivilegeOrganizer {
		// Insufficient privilege
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	c, ok := parseRequestContest(r)
	if !ok {
		// Malformed request
		w.WriteHeader(400)
		fmt.Fprintf(w, "{}")
		return
	}

	c.Owner = u.Id
	if err := c.Create(); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "{\"id\": %d}", c.Id)
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

func contestEditHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}

	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}
	if u.Id != c.Owner {
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	cNew, ok := parseRequestContest(r)
	if !ok {
		// Malformed request
		w.WriteHeader(400)
		fmt.Fprintf(w, "{}")
		return
	}

	cNew.Id = c.Id
	cNew.Owner = u.Id
	if err := cNew.Update(); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "{}")
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

func contestMatchesHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}

	ms, err := models.ReadByContest(c.Id)
	if err != nil {
		panic(err)
	}

	w.Write([]byte("["))
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for i, m := range ms {
		if i != 0 {
			w.Write([]byte(","))
		}
		enc.Encode(m.ShortRepresentation())
	}
	w.Write([]byte("]\n"))
}

// curl http://localhost:3434/contest/1/match/manual -i -H "Cookie: auth=..." -d "submissions=1,2,3"
func contestMatchManualHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return
	}
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}
	if c.ParticipationOf(u) != models.ParticipationTypeModerator {
		// No privilege
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	sids := strings.Split(r.PostFormValue("submissions"), ",")
	m := models.Match{Contest: c.Id, Report: "{\"winner\": \"In queue\"}"}
	for _, sid := range sids {
		sidN, err := strconv.Atoi(sid)
		if err != nil {
			// Malformed request
			w.WriteHeader(400)
			fmt.Fprintf(w, "{}")
			return
		}
		m.Rel.Parties = append(m.Rel.Parties,
			models.Submission{Id: int32(sidN)})
	}

	if err := m.Create(); err != nil {
		panic(err)
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(m.ShortRepresentation())
}

func contestMatchDetailsHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}

	mid, _ := strconv.Atoi(mux.Vars(r)["mid"])
	m := models.Match{Id: int32(mid)}
	if err := m.Read(); err != nil {
		if err == sql.ErrNoRows {
			// No such match
			w.WriteHeader(404)
			return
		} else {
			panic(err)
		}
	}
	if m.Contest != c.Id {
		// Match not in contest
		w.WriteHeader(404)
		return
	}
	m.LoadRel()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(m.Representation())
}

func init() {
	registerRouterFunc("/contest/list", contestListHandler, "GET")
	registerRouterFunc("/contest/create", contestCreateHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/edit", contestEditHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/publish", contestPublishHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/info", contestInfoHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/join", contestJoinHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/submit", contestSubmitHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/submission/{sid:[0-9]+}", contestSubmissionHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/my", contestSubmissionHistoryHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/ranklist", contestRanklistHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/matches", contestMatchesHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/match/manual", contestMatchManualHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/match/{mid:[0-9]+}", contestMatchDetailsHandler, "GET")
}
