package unit

import (
	"testing"

	"github.com/lcalzada-xor/dupdurl/pkg/deduplicator"
	"github.com/lcalzada-xor/dupdurl/pkg/stats"
)

func TestDeduplicatorBasic(t *testing.T) {
	st := stats.NewStatistics()
	dedup := deduplicator.New(st)

	// Add first URL
	dedup.Add("key1", "https://example.com/page1")
	if dedup.Count() != 1 {
		t.Errorf("Count after first add = %d; want 1", dedup.Count())
	}

	// Add duplicate (same key, different URL)
	dedup.Add("key1", "https://example.com/page1?param=value")
	if dedup.Count() != 1 {
		t.Errorf("Count after duplicate = %d; want 1", dedup.Count())
	}

	// Add second unique URL
	dedup.Add("key2", "https://example.com/page2")
	if dedup.Count() != 2 {
		t.Errorf("Count after second unique = %d; want 2", dedup.Count())
	}

	entries := dedup.GetEntries()
	if len(entries) != 2 {
		t.Errorf("GetEntries() length = %d; want 2", len(entries))
	}

	// Verify first entry
	if entries[0].URL != "https://example.com/page1" {
		t.Errorf("First entry URL = %q; want https://example.com/page1", entries[0].URL)
	}
	if entries[0].Count != 2 {
		t.Errorf("First entry count = %d; want 2", entries[0].Count)
	}

	// Verify second entry
	if entries[1].URL != "https://example.com/page2" {
		t.Errorf("Second entry URL = %q; want https://example.com/page2", entries[1].URL)
	}
	if entries[1].Count != 1 {
		t.Errorf("Second entry count = %d; want 1", entries[1].Count)
	}
}

func TestDeduplicatorOrder(t *testing.T) {
	st := stats.NewStatistics()
	dedup := deduplicator.New(st)

	// Add URLs in specific order
	urls := []struct {
		key string
		url string
	}{
		{"key3", "url3"},
		{"key1", "url1"},
		{"key2", "url2"},
		{"key1", "url1_dup"}, // Duplicate
	}

	for _, u := range urls {
		dedup.Add(u.key, u.url)
	}

	entries := dedup.GetEntries()

	// Verify order matches first-seen order
	expectedOrder := []string{"url3", "url1", "url2"}
	if len(entries) != len(expectedOrder) {
		t.Errorf("GetEntries() length = %d; want %d", len(entries), len(expectedOrder))
	}

	for i, expected := range expectedOrder {
		if entries[i].URL != expected {
			t.Errorf("Entry[%d] URL = %q; want %q", i, entries[i].URL, expected)
		}
	}
}

func TestDeduplicatorStatistics(t *testing.T) {
	st := stats.NewStatistics()
	dedup := deduplicator.New(st)

	// Add URLs
	dedup.Add("key1", "url1")
	dedup.Add("key1", "url1_dup")
	dedup.Add("key2", "url2")
	dedup.Add("key3", "url3")
	dedup.Add("key2", "url2_dup")

	// Verify statistics
	if st.UniqueURLs != 3 {
		t.Errorf("UniqueURLs = %d; want 3", st.UniqueURLs)
	}
	if st.Duplicates != 2 {
		t.Errorf("Duplicates = %d; want 2", st.Duplicates)
	}
}

func TestDeduplicatorClear(t *testing.T) {
	st := stats.NewStatistics()
	dedup := deduplicator.New(st)

	// Add URLs
	dedup.Add("key1", "url1")
	dedup.Add("key2", "url2")

	// Clear
	dedup.Clear()

	if dedup.Count() != 0 {
		t.Errorf("Count after Clear() = %d; want 0", dedup.Count())
	}

	entries := dedup.GetEntries()
	if len(entries) != 0 {
		t.Errorf("GetEntries() after Clear() length = %d; want 0", len(entries))
	}
}
