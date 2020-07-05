package controllers

import (
	"github.com/kawa-yoiko/botany/app/globals"
	"github.com/kawa-yoiko/botany/app/models"
	"mime"

	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
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

// Call r.ParseMultipartForm() beforehand
func maybePostFormValue(r *http.Request, key string) (string, bool) {
	if vs := r.PostForm[key]; len(vs) > 0 {
		return vs[0], true
	}
	return "", false
}

func parseRequestContest(r *http.Request, c models.Contest) (models.Contest, []int64, bool) {
	c.LoadPlayback()

	// 32 MB, in accordance with standard library's implementation
	// https://golang.org/src/net/http/request.go - PostFormValue
	r.ParseMultipartForm(32 << 20)

	if s, has := maybePostFormValue(r, "title"); has {
		c.Title = s
	}
	if s, has := maybePostFormValue(r, "start_time"); has {
		if startTime, err := strconv.ParseInt(s, 10, 64); err == nil {
			c.StartTime = startTime
		}
	}
	if s, has := maybePostFormValue(r, "end_time"); has {
		if endTime, err := strconv.ParseInt(s, 10, 64); err == nil {
			c.EndTime = endTime
		}
	}
	if s, has := maybePostFormValue(r, "desc"); has {
		c.Desc = s
	}
	if s, has := maybePostFormValue(r, "details"); has {
		c.Details = s
	}
	if s, has := maybePostFormValue(r, "is_visible"); has {
		c.IsVisible = (s == "true")
	}
	if s, has := maybePostFormValue(r, "is_reg_open"); has {
		c.IsRegOpen = (s == "true")
	}
	if s, has := maybePostFormValue(r, "script"); has {
		c.Script = s
	}
	if s, has := maybePostFormValue(r, "playback"); has {
		c.Playback = s
	}
	if c.StartTime >= c.EndTime {
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

	c, mods, ok := parseRequestContest(r, models.Contest{})
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
	_, c := middlewareContestModeratorVerify(w, r)
	if c.Id == -1 {
		return
	}

	cNew, mods, ok := parseRequestContest(r, c)
	if !ok {
		// Malformed request
		w.WriteHeader(400)
		fmt.Fprintf(w, "{}")
		return
	}

	if err := cNew.Update(); err != nil {
		panic(err)
	}
	if err := cNew.UpdatePlayback(); err != nil {
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
	if err := c.ExecuteMatchScriptOnSubmission(u.Id); err != nil {
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
	participation := c.ParticipationOf(u)
	if !c.HasStarted() && participation != models.ParticipationTypeModerator {
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
		ss := []map[string]interface{}{}
		if participation == -1 {
			// Querying own submission history, but did not participate
			w.WriteHeader(403)
			fmt.Fprintf(w, "[]")
			return
		} else if participation == models.ParticipationTypeContestant {
			s, _, err := models.SubmissionHistory(u.Id, c.Id, -1, 0)
			if err != nil {
				panic(err)
			}
			ss = append(ss, s...)
		} else {
			mods, err := c.ReadModerators()
			if err != nil {
				panic(err)
			}
			sus, err := models.AllSuperusers()
			if err != nil {
				panic(err)
			}
			mods = append(mods, sus...)
			for i := range mods {
				s, _, err := models.SubmissionHistory(mods[i], c.Id, -1, 0)
				if err != nil {
					panic(err)
				}
				ss = append(ss, s...)
			}
		}
		sort.Slice(ss, func(i, j int) bool {
			return ss[i]["id"].(int32) > ss[j]["id"].(int32)
		})
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

func contestJudgeIdHandler(w http.ResponseWriter, r *http.Request) {
	_, c := middlewareContestModeratorVerify(w, r)
	if c.Id == -1 {
		return
	}
	if err := c.Read(); err != nil {
		panic(err)
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]interface{}{
		"judge": c.Judge,
	})
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

	ms, total, err := c.Matches(limit, offset)
	if err != nil {
		panic(err)
	}

	msr := []map[string]interface{}{}
	for _, m := range ms {
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
	if err := c.ExecuteMatchScriptOnManual(arg); err != nil {
		if !errors.Is(err, models.ErrLua) {
			panic(err)
		}
	}

	w.WriteHeader(200)
}

func middlewareMatchDetails(w http.ResponseWriter, r *http.Request) (models.Contest, models.Match) {
	u := middlewareAuthRetrieve(w, r)
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return models.Contest{Id: -1}, models.Match{Id: -1}
	}

	mid, _ := strconv.Atoi(mux.Vars(r)["mid"])
	m := models.Match{Id: int32(mid)}

	if mid != 0 {
		if err := m.Read(); err != nil {
			if err == sql.ErrNoRows {
				// No such match
				w.WriteHeader(404)
				return models.Contest{Id: -1}, models.Match{Id: -1}
			} else {
				panic(err)
			}
		}
		if m.Contest != c.Id {
			// Match not in contest
			w.WriteHeader(404)
			return models.Contest{Id: -1}, models.Match{Id: -1}
		}
	}

	return c, m
}

func contestMatchDetailsHandler(w http.ResponseWriter, r *http.Request) {
	_, m := middlewareMatchDetails(w, r)
	if m.Id == -1 || m.Id == 0 {
		return
	}

	m.LoadRel()
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(m.Representation())
}

func contestMatchLogHandler(w http.ResponseWriter, r *http.Request) {
	_, m := middlewareMatchDetails(w, r)
	if m.Id == -1 || m.Id == 0 {
		return
	}

	party, _ := strconv.Atoi(mux.Vars(r)["party"])
	p := models.MatchParty{Match: m.Id, Index: int32(party)}
	if err := p.LoadLog(); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		} else {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(p.Log))
}

func contestMatchPlaybackHandler(w http.ResponseWriter, r *http.Request) {
	c, m := middlewareMatchDetails(w, r)
	if m.Id == -1 {
		return
	}

	c.LoadPlayback()
	s := c.Playback
	if m.Id != 0 {
		s = strings.Replace(s, "<% report %>", m.Report, -1)
		if strings.Contains(s, "<% report js str %>") {
			jsStr, err := json.Marshal(m.Report)
			if err != nil {
				panic(err)
			}
			s = strings.Replace(s, "<% report js str %>", string(jsStr), -1)
		}
		if strings.Contains(s, "<% num parties %>") {
			partiesCount, err := m.PartiesCount()
			if err != nil {
				panic(err)
			}
			s = strings.Replace(s, "<% num parties %>", strconv.Itoa(partiesCount), -1)
		}
		r := regexp.MustCompile("<% party [0-9]{1,3} %>")
		loaded := false
		s = r.ReplaceAllStringFunc(s, func(t string) string {
			if !loaded {
				if err := m.LoadRel(); err != nil {
					panic(err)
				}
				loaded = true
			}
			t = strings.TrimPrefix(t, "<% party ")
			t = strings.TrimSuffix(t, " %>")
			n, _ := strconv.Atoi(t)
			if n < 0 || n >= len(m.Rel.Parties) {
				return ""
			} else {
				return m.Rel.Parties[n].Rel.User.Nickname
			}
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(s))
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

func contestBannerHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		w.WriteHeader(404)
		return
	}
	f, err := c.LoadBanner()
	if err != nil {
		panic(err)
	}

	if f.Id == -1 {
		f = globals.DefaultBanner
	}

	w.Header().Set("Content-Type", f.Type)
	w.Write(f.Content)
}

func contestBannerUploadHandler(w http.ResponseWriter, r *http.Request) {
	u := middlewareAuthRetrieve(w, r)
	if u.Id == -1 {
		return
	}
	c := middlewareReferredContest(w, r, u)
	if c.Id == -1 {
		return
	}
	if c.Owner != u.Id && u.Privilege != models.UserPrivilegeSuperuser {
		w.WriteHeader(403)
		return
	}
	err, ext, buf := middlewareMultipartFormFile(r)
	if errors.Is(err, http.ErrMissingFile) {
		w.WriteHeader(400)
		return
	} else if err != nil {
		panic(err)
	}

	mimeType := mime.TypeByExtension("." + ext)
	if !strings.HasPrefix(mimeType, "image/") {
		w.WriteHeader(400)
		return
	}
	f, err := c.LoadBanner()
	if err != nil {
		panic(err)
	}
	f.Type = mimeType
	f.Content = buf.Bytes()

	if f.Id == -1 {
		err = f.Create()
		if err == nil {
			c.Banner = f.Id
			err = c.UpdateBanner()
		}
	} else {
		err = f.Update()
	}
	if err != nil {
		panic(err)
	}
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
	registerRouterFunc("/contest/{cid:[0-9]+}/match/{mid:[0-9]+}/log/{party:[0-9]+}", contestMatchLogHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/match/{mid:[0-9]+}/playback", contestMatchPlaybackHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/script_log", contestScriptLogHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/banner", contestBannerHandler, "GET")
	registerRouterFunc("/contest/{cid:[0-9]+}/banner/upload", contestBannerUploadHandler, "POST")
	registerRouterFunc("/contest/{cid:[0-9]+}/judge_id", contestJudgeIdHandler, "GET")
}
