package globals

import (
	"encoding/json"
	"log"
	"os"
)

type GlobalConfig struct {
	AppPort       int    `json:"app_port"`
	DbName        string `json:"db_name"`
	DbUser        string `json:"db_user"`
	DbPassword    string `json:"db_password"`
	RedisPort     int    `json:"redis_port"`
	RedisPassword string `json:"redis_password"`

	CookieKeyPairs []string `json:"cookie_key_pairs"`

	ApiPrefix string `json:"api_prefix"`
}

var configInitialized = false
var config GlobalConfig

func LoadConfiguration(file string) GlobalConfig {
	var c GlobalConfig
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(f)
	decoder.Decode(&c)
	return c
}

func Config() *GlobalConfig {
	if !configInitialized {
		config = LoadConfiguration("config.json")
	}
	return &config
}
