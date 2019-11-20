package globals

import (
	"database/sql"
	"github.com/gorilla/sessions"
)

var DB *sql.DB
var SessionStore *sessions.CookieStore
