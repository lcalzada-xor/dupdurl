package scope

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Checker handles scope verification for URLs
type Checker struct {
	includes []pattern // Patterns to include
	excludes []pattern // Patterns to exclude
}

// pattern represents a scope pattern with wildcard support
type pattern struct {
	raw       string // Original pattern
	parts     []string
	hasPrefix bool // Starts with *
	hasSuffix bool // Ends with *
}

// NewChecker creates a new scope checker
func NewChecker() *Checker {
	return &Checker{
		includes: []pattern{},
		excludes: []pattern{},
	}
}

// LoadFromFile loads scope rules from a file
func (c *Checker) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open scope file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if it's an exclusion (starts with !)
		if strings.HasPrefix(line, "!") {
			pattern := strings.TrimSpace(line[1:])
			c.AddExclude(pattern)
		} else {
			c.AddInclude(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading scope file: %w", err)
	}

	return nil
}

// AddInclude adds an inclusion pattern
func (c *Checker) AddInclude(pattern string) {
	c.includes = append(c.includes, parsePattern(pattern))
}

// AddExclude adds an exclusion pattern
func (c *Checker) AddExclude(pattern string) {
	c.excludes = append(c.excludes, parsePattern(pattern))
}

// parsePattern parses a pattern with wildcard support
func parsePattern(raw string) pattern {
	p := pattern{
		raw: raw,
	}

	// Check for wildcards
	p.hasPrefix = strings.HasPrefix(raw, "*")
	p.hasSuffix = strings.HasSuffix(raw, "*")

	// Remove leading/trailing asterisks but keep internal structure
	clean := raw
	if p.hasPrefix {
		clean = clean[1:] // Remove leading *
	}
	if p.hasSuffix {
		clean = clean[:len(clean)-1] // Remove trailing *
	}

	// Remove leading dot if present (from *.example.com)
	clean = strings.TrimPrefix(clean, ".")

	if clean != "" {
		p.parts = strings.Split(clean, "*")
	}

	return p
}

// IsInScope checks if a host is in scope
func (c *Checker) IsInScope(host string) bool {
	// Normalize host (remove port if present)
	host = normalizeHost(host)

	// If no includes defined, everything is in scope by default
	if len(c.includes) == 0 {
		// But still check excludes
		for _, excl := range c.excludes {
			if matchPattern(host, excl) {
				return false
			}
		}
		return true
	}

	// Check if matches any include pattern
	inScope := false
	for _, incl := range c.includes {
		if matchPattern(host, incl) {
			inScope = true
			break
		}
	}

	// If not in includes, out of scope
	if !inScope {
		return false
	}

	// Check if matches any exclude pattern
	for _, excl := range c.excludes {
		if matchPattern(host, excl) {
			return false
		}
	}

	return true
}

// normalizeHost removes port and normalizes the host
func normalizeHost(host string) string {
	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Convert to lowercase
	host = strings.ToLower(host)

	// Remove www. prefix for comparison
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}

	return host
}

// matchPattern checks if a host matches a pattern
func matchPattern(host string, p pattern) bool {
	// Exact match (no wildcards)
	if !p.hasPrefix && !p.hasSuffix && len(p.parts) == 1 {
		return host == p.parts[0]
	}

	// Prefix wildcard: *.example.com
	if p.hasPrefix && !p.hasSuffix && len(p.parts) == 1 {
		// Match exact or subdomain
		domain := p.parts[0]
		return host == domain || strings.HasSuffix(host, "."+domain)
	}

	// Suffix wildcard: example.*
	if !p.hasPrefix && p.hasSuffix && len(p.parts) == 1 {
		return strings.HasPrefix(host, p.parts[0])
	}

	// Both wildcards: *example*
	if p.hasPrefix && p.hasSuffix && len(p.parts) == 1 {
		return strings.Contains(host, p.parts[0])
	}

	// Multiple wildcards: api*.example.*
	if len(p.parts) > 1 {
		pos := 0
		for i, part := range p.parts {
			if part == "" {
				continue
			}

			idx := strings.Index(host[pos:], part)
			if idx == -1 {
				return false
			}

			// First part with no prefix wildcard must match at start
			if i == 0 && !p.hasPrefix && idx != 0 {
				return false
			}

			pos += idx + len(part)
		}

		// Last part with no suffix wildcard must match at end
		if !p.hasSuffix && pos != len(host) {
			return false
		}

		return true
	}

	return false
}

// GetStats returns scope statistics
func (c *Checker) GetStats() ScopeStats {
	return ScopeStats{
		IncludePatterns: len(c.includes),
		ExcludePatterns: len(c.excludes),
	}
}

// ScopeStats holds scope statistics
type ScopeStats struct {
	IncludePatterns int
	ExcludePatterns int
	InScope         int
	OutOfScope      int
}

// HasRules returns true if any scope rules are defined
func (c *Checker) HasRules() bool {
	return len(c.includes) > 0 || len(c.excludes) > 0
}
