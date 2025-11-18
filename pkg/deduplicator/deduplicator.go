package deduplicator

import (
	"github.com/lcalzada-xor/dupdurl/pkg/locale"
	"github.com/lcalzada-xor/dupdurl/pkg/stats"
)

// Entry represents a deduplicated URL with its count
type Entry struct {
	URL   string `json:"url"`
	Count int    `json:"count"`
}

// Deduplicator handles URL deduplication
type Deduplicator struct {
	seen          map[string]string            // dedup key -> first full URL with values
	counts        map[string]int               // dedup key -> occurrence count
	order         []string                     // preserve first-appearance order
	stats         *stats.Statistics
	localeGroups  map[string]*locale.LocaleGroup // locale-aware grouping
	grouper       *locale.Grouper
	localeAware   bool
	originalURLs  map[string]string            // dedup key -> original URL before normalization
}

// New creates a new Deduplicator instance
func New(s *stats.Statistics) *Deduplicator {
	return &Deduplicator{
		seen:         make(map[string]string),
		counts:       make(map[string]int),
		order:        make([]string, 0),
		stats:        s,
		localeGroups: make(map[string]*locale.LocaleGroup),
		grouper:      nil,
		localeAware:  false,
		originalURLs: make(map[string]string),
	}
}

// NewWithLocaleSupport creates a new Deduplicator with locale awareness
func NewWithLocaleSupport(s *stats.Statistics, localePriority []string) *Deduplicator {
	if len(localePriority) == 0 {
		localePriority = []string{"en"}
	}

	return &Deduplicator{
		seen:         make(map[string]string),
		counts:       make(map[string]int),
		order:        make([]string, 0),
		stats:        s,
		localeGroups: make(map[string]*locale.LocaleGroup),
		grouper:      locale.NewGrouper(localePriority),
		localeAware:  true,
		originalURLs: make(map[string]string),
	}
}

// SetLocaleAware enables or disables locale awareness
func (d *Deduplicator) SetLocaleAware(enabled bool, priority []string) {
	d.localeAware = enabled
	if enabled && d.grouper == nil {
		if len(priority) == 0 {
			priority = []string{"en"}
		}
		d.grouper = locale.NewGrouper(priority)
	}
}

// Add adds a URL to the deduplicator
// dedupKey is used for comparison, normalizedURL is stored for output
func (d *Deduplicator) Add(dedupKey, normalizedURL string) {
	// Standard deduplication logic
	if _, exists := d.seen[dedupKey]; !exists {
		d.seen[dedupKey] = normalizedURL
		d.order = append(d.order, dedupKey)
		d.originalURLs[dedupKey] = normalizedURL
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

// AddWithOriginal adds a URL with both normalized and original versions
func (d *Deduplicator) AddWithOriginal(dedupKey, normalizedURL, originalURL string) {
	// If locale-aware mode is enabled, also track in grouper
	if d.localeAware && d.grouper != nil {
		// Add original URL to locale grouper
		d.grouper.Add(originalURL)
	}

	// Standard deduplication logic
	if _, exists := d.seen[dedupKey]; !exists {
		d.seen[dedupKey] = normalizedURL
		d.order = append(d.order, dedupKey)
		d.originalURLs[dedupKey] = originalURL
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
	// If locale-aware mode is enabled, get best URLs from grouper
	if d.localeAware && d.grouper != nil {
		bestURLs := d.grouper.GetBestURLs()
		entries := make([]Entry, 0, len(bestURLs))

		// For each best URL, find its dedup key and get the count
		seenKeys := make(map[string]bool)

		for _, locURL := range bestURLs {
			// Find the dedup key for this URL
			for key, origURL := range d.originalURLs {
				if origURL == locURL.OriginalURL && !seenKeys[key] {
					entries = append(entries, Entry{
						URL:   d.seen[key],
						Count: d.counts[key],
					})
					seenKeys[key] = true
					break
				}
			}
		}

		return entries
	}

	// Standard mode: return all entries
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
	d.localeGroups = make(map[string]*locale.LocaleGroup)
	d.originalURLs = make(map[string]string)
	if d.localeAware && d.grouper != nil {
		// Reset grouper
		priority := d.grouper.Priority
		d.grouper = locale.NewGrouper(priority)
	}
}

// GetLocaleGroups returns locale groups for debugging/stats
func (d *Deduplicator) GetLocaleGroups() map[string]*locale.LocaleGroup {
	if d.grouper != nil {
		return d.grouper.GetGroups()
	}
	return nil
}
