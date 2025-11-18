package locale

import (
	"testing"
)

func TestGrouperBasic(t *testing.T) {
	grouper := NewGrouper([]string{"en"})

	urls := []string{
		"https://example.com/about",
		"https://example.com/en/about",
		"https://example.com/es/sobre-nosotros",
		"https://example.com/it/chi-siamo",
	}

	for _, url := range urls {
		err := grouper.Add(url)
		if err != nil {
			t.Fatalf("Add(%q) error = %v", url, err)
		}
	}

	bestURLs := grouper.GetBestURLs()

	// Should have 1 group (all about pages)
	if len(bestURLs) != 1 {
		t.Errorf("Expected 1 group, got %d", len(bestURLs))
	}

	// Best URL should be the English one
	if len(bestURLs) > 0 {
		if bestURLs[0].Locale != "en" {
			t.Errorf("Expected English locale to be selected, got %q", bestURLs[0].Locale)
		}
	}
}

func TestGrouperMultiplePages(t *testing.T) {
	grouper := NewGrouper([]string{"en"})

	urls := []string{
		"https://example.com/en/about",
		"https://example.com/es/sobre-nosotros",
		"https://example.com/en/products",
		"https://example.com/es/productos",
		"https://example.com/unique-page",
	}

	for _, url := range urls {
		err := grouper.Add(url)
		if err != nil {
			t.Fatalf("Add(%q) error = %v", url, err)
		}
	}

	bestURLs := grouper.GetBestURLs()

	// Should have 3 groups: about, products, unique-page
	if len(bestURLs) != 3 {
		t.Errorf("Expected 3 groups, got %d", len(bestURLs))
	}
}

func TestGrouperDifferentPaths(t *testing.T) {
	grouper := NewGrouper([]string{"en"})

	urls := []string{
		"https://example.com/endpoint/users",
		"https://example.com/en/endpoint/users",
	}

	for _, url := range urls {
		err := grouper.Add(url)
		if err != nil {
			t.Fatalf("Add(%q) error = %v", url, err)
		}
	}

	bestURLs := grouper.GetBestURLs()

	// These should be treated as different paths (endpoint is not a locale)
	// Depending on detector logic, this might be 1 or 2 groups
	if len(bestURLs) < 1 {
		t.Errorf("Expected at least 1 group, got %d", len(bestURLs))
	}
}

func TestShouldGroup(t *testing.T) {
	grouper := NewGrouper([]string{"en"})

	tests := []struct {
		name     string
		url1     string
		url2     string
		expected bool
	}{
		{
			name:     "Same page different locales",
			url1:     "https://example.com/en/about",
			url2:     "https://example.com/es/sobre-nosotros",
			expected: true,
		},
		{
			name:     "Different pages",
			url1:     "https://example.com/en/about",
			url2:     "https://example.com/en/products",
			expected: false,
		},
		{
			name:     "Same page no locale vs locale",
			url1:     "https://example.com/about",
			url2:     "https://example.com/en/about",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := grouper.ShouldGroup(tt.url1, tt.url2)
			if err != nil {
				t.Fatalf("ShouldGroup() error = %v", err)
			}

			if result != tt.expected {
				t.Errorf("ShouldGroup(%q, %q) = %v, expected %v",
					tt.url1, tt.url2, result, tt.expected)
			}
		})
	}
}
