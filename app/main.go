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

	CookieKeyPairs []string `json:"cookie_key_pairs"`
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

	keyPairs := [][]byte{}
	for i, s := range config.CookieKeyPairs {
		if i%2 == 1 && len(s) != 16 && len(s) != 24 && len(s) != 32 {
			log.Fatalf("Size of key should be 16, 24 or 32 bytes ('%s' is %d bytes)", s, len(s))
		}
		keyPairs = append(keyPairs, []byte(s))
	}
	globals.SessionStore = sessions.NewCookieStore(keyPairs...)

	models.InitializeSchemata(db)
	http.HandleFunc("/", controllers.GetGlobalRouterFunc())

	log.Printf("Listening on http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
