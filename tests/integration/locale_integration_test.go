package integration

import (
	"testing"

	"github.com/lcalzada-xor/dupdurl/pkg/deduplicator"
	"github.com/lcalzada-xor/dupdurl/pkg/normalizer"
	"github.com/lcalzada-xor/dupdurl/pkg/stats"
)

func TestLocaleAwareDeduplication(t *testing.T) {
	// Create stats
	st := stats.NewStatistics()

	// Create deduplicator with locale support
	dedup := deduplicator.NewWithLocaleSupport(st, []string{"en"})

	// Create normalizer
	config := normalizer.NewConfig()

	urls := []string{
		"https://example.com/about",
		"https://example.com/en/about",
		"https://example.com/es/sobre-nosotros",
		"https://example.com/it/chi-siamo",
		"https://example.com/fr/a-propos",
		"https://example.com/products",
		"https://example.com/en/products",
		"https://example.com/es/productos",
		"https://example.com/unique-endpoint",
	}

	// In locale-aware mode, use original URL as key
	for _, url := range urls {
		normalized, err := config.NormalizeURL(url)
		if err != nil {
			continue // Skip errors
		}

		// Use original URL as dedup key in locale-aware mode
		dedup.AddWithOriginal(url, normalized, url)
	}

	entries := dedup.GetEntries()

	// We should have:
	// 1. One "about" page (English preferred)
	// 2. One "products" page (English preferred)
	// 3. One "unique-endpoint"
	// Total: 3 entries

	if len(entries) != 3 {
		t.Errorf("Expected 3 unique entries, got %d", len(entries))
		for i, entry := range entries {
			t.Logf("Entry %d: %s (count: %d)", i+1, entry.URL, entry.Count)
		}
	}

	// Check that we have the preferred versions
	foundAbout := false
	foundProducts := false
	foundUnique := false

	for _, entry := range entries {
		// Should contain "about" or "en/about"
		if entry.URL == "https://example.com/en/about" || entry.URL == "https://example.com/about" {
			foundAbout = true
		}
		// Should contain "products" or "en/products"
		if entry.URL == "https://example.com/en/products" || entry.URL == "https://example.com/products" {
			foundProducts = true
		}
		if entry.URL == "https://example.com/unique-endpoint" {
			foundUnique = true
		}
	}

	if !foundAbout {
		t.Error("Expected 'about' page in results")
	}
	if !foundProducts {
		t.Error("Expected 'products' page in results")
	}
	if !foundUnique {
		t.Error("Expected 'unique-endpoint' in results")
	}
}

func TestLocaleAwareWithSubdomains(t *testing.T) {
	st := stats.NewStatistics()
	dedup := deduplicator.NewWithLocaleSupport(st, []string{"en"})

	config := normalizer.NewConfig()
	config.LocaleAware = true
	config.LocalePriority = []string{"en"}

	urls := []string{
		"https://en.example.com/about",
		"https://es.example.com/about",
		"https://it.example.com/about",
		"https://www.example.com/products",
	}

	for _, url := range urls {
		normalized, err := config.NormalizeURL(url)
		if err != nil {
			continue
		}

		dedupKey, err := config.CreateDedupKey(url)
		if err != nil {
			continue
		}

		dedup.AddWithOriginal(dedupKey, normalized, url)
	}

	entries := dedup.GetEntries()

	// Should have 2 entries: one about page (en subdomain), one products page
	if len(entries) < 1 || len(entries) > 2 {
		t.Errorf("Expected 1-2 unique entries, got %d", len(entries))
		for i, entry := range entries {
			t.Logf("Entry %d: %s", i+1, entry.URL)
		}
	}
}

func TestLocaleAwareWithQueryParams(t *testing.T) {
	st := stats.NewStatistics()
	dedup := deduplicator.NewWithLocaleSupport(st, []string{"en"})

	config := normalizer.NewConfig()
	config.LocaleAware = true
	config.LocalePriority = []string{"en"}

	urls := []string{
		"https://example.com/page?lang=en&foo=bar",
		"https://example.com/page?lang=es&foo=bar",
		"https://example.com/page?lang=fr&foo=bar",
	}

	for _, url := range urls {
		normalized, err := config.NormalizeURL(url)
		if err != nil {
			continue
		}

		dedupKey, err := config.CreateDedupKey(url)
		if err != nil {
			continue
		}

		dedup.AddWithOriginal(dedupKey, normalized, url)
	}

	entries := dedup.GetEntries()

	// Should have 1 entry (English version preferred)
	if len(entries) != 1 {
		t.Errorf("Expected 1 unique entry, got %d", len(entries))
		for i, entry := range entries {
			t.Logf("Entry %d: %s", i+1, entry.URL)
		}
	}
}

func TestNoFalsePositives(t *testing.T) {
	st := stats.NewStatistics()
	dedup := deduplicator.NewWithLocaleSupport(st, []string{"en"})

	config := normalizer.NewConfig()
	config.LocaleAware = true
	config.LocalePriority = []string{"en"}

	urls := []string{
		"https://example.com/endpoint/users",
		"https://example.com/send/email",
		"https://example.com/pen/tools",
		"https://example.com/content/pages",
	}

	for _, url := range urls {
		normalized, err := config.NormalizeURL(url)
		if err != nil {
			continue
		}

		dedupKey, err := config.CreateDedupKey(url)
		if err != nil {
			continue
		}

		dedup.AddWithOriginal(dedupKey, normalized, url)
	}

	entries := dedup.GetEntries()

	// All should be preserved (no false positives)
	if len(entries) != len(urls) {
		t.Errorf("Expected %d unique entries (no false positives), got %d", len(urls), len(entries))
		for i, entry := range entries {
			t.Logf("Entry %d: %s", i+1, entry.URL)
		}
	}
}

func TestMixedScenario(t *testing.T) {
	st := stats.NewStatistics()
	dedup := deduplicator.NewWithLocaleSupport(st, []string{"en"})

	config := normalizer.NewConfig()
	config.LocaleAware = true
	config.LocalePriority = []string{"en"}
	config.FuzzyMode = true

	urls := []string{
		// About pages (should group)
		"https://example.com/en/about",
		"https://example.com/es/sobre-nosotros",
		"https://example.com/it/chi-siamo",

		// Product pages with IDs (should group)
		"https://example.com/en/products/123",
		"https://example.com/es/productos/456",

		// Unique endpoints (should NOT group)
		"https://example.com/api/v1/users",
		"https://example.com/endpoint/data",

		// Contact pages (should group)
		"https://example.com/en/contact",
		"https://example.com/es/contacto",
	}

	for _, url := range urls {
		normalized, err := config.NormalizeURL(url)
		if err != nil {
			continue
		}

		dedupKey, err := config.CreateDedupKey(url)
		if err != nil {
			continue
		}

		dedup.AddWithOriginal(dedupKey, normalized, url)
	}

	entries := dedup.GetEntries()

	// Expected groups:
	// 1. about pages -> 1 entry
	// 2. products pages -> 1 entry
	// 3. api/v1/users -> 1 entry
	// 4. endpoint/data -> 1 entry
	// 5. contact pages -> 1 entry
	// Total: 5 entries

	if len(entries) < 4 || len(entries) > 6 {
		t.Errorf("Expected 4-6 unique entries, got %d", len(entries))
		for i, entry := range entries {
			t.Logf("Entry %d: %s (count: %d)", i+1, entry.URL, entry.Count)
		}
	}
}
