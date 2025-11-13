package deduplicator

import (
	"github.com/lcalzada-xor/dupdurl/pkg/stats"
)

// Entry represents a deduplicated URL with its count
type Entry struct {
	URL   string `json:"url"`
	Count int    `json:"count"`
}

// Deduplicator handles URL deduplication
type Deduplicator struct {
	seen   map[string]string // dedup key -> first full URL with values
	counts map[string]int    // dedup key -> occurrence count
	order  []string          // preserve first-appearance order
	stats  *stats.Statistics
}

// New creates a new Deduplicator instance
func New(s *stats.Statistics) *Deduplicator {
	return &Deduplicator{
		seen:   make(map[string]string),
		counts: make(map[string]int),
		order:  make([]string, 0),
		stats:  s,
	}
}

// Add adds a URL to the deduplicator
// dedupKey is used for comparison, normalizedURL is stored for output
func (d *Deduplicator) Add(dedupKey, normalizedURL string) {
	if _, exists := d.seen[dedupKey]; !exists {
		d.seen[dedupKey] = normalizedURL
		d.order = append(d.order, dedupKey)
		if d.stats != nil {
			d.stats.UniqueURLs++
		}
	} else {
		if d.stats != nil {
			d.stats.Duplicates++
		}
	}
	d.counts[dedupKey]++
}

// GetEntries returns all deduplicated entries in first-seen order
func (d *Deduplicator) GetEntries() []Entry {
	entries := make([]Entry, len(d.order))
	for i, key := range d.order {
		entries[i] = Entry{
			URL:   d.seen[key],
			Count: d.counts[key],
		}
	}
	return entries
}

// Count returns the number of unique entries
func (d *Deduplicator) Count() int {
	return len(d.order)
}

// Clear resets the deduplicator state
func (d *Deduplicator) Clear() {
	d.seen = make(map[string]string)
	d.counts = make(map[string]int)
	d.order = make([]string, 0)
}
