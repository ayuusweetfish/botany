package main

import (
	"github.com/kawa-yoiko/botany/app/controllers"
	"github.com/kawa-yoiko/botany/app/globals"
	"github.com/kawa-yoiko/botany/app/models"

	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

type Config struct {
	AppPort    int    `json:"app_port"`
	DbName     string `json:"db_name"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
}

func LoadConfiguration(file string) Config {
	var c Config
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(f)
	decoder.Decode(&c)
	return c
}

func main() {
	config := LoadConfiguration("config.json")
	port := config.AppPort

	db, err := sql.Open("postgres",
		"sslmode=disable dbname="+config.DbName+
			" user="+config.DbUser+
			" password="+config.DbPassword)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	globals.SessionStore = sessions.NewCookieStore([]byte("vertraulich"))

	models.InitializeSchemata(db)
	http.HandleFunc("/", controllers.GetGlobalRouterFunc())

	log.Printf("Listening on http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
