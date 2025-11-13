package unit

import (
	"testing"

	"github.com/lcalzada-xor/dupdurl/pkg/normalizer"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty path", "", "/"},
		{"root path", "/", "/"},
		{"simple path", "/api/users", "/api/users"},
		{"trailing slash", "/api/users/", "/api/users"},
		{"multiple slashes", "/api//users///profile", "/api/users/profile"},
		{"no leading slash", "api/users", "/api/users"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizer.NormalizePath(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePath(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFuzzyPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"numeric ID", "/api/users/123/profile", "/api/users/{id}/profile"},
		{"multiple IDs", "/api/123/items/456", "/api/{id}/items/{id}"},
		{"no ID", "/api/users/profile", "/api/users/profile"},
		{"trailing ID", "/api/users/123", "/api/users/{id}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizer.FuzzyPath(tt.input)
			if result != tt.expected {
				t.Errorf("FuzzyPath(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		config   *normalizer.Config
		expected string
		wantErr  bool
	}{
		{
			name:  "basic normalization",
			input: "https://www.example.com/path",
			config: &normalizer.Config{
				IgnoreFragment: true,
				TrimSpaces:     true,
				KeepWWW:        false,
				IgnoreParams:   make(map[string]struct{}),
				AllowDomains:   make(map[string]struct{}),
				BlockDomains:   make(map[string]struct{}),
				IgnoreExtensions: make(map[string]struct{}),
			},
			expected: "https://example.com/path",
			wantErr:  false,
		},
		{
			name:  "remove fragment",
			input: "https://example.com/path#section",
			config: &normalizer.Config{
				IgnoreFragment: true,
				TrimSpaces:     true,
				KeepWWW:        false,
				IgnoreParams:   make(map[string]struct{}),
				AllowDomains:   make(map[string]struct{}),
				BlockDomains:   make(map[string]struct{}),
				IgnoreExtensions: make(map[string]struct{}),
			},
			expected: "https://example.com/path",
			wantErr:  false,
		},
		{
			name:  "keep www",
			input: "https://www.example.com/path",
			config: &normalizer.Config{
				IgnoreFragment: true,
				TrimSpaces:     true,
				KeepWWW:        true,
				IgnoreParams:   make(map[string]struct{}),
				AllowDomains:   make(map[string]struct{}),
				BlockDomains:   make(map[string]struct{}),
				IgnoreExtensions: make(map[string]struct{}),
			},
			expected: "https://www.example.com/path",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.config.NormalizeURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("NormalizeURL(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildSortedQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string][]string
		expected string
	}{
		{
			name:     "empty query",
			params:   map[string][]string{},
			expected: "",
		},
		{
			name: "single param",
			params: map[string][]string{
				"foo": {"bar"},
			},
			expected: "foo=bar",
		},
		{
			name: "sorted params",
			params: map[string][]string{
				"z": {"1"},
				"a": {"2"},
				"m": {"3"},
			},
			expected: "a=2&m=3&z=1",
		},
		{
			name: "multiple values sorted",
			params: map[string][]string{
				"sort": {"name", "date", "id"},
			},
			expected: "sort=date&sort=id&sort=name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizer.BuildSortedQuery(tt.params)
			if result != tt.expected {
				t.Errorf("BuildSortedQuery() = %q; want %q", result, tt.expected)
			}
		})
	}
}

func TestParseSet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]struct{}
	}{
		{
			name:     "empty string",
			input:    "",
			expected: map[string]struct{}{},
		},
		{
			name:  "single item",
			input: "foo",
			expected: map[string]struct{}{
				"foo": {},
			},
		},
		{
			name:  "multiple items",
			input: "foo,bar,baz",
			expected: map[string]struct{}{
				"foo": {},
				"bar": {},
				"baz": {},
			},
		},
		{
			name:  "with spaces",
			input: "foo , bar , baz",
			expected: map[string]struct{}{
				"foo": {},
				"bar": {},
				"baz": {},
			},
		},
		{
			name:  "case insensitive",
			input: "FOO,Bar,BAZ",
			expected: map[string]struct{}{
				"foo": {},
				"bar": {},
				"baz": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizer.ParseSet(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseSet() length = %d; want %d", len(result), len(tt.expected))
			}
			for k := range tt.expected {
				if _, ok := result[k]; !ok {
					t.Errorf("ParseSet() missing key %q", k)
				}
			}
		})
	}
}

func TestExtensionFilter(t *testing.T) {
	config := &normalizer.Config{
		IgnoreFragment: true,
		TrimSpaces:     true,
		IgnoreParams:   make(map[string]struct{}),
		AllowDomains:   make(map[string]struct{}),
		BlockDomains:   make(map[string]struct{}),
		IgnoreExtensions: map[string]struct{}{
			"jpg": {},
			"png": {},
			"css": {},
		},
	}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"image jpg", "https://example.com/image.jpg", true},
		{"image png", "https://example.com/path/to/image.png", true},
		{"css file", "https://example.com/style.css", true},
		{"html file", "https://example.com/page.html", false},
		{"no extension", "https://example.com/api/users", false},
		{"extension in path", "https://example.com/api.v1/users", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := config.NormalizeURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeURL() with extension filter error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
