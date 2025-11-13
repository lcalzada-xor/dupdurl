package storage

import (
	"github.com/lcalzada-xor/dupdurl/pkg/deduplicator"
)

// MemoryBackend stores URLs in memory (default behavior)
type MemoryBackend struct {
	seen   map[string]string // dedup key -> first full URL with values
	counts map[string]int    // dedup key -> occurrence count
	order  []string          // preserve first-appearance order
}

// NewMemoryBackend creates a new in-memory storage backend
func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{
		seen:   make(map[string]string),
		counts: make(map[string]int),
		order:  make([]string, 0),
	}
}

// Add stores a URL in memory
func (m *MemoryBackend) Add(dedupKey, url string) error {
	if _, exists := m.seen[dedupKey]; !exists {
		m.seen[dedupKey] = url
		m.order = append(m.order, dedupKey)
	}
	m.counts[dedupKey]++
	return nil
}

// GetEntries returns all stored entries in first-seen order
func (m *MemoryBackend) GetEntries() ([]deduplicator.Entry, error) {
	entries := make([]deduplicator.Entry, len(m.order))
	for i, key := range m.order {
		entries[i] = deduplicator.Entry{
			URL:   m.seen[key],
			Count: m.counts[key],
		}
	}
	return entries, nil
}

// Count returns the number of unique entries
func (m *MemoryBackend) Count() int {
	return len(m.order)
}

// Close is a no-op for memory backend
func (m *MemoryBackend) Close() error {
	return nil
}
