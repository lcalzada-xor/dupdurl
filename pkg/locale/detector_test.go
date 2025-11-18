package locale

import (
	"testing"
)

func TestDetectPathPrefix(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name           string
		url            string
		expectedLocale string
		expectedType   LocaleType
	}{
		{
			name:           "English path prefix",
			url:            "https://example.com/en/about",
			expectedLocale: "en",
			expectedType:   LocaleTypePath,
		},
		{
			name:           "Spanish path prefix",
			url:            "https://example.com/es/productos",
			expectedLocale: "es",
			expectedType:   LocaleTypePath,
		},
		{
			name:           "Italian path prefix",
			url:            "https://example.com/it/chi-siamo",
			expectedLocale: "it",
			expectedType:   LocaleTypePath,
		},
		{
			name:           "No locale",
			url:            "https://example.com/about",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},
		{
			name:           "False positive - endpoint",
			url:            "https://example.com/endpoint/users",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},
		{
			name:           "False positive - id",
			url:            "https://example.com/id/users",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(tt.url)
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
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

func TestDetectSubdomain(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name           string
		url            string
		expectedLocale string
		expectedType   LocaleType
	}{
		{
			name:           "English subdomain",
			url:            "https://en.example.com/about",
			expectedLocale: "en",
			expectedType:   LocaleTypeSubdomain,
		},
		{
			name:           "Spanish subdomain",
			url:            "https://es.example.com/productos",
			expectedLocale: "es",
			expectedType:   LocaleTypeSubdomain,
		},
		{
			name:           "No locale subdomain",
			url:            "https://www.example.com/about",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(tt.url)
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
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

func TestDetectQueryParam(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name           string
		url            string
		expectedLocale string
		expectedType   LocaleType
	}{
		{
			name:           "Lang parameter",
			url:            "https://example.com/about?lang=en",
			expectedLocale: "en",
			expectedType:   LocaleTypeQuery,
		},
		{
			name:           "Locale parameter",
			url:            "https://example.com/about?locale=es",
			expectedLocale: "es",
			expectedType:   LocaleTypeQuery,
		},
		{
			name:           "No locale parameter",
			url:            "https://example.com/about?foo=bar",
			expectedLocale: "",
			expectedType:   LocaleTypeNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(tt.url)
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
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

func TestRemoveLocale(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name        string
		url         string
		expectedBase string
	}{
		{
			name:        "Remove path locale",
			url:         "https://example.com/en/about",
			expectedBase: "https://example.com/about",
		},
		{
			name:        "Remove subdomain locale",
			url:         "https://en.example.com/about",
			expectedBase: "https://example.com/about",
		},
		{
			name:        "Remove query locale",
			url:         "https://example.com/about?lang=en&foo=bar",
			expectedBase: "https://example.com/about?foo=bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(tt.url)
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			if result.BaseURL != tt.expectedBase {
				t.Errorf("Expected base URL %q, got %q", tt.expectedBase, result.BaseURL)
			}
		})
	}
}
