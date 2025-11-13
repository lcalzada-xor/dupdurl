package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/lcalzada-xor/dupdurl/pkg/deduplicator"
)

// SQLiteBackend stores URLs in SQLite database for massive datasets
type SQLiteBackend struct {
	db *sql.DB
}

// NewSQLiteBackend creates a new SQLite storage backend
func NewSQLiteBackend(dbPath string) (*SQLiteBackend, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	backend := &SQLiteBackend{db: db}
	if err := backend.initialize(); err != nil {
		db.Close()
		return nil, err
	}

	return backend, nil
}

// initialize creates the necessary tables
func (s *SQLiteBackend) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dedup_key TEXT UNIQUE NOT NULL,
		url TEXT NOT NULL,
		count INTEGER DEFAULT 1,
		first_seen INTEGER DEFAULT (strftime('%s', 'now'))
	);
	CREATE INDEX IF NOT EXISTS idx_dedup_key ON urls(dedup_key);
	CREATE INDEX IF NOT EXISTS idx_first_seen ON urls(first_seen);
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// Add stores or updates a URL in the database
func (s *SQLiteBackend) Add(dedupKey, url string) error {
	query := `
	INSERT INTO urls (dedup_key, url, count)
	VALUES (?, ?, 1)
	ON CONFLICT(dedup_key) DO UPDATE SET count = count + 1
	`

	_, err := s.db.Exec(query, dedupKey, url)
	if err != nil {
		return fmt.Errorf("failed to insert URL: %w", err)
	}

	return nil
}

// GetEntries retrieves all stored entries ordered by first-seen
func (s *SQLiteBackend) GetEntries() ([]deduplicator.Entry, error) {
	query := `SELECT url, count FROM urls ORDER BY first_seen`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query URLs: %w", err)
	}
	defer rows.Close()

	var entries []deduplicator.Entry
	for rows.Next() {
		var entry deduplicator.Entry
		if err := rows.Scan(&entry.URL, &entry.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return entries, nil
}

// Count returns the number of unique entries
func (s *SQLiteBackend) Count() int {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM urls").Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

// Close closes the database connection
func (s *SQLiteBackend) Close() error {
	return s.db.Close()
}

// Clear removes all entries from the database
func (s *SQLiteBackend) Clear() error {
	_, err := s.db.Exec("DELETE FROM urls")
	return err
}
