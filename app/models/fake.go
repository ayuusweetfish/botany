package models

func FakeDatabase() {
	// Clear database
	for _, schema := range schemata {
		db.Exec("DROP TABLE IF EXISTS " + schema.table + " CASCADE")
	}

	InitializeSchemata(db)
}
