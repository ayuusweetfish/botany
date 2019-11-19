package models

import (
	"../globals"

	"strings"
)

var schemata []string

func registerSchema(sql string) {
	schemata = append(schemata, sql)
}

func InitializeSchemata() {
	s := strings.Join(schemata, "\n")
	globals.DB.Exec(schema)
	println(s)
}
