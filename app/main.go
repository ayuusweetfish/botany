package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const HTTPListenPort = 3434

var db *sql.DB
var store = sessions.NewCookieStore([]byte("vertraulich"))
var uniqId = 0

// A middleware-like function which retrieves the unique ID from cookies,
// assigning a new one if none exist, and write it to the response.
func middlewareProcessSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "QAQ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, ok := session.Values["uniq-id"].(int)
	if !ok || session.IsNew {
		uniqId++
		id = uniqId
		session.Values["uniq-id"] = id
		// Save session
		// Note that the cookie store writes to the HTTP header
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	// All headers go before actual content
	fmt.Fprintf(w, "Your unique ID is %d\n", id)
}

// Routed to /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)
	fmt.Fprintf(w, "This is the home page!")
}

// Routed to /{name:[a-z]+}
func nameHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)

	vars := mux.Vars(r)
	name := vars["name"]
	var count int
	row := db.QueryRow("SELECT count FROM b_user WHERE username = $1", name)
	err := row.Scan(&count)
	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO b_user(username, count) VALUES ($1, 1)", name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		count = 1
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else {
		count += 1
		_, err = db.Exec("UPDATE b_user SET count = $1 WHERE username = $2", count, name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	fmt.Fprintf(w, "Hi %s, your #%d visit!", name, count)
}

const (
	dbuser     = "list"
	dbpassword = ""
	dbname     = "botanyDatabase"
)

func dbConnect() *sql.DB {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("sslmode=disable dbname=%s user=%s", dbname, dbuser))
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func main() {
	//$ initdb -D ./data
	//$ pg_ctl -D ./data start
	//$ createdb dbqwq -U uwu
	db = dbConnect()
	defer db.Close()
	db.Exec(schema)
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/{name:[a-z]+}", nameHandler)
	http.Handle("/", r)
	log.Printf("Listening on http://localhost:%d/\n", HTTPListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", HTTPListenPort), nil))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
