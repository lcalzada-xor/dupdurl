package normalizer

import (
	"regexp"
	"strings"
)

// Path normalization constants and patterns
const (
	maxPathSegmentLength = 1024
)

var (
	// Numeric ID pattern - matches pure numeric path segments
	numericIDRegex = regexp.MustCompile(`/\d+(/|$)`)

	// UUID pattern - matches UUIDs in various formats
	uuidRegex = regexp.MustCompile(`/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}(/|$)`)

	// Hash pattern - matches MD5/SHA1-like hashes
	hashRegex = regexp.MustCompile(`/[0-9a-f]{32,40}(/|$)`)

	// Alphanumeric token pattern - matches long alphanumeric strings
	tokenRegex = regexp.MustCompile(`/[a-zA-Z0-9]{16,}(/|$)`)
)

// FuzzyPattern represents a pattern for fuzzy matching
type FuzzyPattern struct {
	Name        string
	Regex       *regexp.Regexp
	Placeholder string
	Enabled     bool
}

// GetDefaultPatterns returns the default fuzzy matching patterns
func GetDefaultPatterns() []FuzzyPattern {
	return []FuzzyPattern{
		{Name: "numeric", Regex: numericIDRegex, Placeholder: "{id}", Enabled: true},
		{Name: "uuid", Regex: uuidRegex, Placeholder: "{uuid}", Enabled: false},
		{Name: "hash", Regex: hashRegex, Placeholder: "{hash}", Enabled: false},
		{Name: "token", Regex: tokenRegex, Placeholder: "{token}", Enabled: false},
	}
}

// NormalizePath normalizes a URL path
func NormalizePath(p string) string {
	if p == "" {
		return "/"
	}

	// Collapse multiple slashes
	p = collapseSlashes(p)

	// Remove trailing slash (except root)
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}

	return p
}

// collapseSlashes removes consecutive slashes from path
func collapseSlashes(p string) string {
	if p == "" {
		return "/"
	}

	parts := strings.Split(p, "/")
	out := make([]string, 0, len(parts))

	for _, seg := range parts {
		if seg == "" {
			if len(out) == 0 {
				out = append(out, "")
			}
			continue
		}
		out = append(out, seg)
	}

	res := strings.Join(out, "/")
	if !strings.HasPrefix(res, "/") {
		res = "/" + res
	}
	return res
}

// ApplyFuzzyPatterns applies fuzzy matching patterns to a path
func ApplyFuzzyPatterns(p string, patterns []FuzzyPattern) string {
	result := p
	for _, pattern := range patterns {
		if pattern.Enabled {
			result = pattern.Regex.ReplaceAllString(result, "/"+pattern.Placeholder+"$1")
		}
	}
	return result
}

// FuzzyPath replaces numeric path segments with {id}
// This is the legacy method for backward compatibility
func FuzzyPath(p string) string {
	return numericIDRegex.ReplaceAllString(p, "/{id}$1")
}

// EnablePattern enables a fuzzy pattern by name
func EnablePattern(patterns []FuzzyPattern, name string) {
	for i := range patterns {
		if patterns[i].Name == name {
			patterns[i].Enabled = true
			return
		}
	}
}

// EnablePatterns enables multiple fuzzy patterns by name
func EnablePatterns(patterns []FuzzyPattern, names []string) {
	for _, name := range names {
		EnablePattern(patterns, name)
	}
}
