package locale

import (
	"net/url"
	"strings"
)

// LocaleGroup represents a group of URLs that are translations of each other
type LocaleGroup struct {
	BaseKey     string                    // Normalized base key for grouping
	URLs        map[string]*LocalizedURL  // locale -> LocalizedURL
	BestURL     *LocalizedURL             // The selected "best" URL
	Priority    []string                  // Priority order for locale selection
}

// Grouper handles grouping of localized URLs
type Grouper struct {
	detector           *Detector
	translationMatcher *TranslationMatcher
	groups             map[string]*LocaleGroup
	Priority           []string // Exported for access
}

// NewGrouper creates a new locale grouper
func NewGrouper(priority []string) *Grouper {
	if len(priority) == 0 {
		priority = []string{"en"} // Default priority: English
	}

	return &Grouper{
		detector:           NewDetector(),
		translationMatcher: NewTranslationMatcher(),
		groups:             make(map[string]*LocaleGroup),
		Priority:           priority,
	}
}

// Add adds a URL to the grouper
func (g *Grouper) Add(rawURL string) error {
	localized, err := g.detector.Detect(rawURL)
	if err != nil {
		return err
	}

	// Generate a grouping key
	groupKey := g.generateGroupKey(localized)

	// Get or create group
	group, exists := g.groups[groupKey]
	if !exists {
		group = &LocaleGroup{
			BaseKey:  groupKey,
			URLs:     make(map[string]*LocalizedURL),
			Priority: g.Priority,
		}
		g.groups[groupKey] = group
	}

	// Add URL to group
	locale := localized.Locale
	if locale == "" {
		locale = "default"
	}

	// Only keep first occurrence of each locale
	if _, exists := group.URLs[locale]; !exists {
		group.URLs[locale] = localized
	}

	// Update best URL
	g.updateBestURL(group)

	return nil
}

// generateGroupKey creates a unique key for grouping similar URLs
func (g *Grouper) generateGroupKey(localized *LocalizedURL) string {
	baseURL := localized.BaseURL

	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	// Normalize host
	host := strings.ToLower(u.Host)
	host = strings.TrimPrefix(host, "www.")

	// Normalize path with translation awareness
	path := g.normalizePath(u.Path)

	// Build key: host + normalized path
	key := host + path

	// Include sorted query parameter names (not values)
	if u.RawQuery != "" {
		params := u.Query()
		paramNames := make([]string, 0, len(params))
		for name := range params {
			paramNames = append(paramNames, strings.ToLower(name))
		}
		// Sort for consistency
		if len(paramNames) > 0 {
			key += "?" + strings.Join(sortStrings(paramNames), "&")
		}
	}

	return key
}

// normalizePath normalizes a path with translation awareness
func (g *Grouper) normalizePath(path string) string {
	if path == "" || path == "/" {
		return "/"
	}

	segments := strings.Split(strings.Trim(path, "/"), "/")
	normalized := make([]string, len(segments))

	for i, seg := range segments {
		// Convert to lowercase
		segLower := strings.ToLower(seg)

		// Check if it's a known translation
		canonical := g.translationMatcher.GetCanonical(segLower)
		normalized[i] = canonical
	}

	return "/" + strings.Join(normalized, "/")
}

// updateBestURL updates the best URL for a group based on priority
func (g *Grouper) updateBestURL(group *LocaleGroup) {
	// Priority-based selection
	for _, priorityLocale := range g.Priority {
		if url, exists := group.URLs[priorityLocale]; exists {
			group.BestURL = url
			return
		}
	}

	// If no priority match, use "default" (no locale detected)
	if url, exists := group.URLs["default"]; exists {
		group.BestURL = url
		return
	}

	// Otherwise, use first available
	for _, url := range group.URLs {
		group.BestURL = url
		return
	}
}

// GetBestURLs returns the best URL from each group
func (g *Grouper) GetBestURLs() []*LocalizedURL {
	result := make([]*LocalizedURL, 0, len(g.groups))
	for _, group := range g.groups {
		if group.BestURL != nil {
			result = append(result, group.BestURL)
		}
	}
	return result
}

// GetGroups returns all groups
func (g *Grouper) GetGroups() map[string]*LocaleGroup {
	return g.groups
}

// ShouldGroup determines if two URLs should be grouped together
func (g *Grouper) ShouldGroup(url1, url2 string) (bool, error) {
	loc1, err := g.detector.Detect(url1)
	if err != nil {
		return false, err
	}

	loc2, err := g.detector.Detect(url2)
	if err != nil {
		return false, err
	}

	// Generate keys and compare
	key1 := g.generateGroupKey(loc1)
	key2 := g.generateGroupKey(loc2)

	if key1 == key2 {
		// Additional validation: check if they're actually similar enough
		return g.validateSimilarity(loc1, loc2), nil
	}

	return false, nil
}

// validateSimilarity performs additional validation to avoid false positives
func (g *Grouper) validateSimilarity(loc1, loc2 *LocalizedURL) bool {
	u1, err1 := url.Parse(loc1.BaseURL)
	u2, err2 := url.Parse(loc2.BaseURL)

	if err1 != nil || err2 != nil {
		return false
	}

	// Must have same host
	host1 := strings.ToLower(strings.TrimPrefix(u1.Host, "www."))
	host2 := strings.ToLower(strings.TrimPrefix(u2.Host, "www."))
	if host1 != host2 {
		return false
	}

	// Must have same number of path segments
	seg1 := strings.Split(strings.Trim(u1.Path, "/"), "/")
	seg2 := strings.Split(strings.Trim(u2.Path, "/"), "/")

	if len(seg1) != len(seg2) {
		return false
	}

	// Check if segments are either identical or translations
	matchCount := 0
	for i := 0; i < len(seg1); i++ {
		if seg1[i] == seg2[i] {
			matchCount++
		} else if g.translationMatcher.AreTranslations(seg1[i], seg2[i]) {
			matchCount++
		}
	}

	// At least 70% of segments must match or be translations
	threshold := float64(len(seg1)) * 0.7
	return float64(matchCount) >= threshold
}

// sortStrings is a simple bubble sort for string slices
func sortStrings(strs []string) []string {
	result := make([]string, len(strs))
	copy(result, strs)

	n := len(result)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}
