package locale

import (
	"net/url"
	"regexp"
	"strings"
)

// LocaleType represents where the locale is found in the URL
type LocaleType string

const (
	LocaleTypePath      LocaleType = "path"
	LocaleTypeSubdomain LocaleType = "subdomain"
	LocaleTypeQuery     LocaleType = "query"
	LocaleTypeNone      LocaleType = "none"
)

// LocalizedURL represents a URL with locale information
type LocalizedURL struct {
	BaseURL     string     // URL without locale component
	Locale      string     // Detected locale code (e.g., "en", "es")
	LocaleType  LocaleType // Where the locale was found
	OriginalURL string     // Original URL
	Position    int        // Position in path segments (for path type)
}

// ISO 639-1 language codes (comprehensive list)
var localeCodes = map[string]bool{
	"aa": true, "ab": true, "ae": true, "af": true, "ak": true, "am": true, "an": true, "ar": true,
	"as": true, "av": true, "ay": true, "az": true, "ba": true, "be": true, "bg": true, "bh": true,
	"bi": true, "bm": true, "bn": true, "bo": true, "br": true, "bs": true, "ca": true, "ce": true,
	"ch": true, "co": true, "cr": true, "cs": true, "cu": true, "cv": true, "cy": true, "da": true,
	"de": true, "dv": true, "dz": true, "ee": true, "el": true, "en": true, "eo": true, "es": true,
	"et": true, "eu": true, "fa": true, "ff": true, "fi": true, "fj": true, "fo": true, "fr": true,
	"fy": true, "ga": true, "gd": true, "gl": true, "gn": true, "gu": true, "gv": true, "ha": true,
	"he": true, "hi": true, "ho": true, "hr": true, "ht": true, "hu": true, "hy": true, "hz": true,
	"ia": true, "id": true, "ie": true, "ig": true, "ii": true, "ik": true, "io": true, "is": true,
	"it": true, "iu": true, "ja": true, "jv": true, "ka": true, "kg": true, "ki": true, "kj": true,
	"kk": true, "kl": true, "km": true, "kn": true, "ko": true, "kr": true, "ks": true, "ku": true,
	"kv": true, "kw": true, "ky": true, "la": true, "lb": true, "lg": true, "li": true, "ln": true,
	"lo": true, "lt": true, "lu": true, "lv": true, "mg": true, "mh": true, "mi": true, "mk": true,
	"ml": true, "mn": true, "mr": true, "ms": true, "mt": true, "my": true, "na": true, "nb": true,
	"nd": true, "ne": true, "ng": true, "nl": true, "nn": true, "no": true, "nr": true, "nv": true,
	"ny": true, "oc": true, "oj": true, "om": true, "or": true, "os": true, "pa": true, "pi": true,
	"pl": true, "ps": true, "pt": true, "qu": true, "rm": true, "rn": true, "ro": true, "ru": true,
	"rw": true, "sa": true, "sc": true, "sd": true, "se": true, "sg": true, "si": true, "sk": true,
	"sl": true, "sm": true, "sn": true, "so": true, "sq": true, "sr": true, "ss": true, "st": true,
	"su": true, "sv": true, "sw": true, "ta": true, "te": true, "tg": true, "th": true, "ti": true,
	"tk": true, "tl": true, "tn": true, "to": true, "tr": true, "ts": true, "tt": true, "tw": true,
	"ty": true, "ug": true, "uk": true, "ur": true, "uz": true, "ve": true, "vi": true, "vo": true,
	"wa": true, "wo": true, "xh": true, "yi": true, "yo": true, "za": true, "zh": true, "zu": true,
}

// Extended locale codes (language-region combinations like en-US, es-MX, en-us, es-mx)
var extendedLocaleRegex = regexp.MustCompile(`^[a-z]{2}-[a-zA-Z]{2}$`)

// Common query parameter names for locale
var localeQueryParams = []string{"lang", "locale", "language", "hl", "l"}

// Detector handles locale detection in URLs
type Detector struct {
	// Context-based detection to avoid false positives
	contextAware bool
}

// NewDetector creates a new locale detector
func NewDetector() *Detector {
	return &Detector{
		contextAware: true,
	}
}

// Detect analyzes a URL and extracts locale information
func (d *Detector) Detect(rawURL string) (*LocalizedURL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	result := &LocalizedURL{
		OriginalURL: rawURL,
		LocaleType:  LocaleTypeNone,
	}

	// Priority 1: Check subdomain
	if locale := d.detectSubdomain(u.Host); locale != "" {
		result.Locale = locale
		result.LocaleType = LocaleTypeSubdomain
		result.BaseURL = d.removeSubdomainLocale(rawURL, u, locale)
		return result, nil
	}

	// Priority 2: Check path prefix
	if locale, pos := d.detectPathPrefix(u.Path); locale != "" {
		result.Locale = locale
		result.LocaleType = LocaleTypePath
		result.Position = pos
		result.BaseURL = d.removePathLocale(rawURL, u, locale, pos)
		return result, nil
	}

	// Priority 3: Check query parameters
	if locale := d.detectQueryParam(u.Query()); locale != "" {
		result.Locale = locale
		result.LocaleType = LocaleTypeQuery
		result.BaseURL = d.removeQueryLocale(rawURL, u, locale)
		return result, nil
	}

	// No locale detected
	result.BaseURL = rawURL
	return result, nil
}

// detectSubdomain checks if the subdomain is a locale code
func (d *Detector) detectSubdomain(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return ""
	}

	firstPart := strings.ToLower(parts[0])

	// Check if it's a valid locale code
	if localeCodes[firstPart] {
		return firstPart
	}

	// Check extended format (en-us, es-mx)
	if extendedLocaleRegex.MatchString(firstPart) {
		return strings.ToLower(firstPart)
	}

	return ""
}

// detectPathPrefix checks if the path starts with a locale code
func (d *Detector) detectPathPrefix(path string) (string, int) {
	if path == "" || path == "/" {
		return "", -1
	}

	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) == 0 {
		return "", -1
	}

	// Check first segment
	if locale := d.validatePathSegmentAsLocale(segments[0], segments, 0); locale != "" {
		return locale, 0
	}

	// Check second segment (for patterns like /content/en/page)
	if len(segments) > 1 {
		if locale := d.validatePathSegmentAsLocale(segments[1], segments, 1); locale != "" {
			return locale, 1
		}
	}

	return "", -1
}

// validatePathSegmentAsLocale checks if a path segment is a locale with context awareness
func (d *Detector) validatePathSegmentAsLocale(segment string, allSegments []string, position int) string {
	segment = strings.ToLower(segment)

	// Basic check: is it a locale code?
	isLocale := localeCodes[segment] || extendedLocaleRegex.MatchString(segment)
	if !isLocale {
		return ""
	}

	// Context awareness to avoid false positives
	if d.contextAware {
		// Don't treat as locale if it's part of a word
		if strings.Contains(segment, "-") && !extendedLocaleRegex.MatchString(segment) {
			return ""
		}

		// Blacklist common false positives (very conservative)
		// Only reject if it's clearly NOT a locale code
		falsePositives := map[string]bool{
			"id": true, // Often used as identifier, not Indonesian
			"in": true, // Preposition, not Interlingua
			"is": true, // Verb, not Icelandic
			"or": true, // Conjunction, not Oriya
			"to": true, // Preposition, not Tonga
			"ad": true, // Advertisement, not Adyghe
			"as": true, // Conjunction, not Assamese
			"at": true, // Preposition, not ???
			"by": true, // Preposition, not Belarusian
			"go": true, // Verb/language, not ???
			"no": true, // Often "number", not Norwegian
		}

		// If segment is in false positives, reject it
		if falsePositives[segment] {
			return ""
		}

		// Special case: "it" - reject if in API or technical contexts
		if segment == "it" && position > 0 {
			if allSegments[0] == "api" || allSegments[0] == "tech" || allSegments[0] == "technology" {
				return ""
			}
		}

		// API endpoints usually don't have locale in path
		if position > 0 && allSegments[0] == "api" {
			// Unless there's clear evidence (more segments suggesting locale)
			if len(allSegments) < 3 {
				return ""
			}
		}
	}

	return segment
}

// detectQueryParam checks query parameters for locale
func (d *Detector) detectQueryParam(query url.Values) string {
	for _, param := range localeQueryParams {
		if val := query.Get(param); val != "" {
			val = strings.ToLower(val)
			if localeCodes[val] || extendedLocaleRegex.MatchString(val) {
				return val
			}
		}
	}
	return ""
}

// removeSubdomainLocale removes locale subdomain from URL
func (d *Detector) removeSubdomainLocale(rawURL string, u *url.URL, locale string) string {
	parts := strings.Split(u.Host, ".")
	if len(parts) < 2 {
		return rawURL
	}

	// Remove first part (locale)
	newHost := strings.Join(parts[1:], ".")

	newURL := *u
	newURL.Host = newHost
	return newURL.String()
}

// removePathLocale removes locale from path
func (d *Detector) removePathLocale(rawURL string, u *url.URL, locale string, position int) string {
	segments := strings.Split(strings.Trim(u.Path, "/"), "/")

	// Remove the locale segment
	newSegments := make([]string, 0, len(segments)-1)
	for i, seg := range segments {
		if i != position {
			newSegments = append(newSegments, seg)
		}
	}

	newPath := "/" + strings.Join(newSegments, "/")
	if len(newSegments) == 0 {
		newPath = "/"
	}

	newURL := *u
	newURL.Path = newPath
	return newURL.String()
}

// removeQueryLocale removes locale query parameter from URL
func (d *Detector) removeQueryLocale(rawURL string, u *url.URL, locale string) string {
	q := u.Query()

	// Remove all locale-related parameters
	for _, param := range localeQueryParams {
		if strings.ToLower(q.Get(param)) == locale {
			q.Del(param)
		}
	}

	newURL := *u
	newURL.RawQuery = q.Encode()
	return newURL.String()
}

// IsLocaleCode checks if a string is a valid locale code
func IsLocaleCode(code string) bool {
	code = strings.ToLower(code)
	return localeCodes[code] || extendedLocaleRegex.MatchString(code)
}
