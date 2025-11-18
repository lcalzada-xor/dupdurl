package locale

import (
	"testing"
)

func TestTranslationMatcher(t *testing.T) {
	matcher := NewTranslationMatcher()

	tests := []struct {
		name     string
		seg1     string
		seg2     string
		expected bool
	}{
		{
			name:     "English and Spanish about",
			seg1:     "about",
			seg2:     "sobre-nosotros",
			expected: true,
		},
		{
			name:     "English and Italian about",
			seg1:     "about",
			seg2:     "chi-siamo",
			expected: true,
		},
		{
			name:     "Products in English and Spanish",
			seg1:     "products",
			seg2:     "productos",
			expected: true,
		},
		{
			name:     "Different words",
			seg1:     "about",
			seg2:     "contact",
			expected: false,
		},
		{
			name:     "Same word",
			seg1:     "about",
			seg2:     "about",
			expected: true,
		},
		{
			name:     "Contact in multiple languages",
			seg1:     "contact",
			seg2:     "contacto",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.AreTranslations(tt.seg1, tt.seg2)
			if result != tt.expected {
				t.Errorf("AreTranslations(%q, %q) = %v, expected %v",
					tt.seg1, tt.seg2, result, tt.expected)
			}
		})
	}
}

func TestGetCanonical(t *testing.T) {
	matcher := NewTranslationMatcher()

	tests := []struct {
		name     string
		segment  string
		expected string
	}{
		{
			name:     "Spanish about",
			segment:  "sobre-nosotros",
			expected: "about",
		},
		{
			name:     "English about",
			segment:  "about",
			expected: "about",
		},
		{
			name:     "Spanish products",
			segment:  "productos",
			expected: "product",
		},
		{
			name:     "Unknown word",
			segment:  "unknown-word",
			expected: "unknown-word",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.GetCanonical(tt.segment)
			if result != tt.expected {
				t.Errorf("GetCanonical(%q) = %q, expected %q",
					tt.segment, result, tt.expected)
			}
		})
	}
}
