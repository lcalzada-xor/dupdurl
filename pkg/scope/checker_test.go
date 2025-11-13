package scope

import (
	"testing"
)

func TestScopeChecker_IsInScope(t *testing.T) {
	tests := []struct {
		name     string
		includes []string
		excludes []string
		host     string
		expected bool
	}{
		{
			name:     "exact match",
			includes: []string{"example.com"},
			host:     "example.com",
			expected: true,
		},
		{
			name:     "wildcard subdomain match",
			includes: []string{"*.example.com"},
			host:     "api.example.com",
			expected: true,
		},
		{
			name:     "wildcard with base domain",
			includes: []string{"example.com", "*.example.com"},
			host:     "api.example.com",
			expected: true,
		},
		{
			name:     "excluded subdomain",
			includes: []string{"*.example.com"},
			excludes: []string{"dev.example.com"},
			host:     "dev.example.com",
			expected: false,
		},
		{
			name:     "out of scope",
			includes: []string{"example.com"},
			host:     "attacker.com",
			expected: false,
		},
		{
			name:     "www removal",
			includes: []string{"example.com"},
			host:     "www.example.com",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewChecker()
			for _, inc := range tt.includes {
				checker.AddInclude(inc)
			}
			for _, exc := range tt.excludes {
				checker.AddExclude(exc)
			}

			got := checker.IsInScope(tt.host)
			if got != tt.expected {
				t.Errorf("IsInScope(%q) = %v; want %v", tt.host, got, tt.expected)
			}
		})
	}
}
