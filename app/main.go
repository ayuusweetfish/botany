package main

import (
	"github.com/kawa-yoiko/botany/app/controllers"
	"github.com/kawa-yoiko/botany/app/globals"
	"github.com/kawa-yoiko/botany/app/models"

	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func main() {
	config := globals.Config()
	port := config.AppPort

	db, err := sql.Open("postgres",
		"sslmode=disable dbname="+config.DbName+
			" user="+config.DbUser+
			" password="+config.DbPassword)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if config.RedisPort != 0 {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%d", config.RedisPort),
			Password: config.RedisPassword,
			DB:       0,
		})
		_, err = redisClient.Ping().Result()
		if err != nil {
			log.Fatalln(err)
		}
		models.InitializeRedis(redisClient)
	}

	keyPairs := [][]byte{}
	for i, s := range config.CookieKeyPairs {
		if i%2 == 1 && len(s) != 16 && len(s) != 24 && len(s) != 32 {
			log.Fatalf("Size of key should be 16, 24 or 32 bytes ('%s' is %d bytes)", s, len(s))
		}
		keyPairs = append(keyPairs, []byte(s))
	}
	globals.SessionStore = sessions.NewCookieStore(keyPairs...)

	models.InitializeSchemata(db)
	http.HandleFunc("/", controllers.GetRootRouterFunc())

	log.Printf("Listening on http://localhost:%d%s\n", port, config.ApiPrefix)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
