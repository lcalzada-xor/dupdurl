package storage

import "github.com/lcalzada-xor/dupdurl/pkg/deduplicator"

// Backend defines the interface for storage backends
type Backend interface {
	// Add stores a URL with its dedup key
	Add(dedupKey, url string) error

	// GetEntries retrieves all stored entries
	GetEntries() ([]deduplicator.Entry, error)

	// Count returns the number of unique entries
	Count() int

	// Close closes the backend and releases resources
	Close() error
}
