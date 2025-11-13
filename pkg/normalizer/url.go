package normalizer

import (
	"fmt"
	"net/url"
	"strings"
)

// Config holds URL normalization configuration
type Config struct {
	Mode             string
	IgnoreParams     map[string]struct{}
	SortParams       bool
	IgnoreFragment   bool
	CaseSensitive    bool
	KeepWWW          bool
	KeepScheme       bool
	TrimSpaces       bool
	FuzzyMode        bool
	FuzzyPatterns    []FuzzyPattern
	PathIncludeQuery bool
	AllowDomains     map[string]struct{}
	BlockDomains     map[string]struct{}
	IgnoreExtensions map[string]struct{}
	FilterExtensions map[string]struct{}
}

// NewConfig creates a default normalization configuration
func NewConfig() *Config {
	return &Config{
		Mode:           "url",
		IgnoreParams:   make(map[string]struct{}),
		IgnoreFragment: true,
		TrimSpaces:     true,
		FuzzyPatterns:  GetDefaultPatterns(),
		AllowDomains:   make(map[string]struct{}),
		BlockDomains:   make(map[string]struct{}),
		IgnoreExtensions: make(map[string]struct{}),
		FilterExtensions: make(map[string]struct{}),
	}
}

// NormalizeURL normalizes a URL according to the configuration
func (c *Config) NormalizeURL(raw string) (string, error) {
	if c.TrimSpaces {
		raw = strings.TrimSpace(raw)
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	// Check domain filtering
	if err := c.checkDomainFilters(u.Host); err != nil {
		return "", err
	}

	// Check extension filtering
	if err := c.checkExtensionFilter(u.Path); err != nil {
		return "", err
	}

	// Normalize scheme
	c.normalizeScheme(u)

	// Normalize host
	c.normalizeHost(u)

	// Remove fragment
	if c.IgnoreFragment {
		u.Fragment = ""
	}

	// Normalize path
	u.Path = NormalizePath(u.Path)

	// Apply fuzzy mode
	if c.FuzzyMode {
		if len(c.FuzzyPatterns) > 0 {
			u.Path = ApplyFuzzyPatterns(u.Path, c.FuzzyPatterns)
		} else {
			u.Path = FuzzyPath(u.Path)
		}
	}

	// Query params handling - keep values by default
	q := u.Query()

	// Delete ignored params
	for p := range c.IgnoreParams {
		q.Del(p)
	}

	if c.SortParams {
		u.RawQuery = BuildSortedQuery(q)
	} else {
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

// CreateDedupKey creates a key for deduplication (parameter names only, no values)
func (c *Config) CreateDedupKey(raw string) (string, error) {
	if c.TrimSpaces {
		raw = strings.TrimSpace(raw)
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	// Apply same normalization
	c.normalizeScheme(u)
	c.normalizeHost(u)

	if c.IgnoreFragment {
		u.Fragment = ""
	}

	u.Path = NormalizePath(u.Path)

	if c.FuzzyMode {
		if len(c.FuzzyPatterns) > 0 {
			u.Path = ApplyFuzzyPatterns(u.Path, c.FuzzyPatterns)
		} else {
			u.Path = FuzzyPath(u.Path)
		}
	}

	// For the dedup key, we only keep parameter NAMES, not values
	q := u.Query()

	// Delete ignored params
	for p := range c.IgnoreParams {
		q.Del(p)
	}

	// Build query string with param names only (no values)
	if len(q) > 0 {
		u.RawQuery = BuildKeyOnlyQuery(q)
	} else {
		u.RawQuery = ""
	}

	return u.String(), nil
}

// NormalizeLine normalizes a line according to the mode
func (c *Config) NormalizeLine(line string) (string, error) {
	if c.TrimSpaces {
		line = strings.TrimSpace(line)
	}
	if line == "" {
		return "", fmt.Errorf("empty line")
	}

	switch c.Mode {
	case "raw":
		if !c.CaseSensitive {
			return strings.ToLower(line), nil
		}
		return line, nil

	case "host":
		return c.extractHost(line)

	case "path":
		return c.extractPath(line)

	case "params":
		return ExtractParams(line)

	case "url":
		return c.NormalizeURL(line)

	default:
		return "", fmt.Errorf("unknown mode: %s", c.Mode)
	}
}

// Helper methods

func (c *Config) normalizeScheme(u *url.URL) {
	if !c.CaseSensitive && !c.KeepScheme {
		u.Scheme = strings.ToLower(u.Scheme)
	} else if !c.KeepScheme {
		u.Scheme = "https"
	}
}

func (c *Config) normalizeHost(u *url.URL) {
	// Normalize case FIRST
	if !c.CaseSensitive {
		u.Host = strings.ToLower(u.Host)
	}

	// Remove default ports
	if u.Scheme == "https" && strings.HasSuffix(u.Host, ":443") {
		u.Host = strings.TrimSuffix(u.Host, ":443")
	} else if u.Scheme == "http" && strings.HasSuffix(u.Host, ":80") {
		u.Host = strings.TrimSuffix(u.Host, ":80")
	}

	// Remove www (after lowercasing)
	if !c.KeepWWW && strings.HasPrefix(u.Host, "www.") {
		u.Host = strings.TrimPrefix(u.Host, "www.")
	}
}

func (c *Config) checkDomainFilters(host string) error {
	normalizedHost := strings.ToLower(host)
	if strings.HasPrefix(normalizedHost, "www.") {
		normalizedHost = strings.TrimPrefix(normalizedHost, "www.")
	}

	if len(c.AllowDomains) > 0 {
		if _, ok := c.AllowDomains[normalizedHost]; !ok {
			return fmt.Errorf("domain not in whitelist: %s", host)
		}
	}

	if len(c.BlockDomains) > 0 {
		if _, ok := c.BlockDomains[normalizedHost]; ok {
			return fmt.Errorf("domain in blacklist: %s", host)
		}
	}

	return nil
}

func (c *Config) checkExtensionFilter(path string) error {
	// Find the last dot in the path
	lastDot := strings.LastIndex(path, ".")
	if lastDot == -1 || lastDot == len(path)-1 {
		// No extension found
		if len(c.FilterExtensions) > 0 {
			// Whitelist mode: reject URLs without extension
			return fmt.Errorf("no extension found (whitelist mode requires extension)")
		}
		return nil
	}

	// Extract extension (without the dot)
	ext := strings.ToLower(path[lastDot+1:])

	// Check if there's a slash after the dot (not a real extension)
	if strings.Contains(ext, "/") {
		if len(c.FilterExtensions) > 0 {
			return fmt.Errorf("no valid extension found (whitelist mode requires extension)")
		}
		return nil
	}

	// Whitelist mode (FilterExtensions): only allow specified extensions
	if len(c.FilterExtensions) > 0 {
		if _, allowed := c.FilterExtensions[ext]; !allowed {
			return fmt.Errorf("extension not in whitelist: .%s", ext)
		}
		return nil
	}

	// Blacklist mode (IgnoreExtensions): ignore specified extensions
	if len(c.IgnoreExtensions) > 0 {
		if _, ignored := c.IgnoreExtensions[ext]; ignored {
			return fmt.Errorf("ignored extension: .%s", ext)
		}
	}

	return nil
}

func (c *Config) extractHost(line string) (string, error) {
	u, err := url.Parse(line)
	if err != nil {
		if !c.CaseSensitive {
			return strings.ToLower(line), nil
		}
		return line, nil
	}

	// Check extension filtering BEFORE processing
	if err := c.checkExtensionFilter(u.Path); err != nil {
		return "", err
	}

	h := u.Host

	// Normalize case FIRST
	if !c.CaseSensitive {
		h = strings.ToLower(h)
	}

	// Remove default ports
	if u.Scheme == "https" && strings.HasSuffix(h, ":443") {
		h = strings.TrimSuffix(h, ":443")
	} else if u.Scheme == "http" && strings.HasSuffix(h, ":80") {
		h = strings.TrimSuffix(h, ":80")
	}

	// Remove www (after lowercasing)
	if !c.KeepWWW && strings.HasPrefix(h, "www.") {
		h = strings.TrimPrefix(h, "www.")
	}

	return h, nil
}

func (c *Config) extractPath(line string) (string, error) {
	u, err := url.Parse(line)
	if err != nil {
		if !c.CaseSensitive {
			return strings.ToLower(line), nil
		}
		return line, nil
	}

	// Check extension filtering BEFORE processing
	if err := c.checkExtensionFilter(u.Path); err != nil {
		return "", err
	}

	host := u.Host

	// Normalize case FIRST
	if !c.CaseSensitive {
		host = strings.ToLower(host)
	}

	// Remove default ports
	if u.Scheme == "https" && strings.HasSuffix(host, ":443") {
		host = strings.TrimSuffix(host, ":443")
	} else if u.Scheme == "http" && strings.HasSuffix(host, ":80") {
		host = strings.TrimSuffix(host, ":80")
	}

	// Remove www (after lowercasing)
	if !c.KeepWWW && strings.HasPrefix(host, "www.") {
		host = strings.TrimPrefix(host, "www.")
	}

	path := NormalizePath(u.Path)
	if !c.CaseSensitive {
		path = strings.ToLower(path)
	}
	if c.FuzzyMode {
		if len(c.FuzzyPatterns) > 0 {
			path = ApplyFuzzyPatterns(path, c.FuzzyPatterns)
		} else {
			path = FuzzyPath(path)
		}
	}

	result := host + path

	// Optionally include normalized query
	if c.PathIncludeQuery && u.RawQuery != "" {
		q := u.Query()
		for p := range c.IgnoreParams {
			q.Del(p)
		}
		if c.SortParams {
			result += "?" + BuildSortedQuery(q)
		} else {
			result += "?" + q.Encode()
		}
	}

	return result, nil
}
