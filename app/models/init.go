package models

import (
	"database/sql"
	"strings"
)

var db *sql.DB
var schemata []string

func registerSchema(sql string) {
	schemata = append(schemata, sql)
}

func InitializeSchemata(dbInput *sql.DB) {
	db = dbInput
	s := strings.Join(schemata, "\n")
	db.Exec(schema)
	println(s)
}
