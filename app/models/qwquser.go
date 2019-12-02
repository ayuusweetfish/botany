package models

type QwQUser struct {
	Name  string
	Count int32
}

func init() {
	registerSchema("qwq_user",
		"name TEXT",
		"count INTEGER",
	)
}

func (u *QwQUser) Create() error {
	_, err := db.Exec("INSERT INTO qwq_user(name, count) VALUES ($1, $2)", u.Name, u.Count)
	return err
}

func (u *QwQUser) Read() error {
	row := db.QueryRow("SELECT count FROM qwq_user WHERE name = $1", u.Name)
	err := row.Scan(&u.Count)
	return err
}

func (u *QwQUser) Update() error {
	_, err := db.Exec("UPDATE qwq_user SET count = $1 WHERE name = $2", u.Count, u.Name)
	return err
}
