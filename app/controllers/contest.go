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
	w.Write([]byte("]"))
}

func contestInfoHandler(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.Atoi(mux.Vars(r)["cid"])

	c := models.Contest{Id: int32(cid)}
	if err := c.Read(); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		} else {
			panic(err)
		}
	}
	c.LoadRel()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(c.Representation())
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
	registerRouterFunc("/contest/create", contestCreateHandler, "POST")
}
