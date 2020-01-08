package controllers

import (
	"github.com/kawa-yoiko/botany/app/models"

	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		enc.Encode(c.ShortRepresentation(u))
	}
	w.Write([]byte("]\n"))
}

func parseRequestContest(r *http.Request) (models.Contest, []int64, bool) {
	title := r.PostFormValue("title")
	banner := r.PostFormValue("banner")
	startTime, err1 := strconv.ParseInt(r.PostFormValue("start_time"), 10, 64)
	endTime, err2 := strconv.ParseInt(r.PostFormValue("end_time"), 10, 64)
	desc := r.PostFormValue("desc")
	details := r.PostFormValue("details")
	isVisible := (r.PostFormValue("is_visible") == "true")
	isRegOpen := (r.PostFormValue("is_reg_open") == "true")
	script := r.PostFormValue("script")

	if err1 != nil || err2 != nil || startTime >= endTime {
		return models.Contest{}, nil, false
	}
	// TODO: Check validity of other parameters

	mods := []int64{}
	for _, mod := range strings.Split(r.PostFormValue("moderators"), ",") {
		if strings.TrimSpace(mod) != "" {
			uid, err := strconv.ParseInt(mod, 10, 64)
			if err != nil {
				return models.Contest{}, nil, false
			}
			mods = append(mods, uid)
		}
	}

	c := models.Contest{
		Title:     title,
		Banner:    banner,
		StartTime: startTime,
		EndTime:   endTime,
		Desc:      desc,
		Details:   details,
		IsVisible: isVisible,
		IsRegOpen: isRegOpen,
		Script:    script,
	}
	return c, mods, true
}

// curl http://localhost:3434/contest/create -i -H "Cookie: auth=..." -d
// "title=Grand+Contest&banner=1.png&start_time=0&end_time=1576000000&desc=Really+big+contest&details=Lorem+ipsum+dolor+sit+amet&is_visible=true&is_reg_open=true&script=function+on_submission()%0Aend&moderators=1,2,4"
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

	c, mods, ok := parseRequestContest(r)
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
	if err := c.UpdateModerators(mods); err != nil {
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

func middlewareContestModeratorVerify(w http.ResponseWriter, r *http.Request) (models.User, models.Contest) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		w.WriteHeader(401)
		return models.User{}, models.Contest{Id: -1}
	}
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return models.User{}, models.Contest{Id: -1}
	}
	if c.ParticipationOf(u) != models.ParticipationTypeModerator {
		// No privilege
		w.WriteHeader(403)
		return models.User{}, models.Contest{Id: -1}
	}

	return u, c
}

func contestEditHandler(w http.ResponseWriter, r *http.Request) {
	u, c := middlewareContestModeratorVerify(w, r)
	if c.Id == -1 {
		return
	}

	cNew, mods, ok := parseRequestContest(r)
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
	if err := c.UpdateModerators(mods); err != nil {
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
	if err := p.Update(); err != nil {
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
		Language: r.PostFormValue("lang"),
		Contents: r.PostFormValue("code"),
	}
	if err := s.Create(); err != nil {
		panic(err)
	}
	if err := s.LoadRel(); err != nil {
		panic(err)
	}

	// Send for compilation
	if err := s.SendToQueue(); err != nil {
		panic(err)
	}

	// Invoke contest script
	if err := c.ExecuteScriptOnSubmission(u.Id); err != nil {
		if !errors.Is(err, models.ErrLua) {
			panic(err)
		}
	}

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

func contestSubmissionHistoryHandlerCommon(w http.ResponseWriter, r *http.Request, u models.User, subType int) {
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 || !c.IsVisibleTo(u) {
		w.WriteHeader(404)
		return
	}
	if !c.HasStarted() {
		w.WriteHeader(403)
		fmt.Fprintf(w, "[]")
		return
	}
	if subType == 0 {
		limit, offset := parsePagination(w, r)
		if limit == -1 || offset == -1 {
			w.WriteHeader(400)
			return
		}
		ss, total, err := models.SubmissionHistory(u.Id, c.Id, limit, offset)
		if err != nil {
			panic(err)
		}
		// XXX: Avoid duplication?
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.Encode(map[string]interface{}{
			"total":       total,
			"submissions": ss,
		})
	} else if subType == 1 {
		if c.ParticipationOf(u) == -1 {
			// Querying own submission history, but did not participate
			w.WriteHeader(403)
			fmt.Fprintf(w, "[]")
			return
		}
		ss, _, err := models.SubmissionHistory(u.Id, c.Id, -1, 0)
		if err != nil {
			panic(err)
		}
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.Encode(ss)
	}
}

func contestSubmissionHistoryHandlerUser(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	contestSubmissionHistoryHandlerCommon(w, r, u, 1)
}

func contestSubmissionHistoryHandlerAll(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	contestSubmissionHistoryHandlerCommon(w, r, u, 0)
}

func contestDelegateHandler(w http.ResponseWriter, r *http.Request) {
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
	if !c.HasStarted() {
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	sid, err := strconv.Atoi(r.PostFormValue("submission"))
	if err != nil {
		// Malformed request
		w.WriteHeader(400)
		fmt.Fprintf(w, "{}")
		return
	}

	if sid != -1 {
		s := models.Submission{Id: int32(sid)}
		if err := s.Read(); err != nil {
			if err == sql.ErrNoRows {
				// Treat as permission error
				w.WriteHeader(403)
				fmt.Fprintf(w, "{}")
				return
			} else {
				panic(err)
			}
		}
		if s.User != u.Id || s.Contest != c.Id || s.Status != models.SubmissionStatusAccepted {
			w.WriteHeader(403)
			fmt.Fprintf(w, "{}")
			return
		}
	}

	p := models.ContestParticipation{User: u.Id, Contest: c.Id}
	if err := p.Read(); err != nil {
		if err == sql.ErrNoRows {
			// Did not participate
			w.WriteHeader(403)
			fmt.Fprintf(w, "{}")
			return
		} else {
			panic(err)
		}
	}

	p.Delegate = int32(sid)
	if err := p.Update(); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "{}")
}

func myDelegateHandler(w http.ResponseWriter, r *http.Request) {
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
	if !c.HasStarted() {
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}

	p := models.ContestParticipation{User: u.Id, Contest: c.Id}
	err := p.Read()
	if err == sql.ErrNoRows {
		w.WriteHeader(403)
		fmt.Fprintf(w, "{}")
		return
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]interface{}{
		"submission": p.Delegate,
	})
}

func contestJudgeSelHandler(w http.ResponseWriter, r *http.Request) {
	_, c := middlewareContestModeratorVerify(w, r)
	if c.Id == -1 {
		return
	}

	sid, err := strconv.Atoi(r.PostFormValue("submission"))
	if err != nil {
		// Malformed request
		w.WriteHeader(400)
		return
	}

	if sid != -1 {
		s := models.Submission{Id: int32(sid)}
		if err := s.Read(); err != nil {
			if err == sql.ErrNoRows {
				// Non-existent submission
				w.WriteHeader(400)
				return
			} else {
				panic(err)
			}
		}
		if s.Contest != c.Id || s.Status != models.SubmissionStatusAccepted {
			// Not an accepted submission in current contest
			w.WriteHeader(400)
			return
		}
	}

	c.Judge = int32(sid)
	if err := c.Update(); err != nil {
		panic(err)
	}

	w.WriteHeader(200)
}

func contestRanklistHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}
	limit, offset := parsePagination(w, r)
	if limit == -1 || offset == -1 {
		w.WriteHeader(400)
		return
	}

	ps, total, err := c.PartParticipation(limit, offset)
	if err != nil {
		panic(err)
	}

	prs := []map[string]interface{}{}
	for _, p := range ps {
		prs = append(prs, p.Representation())
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]interface{}{
		"total":        total,
		"participants": prs,
	})
}

func contestMatchesHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}

	limit, offset := parsePagination(w, r)
	if offset == -1 || limit == -1 {
		w.WriteHeader(400)
		return
	}

	ms, err := models.ReadByContest(c.Id)
	if err != nil {
		panic(err)
	}

	total := len(ms)
	msr := []map[string]interface{}{}
	var begin int
	var end int
	if offset > total {
		begin = total
	} else {
		begin = offset
	}
	if offset+limit > total {
		end = total
	} else {
		end = offset + limit
	}
	for _, m := range ms[begin:end] {
		msr = append(msr, m.ShortRepresentation())
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]interface{}{
		"total":   total,
		"matches": msr,
	})
}

// curl http://localhost:3434/contest/1/match/manual -i -H "Cookie: auth=..." -d "submissions=1,2,3"
func contestMatchManualHandler(w http.ResponseWriter, r *http.Request) {
	_, c := middlewareContestModeratorVerify(w, r)
	if c.Id == -1 {
		return
	}

	sids := strings.Split(r.PostFormValue("submissions"), ",")
	m := models.Match{Contest: c.Id, Report: "Pending"}
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

	// Send for match
	if err := m.SendToQueue(c.Judge); err != nil {
		panic(err)
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(m.ShortRepresentation())
}

func contestMatchManualScriptHandler(w http.ResponseWriter, r *http.Request) {
	_, c := middlewareContestModeratorVerify(w, r)
	if c.Id == -1 {
		return
	}

	arg := r.PostFormValue("arg")
	if err := c.ExecuteScriptOnManual(arg); err != nil {
		if !errors.Is(err, models.ErrLua) {
			panic(err)
		}
	}

	w.WriteHeader(200)
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

func contestScriptLogHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["full"]) < 1 {
		w.WriteHeader(400)
		return
	}
	full := r.URL.Query()["full"][0]
	if full != "0" && full != "1" {
		w.WriteHeader(400)
		return
	}

	_, c := middlewareContestModeratorVerify(w, r)
	if c.Id == -1 {
		return
	}

	var err error
	var s string

	if full == "1" {
		// Full log
		err, s = c.ReadScriptLog()
		if err != nil {
			panic(err)
		}
		// Write file name
		w.Header().Set("Content-Disposition",
			"attachment; filename=\"contest_log_"+strconv.FormatInt(int64(c.Id), 10)+
				"_"+time.Now().Format("20060102150405")+".txt\"")
	} else {
		// Tail log
		s = c.TailLog()
	}

	w.WriteHeader(200)
	w.Write([]byte(s))
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
	registerRouterFunc("/contest/{cid:[0-9]+}/my", contestSubmissionHistoryHandlerUser, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/submission/list", contestSubmissionHistoryHandlerAll, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/delegate", contestDelegateHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/my_delegate", myDelegateHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/judge", contestJudgeSelHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/ranklist", contestRanklistHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/matches", contestMatchesHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/match/manual", contestMatchManualHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/match/manual_script", contestMatchManualScriptHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/match/{mid:[0-9]+}", contestMatchDetailsHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/script_log", contestScriptLogHandler, "GET")
}
