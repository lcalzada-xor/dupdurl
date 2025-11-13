package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/lcalzada-xor/dupdurl/pkg/deduplicator"
)

// DiffReport represents the differences between two URL sets
type DiffReport struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Changed []Change `json:"changed"`
}

// Change represents a URL that exists in both sets but with different counts
type Change struct {
	URL      string `json:"url"`
	OldCount int    `json:"old_count"`
	NewCount int    `json:"new_count"`
}

// Differ compares URL sets
type Differ struct {
	baseline map[string]int // URL -> count
}

// NewDiffer creates a new Differ instance
func NewDiffer() *Differ {
	return &Differ{
		baseline: make(map[string]int),
	}
}

// LoadBaseline loads baseline URLs from a JSON file
func (d *Differ) LoadBaseline(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read baseline file: %w", err)
	}

	var entries []deduplicator.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to parse baseline JSON: %w", err)
	}

	for _, entry := range entries {
		d.baseline[entry.URL] = entry.Count
	}

	return nil
}

// LoadBaselineFromEntries loads baseline from entry slice
func (d *Differ) LoadBaselineFromEntries(entries []deduplicator.Entry) {
	d.baseline = make(map[string]int, len(entries))
	for _, entry := range entries {
		d.baseline[entry.URL] = entry.Count
	}
}

// Compare compares current entries against baseline
func (d *Differ) Compare(current []deduplicator.Entry) *DiffReport {
	report := &DiffReport{
		Added:   []string{},
		Removed: []string{},
		Changed: []Change{},
	}

	// Track which baseline URLs we've seen
	seen := make(map[string]struct{}, len(d.baseline))

	// Check for added and changed URLs
	for _, entry := range current {
		oldCount, existed := d.baseline[entry.URL]

		if !existed {
			// New URL
			report.Added = append(report.Added, entry.URL)
		} else {
			// Existed in baseline
			seen[entry.URL] = struct{}{}

			// Check if count changed
			if entry.Count != oldCount {
				report.Changed = append(report.Changed, Change{
					URL:      entry.URL,
					OldCount: oldCount,
					NewCount: entry.Count,
				})
			}
		}
	}

	// Check for removed URLs (in baseline but not in current)
	for url := range d.baseline {
		if _, stillExists := seen[url]; !stillExists {
			report.Removed = append(report.Removed, url)
		}
	}

	return report
}

// PrintReport prints a human-readable diff report
func (r *DiffReport) PrintReport(w io.Writer) {
	if len(r.Added) > 0 {
		fmt.Fprintf(w, "\n[ADDED] %d new URLs:\n", len(r.Added))
		for _, url := range r.Added {
			fmt.Fprintf(w, "  + %s\n", url)
		}
	}

	if len(r.Removed) > 0 {
		fmt.Fprintf(w, "\n[REMOVED] %d URLs:\n", len(r.Removed))
		for _, url := range r.Removed {
			fmt.Fprintf(w, "  - %s\n", url)
		}
	}

	if len(r.Changed) > 0 {
		fmt.Fprintf(w, "\n[CHANGED] %d URLs with different counts:\n", len(r.Changed))
		for _, change := range r.Changed {
			fmt.Fprintf(w, "  ~ %s (%d -> %d)\n", change.URL, change.OldCount, change.NewCount)
		}
	}

	if len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0 {
		fmt.Fprintln(w, "\nNo differences found.")
	}
}

// ToJSON converts report to JSON
func (r *DiffReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// Summary returns a summary of the diff
func (r *DiffReport) Summary() string {
	return fmt.Sprintf("Added: %d, Removed: %d, Changed: %d",
		len(r.Added), len(r.Removed), len(r.Changed))
}

// SaveBaseline saves current entries as baseline JSON file
func SaveBaseline(entries []deduplicator.Entry, path string) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal entries: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write baseline file: %w", err)
	}

	return nil
}
