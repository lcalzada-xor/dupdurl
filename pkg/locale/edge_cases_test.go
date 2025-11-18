package locale

import (
	"testing"
)

// TestEdgeCases tests various edge cases and corner scenarios
func TestEdgeCases(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name           string
		url            string
		expectedLocale string
		expectedType   LocaleType
		shouldError    bool
	}{
		// Multiple locale indicators
		{
			name:           "Multiple locales in path",
			url:            "https://example.com/en/fr/about",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},
		{
			name:           "Subdomain and path locale",
			url:            "https://en.example.com/es/about",
			expectedLocale: "en",
			expectedType:   LocaleTypeSubdomain,
		},
		{
			name:           "All three types present",
			url:            "https://en.example.com/es/about?lang=fr",
			expectedLocale: "en",
			expectedType:   LocaleTypeSubdomain,
		},

		// Extended locale codes
		{
			name:           "Extended locale en-US",
			url:            "https://example.com/en-us/about",
			expectedLocale: "en-us",
			expectedType:   LocaleTypePath,
		},
		{
			name:           "Extended locale pt-BR",
			url:            "https://example.com/pt-br/produtos",
			expectedLocale: "pt-br",
			expectedType:   LocaleTypePath,
		},

		// Edge cases with numbers
		{
			name:           "Locale with number in path",
			url:            "https://example.com/en/page/123",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},

		// Deep paths
		{
			name:           "Deep path with locale",
			url:            "https://example.com/en/category/subcategory/product/123",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},

		// Special characters
		{
			name:           "Path with dashes and locale",
			url:            "https://example.com/en/sobre-nosotros-empresa",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},
		{
			name:           "Path with underscores",
			url:            "https://example.com/en/about_us",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},

		// Root paths
		{
			name:           "Locale at root",
			url:            "https://example.com/en",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},
		{
			name:           "Locale at root with slash",
			url:            "https://example.com/en/",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},

		// Empty and minimal cases
		{
			name:           "Root URL only",
			url:            "https://example.com",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},
		{
			name:           "Root URL with slash",
			url:            "https://example.com/",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},

		// Query parameters edge cases
		{
			name:           "Multiple lang params",
			url:            "https://example.com/page?lang=en&locale=es",
			expectedLocale: "en",
			expectedType:   LocaleTypeQuery,
		},
		{
			name:           "Case variation in param",
			url:            "https://example.com/page?LANG=en",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},

		// Complex false positives
		{
			name:           "Word containing locale code",
			url:            "https://example.com/broken/page",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},
		{
			name:           "Identifier",
			url:            "https://example.com/item/id/12345",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},
		{
			name:           "Send endpoint",
			url:            "https://example.com/send/notification",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},

		// API endpoints with locale
		{
			name:           "API v1 with locale",
			url:            "https://example.com/api/v1/en/users",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},
		{
			name:           "GraphQL endpoint",
			url:            "https://example.com/graphql/en",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},

		// Unusual but valid locales
		{
			name:           "Chinese locale",
			url:            "https://example.com/zh/about",
			expectedLocale: "zh",
			expectedType:   LocaleTypePath,
		},
		{
			name:           "Japanese locale",
			url:            "https://example.com/ja/about",
			expectedLocale: "ja",
			expectedType:   LocaleTypePath,
		},
		{
			name:           "Korean locale",
			url:            "https://example.com/ko/about",
			expectedLocale: "ko",
			expectedType:   LocaleTypePath,
		},
		{
			name:           "Arabic locale",
			url:            "https://example.com/ar/about",
			expectedLocale: "ar",
			expectedType:   LocaleTypePath,
		},

		// Port numbers
		{
			name:           "URL with port and locale",
			url:            "https://example.com:8080/en/about",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},

		// Fragments
		{
			name:           "URL with fragment",
			url:            "https://example.com/en/about#section",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(tt.url)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.Locale != tt.expectedLocale {
				t.Errorf("Expected locale %q, got %q", tt.expectedLocale, result.Locale)
			}

			if result.LocaleType != tt.expectedType {
				t.Errorf("Expected type %q, got %q", tt.expectedType, result.LocaleType)
			}
		})
	}
}

// TestBaseURLGeneration tests that base URLs are correctly generated
func TestBaseURLGeneration(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Simple path locale removal",
			url:      "https://example.com/en/about",
			expected: "https://example.com/about",
		},
		{
			name:     "Subdomain locale removal",
			url:      "https://en.example.com/about",
			expected: "https://example.com/about",
		},
		{
			name:     "Query param locale removal",
			url:      "https://example.com/about?lang=en&foo=bar",
			expected: "https://example.com/about?foo=bar",
		},
		{
			name:     "Deep path locale removal",
			url:      "https://example.com/en/category/product",
			expected: "https://example.com/category/product",
		},
		{
			name:     "Extended locale removal",
			url:      "https://example.com/en-US/about",
			expected: "https://example.com/about",
		},
		{
			name:     "No locale - URL unchanged",
			url:      "https://example.com/about",
			expected: "https://example.com/about",
		},
		{
			name:     "Preserve query params",
			url:      "https://example.com/en/search?q=test&page=1",
			expected: "https://example.com/search?q=test&page=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(tt.url)
			if err != nil {
				t.Fatalf("Detect error: %v", err)
			}

			if result.BaseURL != tt.expected {
				t.Errorf("Expected base URL %q, got %q", tt.expected, result.BaseURL)
			}
		})
	}
}

// TestConcurrentAccess tests thread safety
func TestConcurrentAccess(t *testing.T) {
	detector := NewDetector()
	grouper := NewGrouper([]string{"en"})

	urls := []string{
		"https://example.com/en/about",
		"https://example.com/es/sobre-nosotros",
		"https://example.com/it/chi-siamo",
		"https://example.com/fr/a-propos",
	}

	// Run concurrent detections
	done := make(chan bool, len(urls))

	for _, url := range urls {
		go func(u string) {
			_, err := detector.Detect(u)
			if err != nil {
				t.Errorf("Concurrent detect error: %v", err)
			}

			err = grouper.Add(u)
			if err != nil {
				t.Errorf("Concurrent add error: %v", err)
			}

			done <- true
		}(url)
	}

	// Wait for all goroutines
	for i := 0; i < len(urls); i++ {
		<-done
	}

	// Verify results
	bestURLs := grouper.GetBestURLs()
	if len(bestURLs) != 1 {
		t.Errorf("Expected 1 group, got %d", len(bestURLs))
	}
}

// TestMalformedURLs tests handling of malformed URLs
func TestMalformedURLs(t *testing.T) {
	detector := NewDetector()

	malformed := []string{
		"not-a-url",
		"htp://missing-t.com",
		"://no-scheme.com",
		"",
		"   ",
	}

	for _, url := range malformed {
		t.Run(url, func(t *testing.T) {
			_, err := detector.Detect(url)
			// Should handle gracefully (either error or process as-is)
			if err != nil {
				// Expected behavior - error on malformed URL
				t.Logf("Correctly errored on malformed URL: %v", err)
			}
		})
	}
}

// TestLargeScale tests with many URLs
func TestLargeScale(t *testing.T) {
	grouper := NewGrouper([]string{"en"})

	// Generate 1000 URLs with different locales
	locales := []string{"en", "es", "fr", "de", "it", "pt", "ja", "zh", "ko", "ar"}
	paths := []string{"about", "products", "contact", "services", "help"}

	count := 0
	for _, locale := range locales {
		for _, path := range paths {
			url := "https://example.com/" + locale + "/" + path
			err := grouper.Add(url)
			if err != nil {
				t.Fatalf("Error adding URL %s: %v", url, err)
			}
			count++
		}
	}

	t.Logf("Added %d URLs", count)

	bestURLs := grouper.GetBestURLs()
	t.Logf("Got %d unique endpoints", len(bestURLs))

	// Should have 5 groups (one per path)
	if len(bestURLs) != 5 {
		t.Errorf("Expected 5 groups, got %d", len(bestURLs))
	}

	// All should be English
	for _, locURL := range bestURLs {
		if locURL.Locale != "en" {
			t.Errorf("Expected English locale, got %q for URL %s", locURL.Locale, locURL.OriginalURL)
		}
	}
}
