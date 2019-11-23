package controllers

import (
	"database/sql"
	"net/http"
)

func rowsToMap(db *sql.DB, sql_info string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(sql_info, args)
	if err != nil {
		return nil, err
	}
	columns, _ := rows.Columns()
	columnLen := len(columns)
	cache := make([]interface{}, columnLen)
	for i, _ := range cache {
		var a interface{}
		cache[i] = &a
	}
	var rlt []map[string]interface{}
	for rows.Next() {
		rows.Scan(cache...)
		data := make(map[string]interface{})
		for i, _ := range cache {
			data[columns[i]] = cache[i]
		}
		rlt = append(rlt, data)
	}
	return rlt, nil
}

func contestlistHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)
	// todo: add models/contest and use Query to get data.
	// rowsToMap()
	// fmt.Fprintf(w, `{"total": `))
}

func init() {
	registerRouterFunc("/gamelist", gamelistHandler)
}
