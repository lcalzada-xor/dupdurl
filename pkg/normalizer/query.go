package normalizer

import (
	"net/url"
	"sort"
	"strings"

	"github.com/lcalzada-xor/dupdurl/pkg/pool"
)

// BuildSortedQuery builds a query string with sorted keys and values
// Uses string builder pool for better performance
func BuildSortedQuery(q url.Values) string {
	if len(q) == 0 {
		return ""
	}

	// Pre-allocate with exact size
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Use pooled string builder
	sb := pool.GetBuilder()
	defer pool.PutBuilder(sb)

	first := true
	for _, k := range keys {
		vs := q[k]
		sort.Strings(vs) // Deterministic order for values too
		for _, v := range vs {
			if !first {
				sb.WriteByte('&')
			}
			sb.WriteString(url.QueryEscape(k))
			if v != "" {
				sb.WriteByte('=')
				sb.WriteString(url.QueryEscape(v))
			}
			first = false
		}
	}
	return sb.String()
}

// BuildKeyOnlyQuery builds a query string with parameter names only (no values)
// Used for deduplication keys
func BuildKeyOnlyQuery(q url.Values) string {
	if len(q) == 0 {
		return ""
	}

	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return strings.Join(keys, "&") + "="
}

// ParseSet parses a comma-separated string into a set
// Pre-allocates map with estimated size for better performance
func ParseSet(s string) map[string]struct{} {
	if s == "" {
		return make(map[string]struct{})
	}

	// Estimate size based on comma count
	estimatedSize := strings.Count(s, ",") + 1
	m := make(map[string]struct{}, estimatedSize)

	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		m[strings.ToLower(item)] = struct{}{}
	}
	return m
}

// ExtractParams extracts and sorts parameter names from a URL
func ExtractParams(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	if len(q) == 0 {
		return "", nil
	}

	// Get unique parameter names, sorted (normalized to lowercase)
	params := make([]string, 0, len(q))
	for k := range q {
		params = append(params, strings.ToLower(k))
	}
	sort.Strings(params)

	return strings.Join(params, ","), nil
}
