package models

type QwQUser struct {
	Name string
	Count int32
}

var schema = `
CREATE TABLE visitor (
	name TEXT,
	count INTEGER
)`

func (u *QwQUser) Create() error {
	_, err := db.Exec("INSERT INTO visitor(name, count) VALUES ($1, $2)", u.Name, u.Count)
	return err
}

func (u *QwQUser) Read() error {
	row := db.QueryRow("SELECT count FROM visitor WHERE name = $1", u.Name)
	err := row.Scan(&u.Count)
	return err
}

func (u *QwQUser) Update() error {
	_, err := db.Exec("UPDATE visitor SET count = $1 WHERE name = $2", u.Count, u.Name)
	return err
}

func init() {
	registerSchema(schema)
}
