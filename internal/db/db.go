package db

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "modernc.org/sqlite"
)

type Request struct {
	ID      int64
	URL     string
	Name    string
	Method  string
	Headers string
	Body    string
}

type Store struct {
	db *sql.DB
}

const schemaVersion = 1

// DefaultPath returns where requests are stored. GOREST_DB overrides it, and a
// database left by older versions at ~/gorest/.db is moved to the new location.
func DefaultPath() (string, error) {
	if path := os.Getenv("GOREST_DB"); path != "" {
		return path, nil
	}

	dir, err := dataDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "gorest", "gorest.db")

	home, err := os.UserHomeDir()
	if err != nil {
		return path, nil
	}
	legacyDir := filepath.Join(home, "gorest")
	legacyPath := filepath.Join(legacyDir, ".db")
	if fileExists(legacyPath) && !fileExists(path) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return "", err
		}
		if err := os.Rename(legacyPath, path); err != nil {
			return legacyPath, nil
		}
		_ = os.Remove(legacyDir)
	}

	return path, nil
}

func dataDir() (string, error) {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return dir, nil
	}
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		return os.UserConfigDir()
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share"), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1)

	store := &Store{db: conn}
	if err := store.migrate(); err != nil {
		conn.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	var version int
	if err := s.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return err
	}
	if version >= schemaVersion {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT,
			name TEXT,
			method TEXT,
			headers TEXT,
			body TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
	}
	for _, stmt := range statements {
		if _, err := tx.Exec(stmt); err != nil {
			return err
		}
	}

	if err := splitLegacyMethods(tx); err != nil {
		return err
	}

	if _, err := tx.Exec("PRAGMA user_version = 1"); err != nil {
		return err
	}
	return tx.Commit()
}

// Older versions stored the method inside the URL field, e.g. "POST http://...".
func splitLegacyMethods(tx *sql.Tx) error {
	rows, err := tx.Query("SELECT id, COALESCE(url, '') FROM requests")
	if err != nil {
		return err
	}

	type fix struct {
		id          int64
		method, url string
	}
	var fixes []fix
	for rows.Next() {
		var id int64
		var url string
		if err := rows.Scan(&id, &url); err != nil {
			rows.Close()
			return err
		}
		method, rest, ok := strings.Cut(strings.TrimSpace(url), " ")
		if ok && isHTTPMethod(method) {
			fixes = append(fixes, fix{id: id, method: strings.ToUpper(method), url: strings.TrimSpace(rest)})
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	for _, f := range fixes {
		if _, err := tx.Exec("UPDATE requests SET method = ?, url = ? WHERE id = ?", f.method, f.url, f.id); err != nil {
			return err
		}
	}
	return nil
}

func isHTTPMethod(s string) bool {
	switch strings.ToUpper(s) {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		return true
	}
	return false
}

func (s *Store) CreateRequest(req *Request) error {
	result, err := s.db.Exec(
		"INSERT INTO requests (url, name, method, headers, body) VALUES (?, ?, ?, ?, ?)",
		req.URL, req.Name, req.Method, req.Headers, req.Body,
	)
	if err != nil {
		return err
	}
	req.ID, err = result.LastInsertId()
	return err
}

func (s *Store) UpdateRequest(req Request) error {
	_, err := s.db.Exec(
		"UPDATE requests SET url = ?, name = ?, method = ?, headers = ?, body = ? WHERE id = ?",
		req.URL, req.Name, req.Method, req.Headers, req.Body, req.ID,
	)
	return err
}

func (s *Store) DeleteRequest(id int64) error {
	_, err := s.db.Exec("DELETE FROM requests WHERE id = ?", id)
	return err
}

func (s *Store) ListRequests() ([]Request, error) {
	rows, err := s.db.Query(`
		SELECT id, COALESCE(url, ''), COALESCE(name, ''), COALESCE(method, ''), COALESCE(headers, ''), COALESCE(body, '')
		FROM requests ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var req Request
		if err := rows.Scan(&req.ID, &req.URL, &req.Name, &req.Method, &req.Headers, &req.Body); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

func (s *Store) Setting(key string) (string, error) {
	var value string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return value, err
}

func (s *Store) SetSetting(key, value string) error {
	_, err := s.db.Exec(
		"INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value",
		key, value,
	)
	return err
}
