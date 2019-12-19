package models

import (
	"database/sql"
	"strings"
)

var db *sql.DB

type TableSchema struct {
	table   string
	columns []string
}

var schemata []TableSchema

func registerSchema(table string, columns ...string) {
	schemata = append(schemata, TableSchema{table, columns})
}

func InitializeSchemata(dbInput *sql.DB) {
	db = dbInput
	for _, schema := range schemata {
		cmd := "CREATE TABLE IF NOT EXISTS " + schema.table + " ()"
		println(cmd)
		db.Exec(cmd)
		for _, columnDesc := range schema.columns {
			columnName := strings.SplitN(columnDesc, " ", 2)[0]
			if columnName != "ADD" {
				// Column
				row := db.QueryRow("SELECT COUNT(*) FROM information_schema.columns "+
					"WHERE table_name = $1 AND column_name = $2",
					schema.table,
					columnName,
				)
				var count int
				if err := row.Scan(&count); err != nil {
					panic(err)
				}
				if count == 0 {
					schema := "ALTER TABLE " + schema.table + " ADD COLUMN " + columnDesc
					println(schema)
					db.Exec(schema)
				} else {
					// println("Column " + columnName + " already exists, skipping")
				}
			}
		}
	}
	for _, schema := range schemata {
		for _, columnDesc := range schema.columns {
			columnName := strings.SplitN(columnDesc, " ", 2)[0]
			if columnName == "ADD" {
				// Constraint
				schema := "ALTER TABLE " + schema.table + " " + columnDesc
				println(schema)
				db.Exec(schema)
			}
		}
	}
}
