package db

import "database/sql"

var db *sql.DB

func createConnection() error {
	if db == nil {
		newConnection, err := sql.Open("sqlite3", "gorest.db")

		if err != nil {
			return err
		}

		db = newConnection
	}

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

func init() {
	err := createConnection()

	if err != nil {
		panic(err)
	}

	// Define the SQL statement to create the table.
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS requests (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT,
        name TEXT,
        method TEXT,
        headers TEXT,
        body TEXT
    );
    `

	// Execute the SQL statement to create the table.
	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
}
