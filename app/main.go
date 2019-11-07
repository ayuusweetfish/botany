package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const HTTPListenPort = 3434

var db *sql.DB

var schema = `
CREATE TABLE visitor (
	name TEXT,
	count INTEGER
)`

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page!")
}

func nameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var count int
	row := db.QueryRow("SELECT count FROM visitor WHERE name = $1", name)
	err := row.Scan(&count)
	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO visitor(name, count) VALUES ($1, 1)", name)
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
		_, err = db.Exec("UPDATE visitor SET count = $1 WHERE name = $2", count, name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	fmt.Fprintf(w, "Hi %s, your #%d visit!", name, count)
}

func main() {
	// $ initdb -D ./data
	// $ pg_ctl -D ./data start
	// $ createdb dbqwq -U uwu
	var err error
	db, err = sql.Open("postgres", "sslmode=disable dbname=dbqwq user=uwu")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	db.Exec(schema)

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/{name:[a-z]+}", nameHandler)
	http.Handle("/", r)

	log.Printf("Listening on http://localhost:%d/\n", HTTPListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", HTTPListenPort), nil))
}
