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

const HTTPListenPort = 3000

const (
	dbuser     = "sakura"
	dbpassword = "123456"
	dbname     = "botanydatabase"
)

func dbConnect() *sql.DB {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("sslmode=disable dbname=%s "+
		"user=%s password=%s", dbname, dbuser, dbpassword))
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func main() {
	//$ initdb -D ./data
	//$ pg_ctl -D ./data start
	//$ createdb dbqwq -U uwu
	db, err := sql.Open("postgres", fmt.Sprintf("sslmode=disable dbname=%s user=%s", dbname, dbuser))
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
