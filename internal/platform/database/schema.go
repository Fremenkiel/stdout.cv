package database

import "database/sql"

const schema = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	)
`

func InitializeSchema(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
