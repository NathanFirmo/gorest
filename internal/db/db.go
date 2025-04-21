package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Request struct {
	ID      int64
	URL     string
	Name    string
	Method  string
	Headers string
	Body    string
}

var db *sql.DB

func getDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dbDir := filepath.Join(home, "gorest")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		panic(err)
	}
	return filepath.Join(dbDir, ".db")
}

func createConnection() (*sql.DB, error) {
	return sql.Open("sqlite3", getDBPath())
}

func init() {
	var err error
	db, err = createConnection()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT,
			name TEXT,
			method TEXT,
			headers TEXT,
			body TEXT
		)
	`)
	if err != nil {
		panic(err)
	}
}

func SaveRequest(req *Request) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO requests (id, url, name, method, headers, body) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(req.ID, req.URL, req.Name, req.Method, req.Headers, req.Body)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func UpdateRequest(req *Request) error {
	stmt, err := db.Prepare("UPDATE requests SET url = ?, name = ?, method = ?, headers = ?, body = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.URL, req.Name, req.Method, req.Headers, req.Body, req.ID)
	return err
}

func DeleteRequest(id int64) error {
	stmt, err := db.Prepare("DELETE FROM requests WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func GetAllRequests() ([]Request, error) {
	rows, err := db.Query("SELECT id, url, name, method, headers, body FROM requests")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var req Request
		err := rows.Scan(&req.ID, &req.URL, &req.Name, &req.Method, &req.Headers, &req.Body)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, nil
}
