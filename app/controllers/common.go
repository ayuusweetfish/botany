package controllers

import (
	"net/http"
	"strconv"
)

// URL query string: count=...&page=...
// Page is 0-indexed
// Returns count and offset in elements (not pages)
func parsePagination(w http.ResponseWriter, r *http.Request) (int, int) {
	limits := r.URL.Query()["count"]
	pages := r.URL.Query()["page"]
	if limits != nil && pages != nil {
		limit, _ := strconv.Atoi(r.URL.Query()["count"][0])
		page, _ := strconv.Atoi(r.URL.Query()["page"][0])
		if limit >= 0 && page >= 0 {
			return limit, page * limit
		}
	}
	return -1, -1
}
