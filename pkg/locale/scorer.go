package locale

import (
	"net/url"
	"strings"
)

// Score represents a URL's priority score
type Score struct {
	URL            string
	LocaleScore    int // Higher priority locales get higher scores
	CompletenessScore int // URLs with more info (query params) score higher
	FirstSeenBonus int // First seen URLs get a bonus
	TotalScore     int
}

// Scorer handles URL scoring for prioritization
type Scorer struct {
	localePriority map[string]int // locale -> priority score
}

// NewScorer creates a new scorer with given locale priorities
func NewScorer(priorities []string) *Scorer {
	s := &Scorer{
		localePriority: make(map[string]int),
	}

	// Assign scores based on priority order (higher index = higher priority)
	// Default locale gets middle score
	s.localePriority["default"] = 50

	// Priority locales get incrementing scores
	for i, locale := range priorities {
		s.localePriority[locale] = 100 + (len(priorities)-i)*10
	}

	return s
}

// Score calculates the score for a localized URL
func (s *Scorer) Score(localized *LocalizedURL, isFirstSeen bool) Score {
	score := Score{
		URL: localized.OriginalURL,
	}

	// Locale score
	locale := localized.Locale
	if locale == "" {
		locale = "default"
	}

	if priorityScore, exists := s.localePriority[locale]; exists {
		score.LocaleScore = priorityScore
	} else {
		// Unknown locale gets low score
		score.LocaleScore = 25
	}

	// Completeness score
	score.CompletenessScore = s.calculateCompleteness(localized.OriginalURL)

	// First seen bonus
	if isFirstSeen {
		score.FirstSeenBonus = 10
	}

	// Calculate total
	score.TotalScore = score.LocaleScore + score.CompletenessScore + score.FirstSeenBonus

	return score
}

// calculateCompleteness scores based on URL completeness
func (s *Scorer) calculateCompleteness(rawURL string) int {
	u, err := url.Parse(rawURL)
	if err != nil {
		return 0
	}

	score := 0

	// Query parameters add to completeness
	params := u.Query()
	score += len(params) * 2

	// Path depth adds to completeness
	pathSegments := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(pathSegments) > 0 && pathSegments[0] != "" {
		score += len(pathSegments)
	}

	// Cap at 20 to avoid overwhelming locale score
	if score > 20 {
		score = 20
	}

	return score
}

// ComparePriority compares two locales and returns which has higher priority
// Returns: -1 if locale1 > locale2, 0 if equal, 1 if locale2 > locale1
func (s *Scorer) ComparePriority(locale1, locale2 string) int {
	if locale1 == "" {
		locale1 = "default"
	}
	if locale2 == "" {
		locale2 = "default"
	}

	score1 := s.localePriority[locale1]
	score2 := s.localePriority[locale2]

	if score1 > score2 {
		return -1
	} else if score1 < score2 {
		return 1
	}
	return 0
}

// GetBestFromGroup selects the best URL from a group
func (s *Scorer) GetBestFromGroup(group *LocaleGroup) *LocalizedURL {
	if len(group.URLs) == 0 {
		return nil
	}

	var bestURL *LocalizedURL
	var bestScore Score
	isFirst := true

	for _, locURL := range group.URLs {
		score := s.Score(locURL, isFirst)
		isFirst = false

		if bestURL == nil || score.TotalScore > bestScore.TotalScore {
			bestURL = locURL
			bestScore = score
		}
	}

	return bestURL
}
