package models

import (
	_ "database/sql"
)

type File struct {
	Id      int32
	Type    string
	Content []byte
}

func init() {
	registerSchema("file",
		"id SERIAL PRIMARY KEY",
		"type TEXT NOT NULL",
		"content BYTEA NOT NULL",
	)
}

func (f *File) Create() error {
	return db.QueryRow("INSERT INTO "+
		"file(type, content) VALUES ($1, $2) RETURNING id",
		f.Type,
		f.Content,
	).Scan(&f.Id)
}

func (f *File) Read() error {
	return db.QueryRow("SELECT type, content FROM file WHERE id = $1",
		f.Id).Scan(&f.Type, &f.Content)
}

func (f *File) Update() error {
	_, err := db.Exec("UPDATE file SET "+
		"type = $1, content = $2 WHERE id = $3",
		f.Type, f.Content, f.Id)
	return err
}
