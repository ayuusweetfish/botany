package main

import (
	"./controllers"
	"./globals"

	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
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

func main() {
	// $ initdb -D ./data
	// $ pg_ctl -D ./data start
	// $ createdb dbqwq -U uwu
	db, err := sql.Open("postgres", "sslmode=disable dbname=dbqwq user=uwu")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	db.Exec(schema)

	globals.DB = db
	globals.SessionStore = sessions.NewCookieStore([]byte("vertraulich"))

	http.Handle("/", controllers.Handler())

	log.Printf("Listening on http://localhost:%d/\n", HTTPListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", HTTPListenPort), nil))
}
