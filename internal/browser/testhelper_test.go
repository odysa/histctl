package browser

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

const safariSchema = `
CREATE TABLE history_items (
	id INTEGER PRIMARY KEY,
	url TEXT
);
CREATE TABLE history_visits (
	id INTEGER PRIMARY KEY,
	history_item INTEGER REFERENCES history_items(id),
	title TEXT,
	visit_time REAL
);`

const chromeSchema = `
CREATE TABLE urls (
	id INTEGER PRIMARY KEY,
	url TEXT,
	title TEXT,
	visit_count INTEGER DEFAULT 0
);
CREATE TABLE visits (
	id INTEGER PRIMARY KEY,
	url INTEGER REFERENCES urls(id),
	visit_time INTEGER
);`

const firefoxSchema = `
CREATE TABLE moz_places (
	id INTEGER PRIMARY KEY,
	url TEXT,
	title TEXT
);
CREATE TABLE moz_historyvisits (
	id INTEGER PRIMARY KEY,
	place_id INTEGER REFERENCES moz_places(id),
	visit_date INTEGER
);`

// createTestDB creates a temporary SQLite database with the given schema
// and returns its path. The database is removed when the test completes.
func createTestDB(t *testing.T, schema string) string {
	t.Helper()
	path := t.TempDir() + "/test.db"
	db, err := sql.Open("sqlite", "file:"+path)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	defer db.Close()
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	return path
}

// seedSafari inserts test rows into a Safari-schema database.
func seedSafari(t *testing.T, path string, rows []safariRow) {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+path)
	if err != nil {
		t.Fatalf("open db for seeding: %v", err)
	}
	defer db.Close()
	for _, r := range rows {
		if _, err := db.Exec("INSERT INTO history_items (id, url) VALUES (?, ?)", r.id, r.url); err != nil {
			t.Fatalf("seed history_items: %v", err)
		}
		if _, err := db.Exec("INSERT INTO history_visits (history_item, title, visit_time) VALUES (?, ?, ?)", r.id, r.title, r.visitTime); err != nil {
			t.Fatalf("seed history_visits: %v", err)
		}
	}
}

type safariRow struct {
	id        int64
	url       *string // nil for NULL
	title     string
	visitTime float64
}

// seedChrome inserts test rows into a Chrome-schema database.
func seedChrome(t *testing.T, path string, rows []chromeRow) {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+path)
	if err != nil {
		t.Fatalf("open db for seeding: %v", err)
	}
	defer db.Close()
	for _, r := range rows {
		if _, err := db.Exec("INSERT INTO urls (id, url, title) VALUES (?, ?, ?)", r.id, r.url, r.title); err != nil {
			t.Fatalf("seed urls: %v", err)
		}
		if _, err := db.Exec("INSERT INTO visits (url, visit_time) VALUES (?, ?)", r.id, r.visitTime); err != nil {
			t.Fatalf("seed visits: %v", err)
		}
	}
}

type chromeRow struct {
	id        int64
	url       *string // nil for NULL
	title     string
	visitTime int64
}

// seedFirefox inserts test rows into a Firefox-schema database.
func seedFirefox(t *testing.T, path string, rows []firefoxRow) {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+path)
	if err != nil {
		t.Fatalf("open db for seeding: %v", err)
	}
	defer db.Close()
	for _, r := range rows {
		if _, err := db.Exec("INSERT INTO moz_places (id, url, title) VALUES (?, ?, ?)", r.id, r.url, r.title); err != nil {
			t.Fatalf("seed moz_places: %v", err)
		}
		if _, err := db.Exec("INSERT INTO moz_historyvisits (place_id, visit_date) VALUES (?, ?)", r.id, r.visitDate); err != nil {
			t.Fatalf("seed moz_historyvisits: %v", err)
		}
	}
}

type firefoxRow struct {
	id        int64
	url       *string // nil for NULL
	title     string
	visitDate int64
}

func strPtr(s string) *string { return &s }
