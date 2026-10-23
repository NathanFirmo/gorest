package db

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func openTestStore(t *testing.T) *Store {
	t.Helper()
	store, err := Open(filepath.Join(t.TempDir(), "gorest.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestRequestsCRUD(t *testing.T) {
	store := openTestStore(t)

	first := Request{Name: "List users", Method: "GET", URL: "https://example.com/users"}
	second := Request{Name: "Create user", Method: "POST", URL: "https://example.com/users", Body: `{"a":1}`}
	for _, req := range []*Request{&first, &second} {
		if err := store.CreateRequest(req); err != nil {
			t.Fatal(err)
		}
	}
	if first.ID == 0 || second.ID == first.ID {
		t.Fatalf("unexpected ids %d and %d", first.ID, second.ID)
	}

	second.Headers = "Content-Type: application/json"
	if err := store.UpdateRequest(second); err != nil {
		t.Fatal(err)
	}
	if err := store.DeleteRequest(first.ID); err != nil {
		t.Fatal(err)
	}

	requests, err := store.ListRequests()
	if err != nil {
		t.Fatal(err)
	}
	if len(requests) != 1 || requests[0] != second {
		t.Fatalf("ListRequests() = %#v, want only %#v", requests, second)
	}
}

func TestSettings(t *testing.T) {
	store := openTestStore(t)

	if value, err := store.Setting("theme"); err != nil || value != "" {
		t.Fatalf("Setting() = %q, %v", value, err)
	}
	for _, theme := range []string{"nord", "dracula"} {
		if err := store.SetSetting("theme", theme); err != nil {
			t.Fatal(err)
		}
	}
	if value, _ := store.Setting("theme"); value != "dracula" {
		t.Fatalf("Setting() = %q, want dracula", value)
	}
}

func TestMigratesLegacyDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.db")
	legacy, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = legacy.Exec(`
		CREATE TABLE requests (id INTEGER PRIMARY KEY AUTOINCREMENT, url TEXT, name TEXT, method TEXT, headers TEXT, body TEXT);
		INSERT INTO requests (id, url, name, method, headers, body) VALUES
			(0, 'POST https://example.com/users', 'Create', 'GET', '', ''),
			(1, 'https://example.com', 'Plain', 'GET', NULL, NULL);
	`)
	legacy.Close()
	if err != nil {
		t.Fatal(err)
	}

	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	requests, err := store.ListRequests()
	if err != nil {
		t.Fatal(err)
	}
	want := []Request{
		{ID: 0, Name: "Create", Method: "POST", URL: "https://example.com/users"},
		{ID: 1, Name: "Plain", Method: "GET", URL: "https://example.com"},
	}
	if len(requests) != len(want) || requests[0] != want[0] || requests[1] != want[1] {
		t.Fatalf("ListRequests() = %#v, want %#v", requests, want)
	}
}
