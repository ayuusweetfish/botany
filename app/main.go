package main

import (
	"github.com/kawa-yoiko/botany/app/controllers"
	"github.com/kawa-yoiko/botany/app/globals"
	"github.com/kawa-yoiko/botany/app/models"

	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

const HTTPListenPort = 3434

func main() {
	// $ initdb -D ./data
	// $ pg_ctl -D ./data start
	// $ createdb dbqwq -U uwu
	db, err := sql.Open("postgres", "sslmode=disable dbname=dbqwq user=uwu")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	globals.SessionStore = sessions.NewCookieStore([]byte("vertraulich"))

	models.InitializeSchemata(db)
	http.HandleFunc("/", controllers.GetGlobalRouterFunc())

	log.Printf("Listening on http://localhost:%d/\n", HTTPListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", HTTPListenPort), nil))
}
