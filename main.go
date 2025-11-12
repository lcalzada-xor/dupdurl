// dedup.go - Advanced URL deduplication tool for bug bounty pipelines
//
// Usage examples:
//   cat urls.txt | go run dedup.go                # print unique URLs (url mode)
//   cat urls.txt | go run dedup.go -mode=path     # dedupe by path (host+path)
//   cat urls.txt | go run dedup.go -mode=params   # dedupe by unique parameter names
//   cat urls.txt | go run dedup.go -counts        # print "count url" sorted by first appearance
//   cat urls.txt | go run dedup.go -ignore-params=utm_source,utm_medium -sort-params
//   cat urls.txt | go run dedup.go -fuzzy         # normalize numeric IDs in paths
//   cat urls.txt | go run dedup.go -ignore-extensions=jpg,png,css -stats
//   cat urls.txt | go run dedup.go -output=json   # JSON output format
//
// Build: go build -o dedup dedup.go

package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	// Core options
	mode          = flag.String("mode", "url", "normalization mode: url|path|host|raw|params")
	ignoreParams  = flag.String("ignore-params", "", "comma-separated query params to remove (e.g. utm_source,fbclid)")
	sortParams    = flag.Bool("sort-params", false, "sort query parameters alphabetically")
	ignoreFrag    = flag.Bool("ignore-fragment", true, "remove URL fragment (#...)")
	caseSensitive = flag.Bool("case-sensitive", false, "consider case when comparing paths/hosts")
	keepWWW       = flag.Bool("keep-www", false, "don't strip leading www. from host")
	keepScheme    = flag.Bool("keep-scheme", false, "distinguish between http:// and https://")
	trimSpaces    = flag.Bool("trim", true, "trim surrounding spaces")
	
	// Output options
	printCounts   = flag.Bool("counts", false, "print counts before each unique entry")
	outputFormat  = flag.String("output", "text", "output format: text|json|csv")
	showStats     = flag.Bool("stats", false, "print statistics at the end")
	verbose       = flag.Bool("verbose", false, "verbose mode: show warnings and parse errors")
	
	// Advanced normalization
	fuzzyMode         = flag.Bool("fuzzy", false, "replace numeric IDs in paths with {id} placeholder")
	pathIncludeQuery  = flag.Bool("path-include-query", false, "in path mode, include normalized query string")
	
	// Filtering
	allowDomains  = flag.String("allow-domains", "", "comma-separated list of allowed domains (whitelist)")
	blockDomains  = flag.String("block-domains", "", "comma-separated list of blocked domains (blacklist)")
	
	// Statistics
	stats Statistics
)

// Statistics tracks processing metrics
type Statistics struct {
	TotalProcessed int
	UniqueURLs     int
	Duplicates     int
	ParseErrors    int
	Filtered       int
}

// URLEntry represents a deduplicated URL with its count
type URLEntry struct {
	URL   string `json:"url"`
	Count int    `json:"count"`
}

func normalizeURL(raw string, ignoredSet, allowedDomains, blockedDomains map[string]struct{}) (string, error) {
	if *trimSpaces {
		raw = strings.TrimSpace(raw)
	}
	
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	// Check domain filtering
	if len(allowedDomains) > 0 {
		host := strings.ToLower(u.Host)
		if strings.HasPrefix(host, "www.") {
			host = strings.TrimPrefix(host, "www.")
		}
		if _, ok := allowedDomains[host]; !ok {
			return "", fmt.Errorf("domain not in whitelist: %s", u.Host)
		}
	}
	
	if len(blockedDomains) > 0 {
		host := strings.ToLower(u.Host)
		if strings.HasPrefix(host, "www.") {
			host = strings.TrimPrefix(host, "www.")
		}
		if _, ok := blockedDomains[host]; ok {
			return "", fmt.Errorf("domain in blacklist: %s", u.Host)
		}
	}

	// Lowercase scheme unless keepScheme or caseSensitive
	if !*caseSensitive && !*keepScheme {
		u.Scheme = strings.ToLower(u.Scheme)
	} else if !*keepScheme {
		u.Scheme = "https" // normalize to https if not keeping scheme
	}

	// Lowercase host unless caseSensitive
	if !*caseSensitive {
		u.Host = strings.ToLower(u.Host)
	}

	// Strip leading "www." by default
	if !*keepWWW {
		if strings.HasPrefix(u.Host, "www.") {
			u.Host = strings.TrimPrefix(u.Host, "www.")
		}
	}

	// Remove fragment
	if *ignoreFrag {
		u.Fragment = ""
	}

	// Normalize path
	u.Path = normalizePath(u.Path)

	// Apply fuzzy mode (replace numeric IDs)
	if *fuzzyMode {
		u.Path = fuzzyPath(u.Path)
	}

	// Query params handling - keep values by default
	q := u.Query()
	
	// Delete ignored params
	for p := range ignoredSet {
		q.Del(p)
	}
	
	if *sortParams {
		// Build sorted query with sorted values
		u.RawQuery = buildSortedQuery(q)
	} else {
		// Keep parameter values as-is
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

func normalizePath(p string) string {
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

var numericIDRegex = regexp.MustCompile(`/\d+(/|$)`)

func fuzzyPath(p string) string {
	// Replace numeric path segments with {id}
	// Example: /user/12345/profile -> /user/{id}/profile
	return numericIDRegex.ReplaceAllString(p, "/{id}$1")
}

func buildSortedQuery(q url.Values) string {
	if len(q) == 0 {
		return ""
	}
	
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	var sb strings.Builder
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

// createDedupKey creates a key for deduplication that includes parameter names but not values
// This is used when we want to deduplicate based on param structure but keep sample values
func createDedupKey(raw string, ignoredSet, allowedDomains, blockedDomains map[string]struct{}) (string, error) {
	if *trimSpaces {
		raw = strings.TrimSpace(raw)
	}
	
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	// Apply same normalization as normalizeURL
	if !*caseSensitive && !*keepScheme {
		u.Scheme = strings.ToLower(u.Scheme)
	} else if !*keepScheme {
		u.Scheme = "https"
	}

	if !*caseSensitive {
		u.Host = strings.ToLower(u.Host)
	}

	if !*keepWWW {
		if strings.HasPrefix(u.Host, "www.") {
			u.Host = strings.TrimPrefix(u.Host, "www.")
		}
	}

	if *ignoreFrag {
		u.Fragment = ""
	}

	u.Path = normalizePath(u.Path)

	if *fuzzyMode {
		u.Path = fuzzyPath(u.Path)
	}

	// For the dedup key, we only keep parameter NAMES, not values
	q := u.Query()
	
	// Delete ignored params
	for p := range ignoredSet {
		q.Del(p)
	}
	
	// Build query string with param names only (no values)
	if len(q) > 0 {
		keys := make([]string, 0, len(q))
		for k := range q {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		u.RawQuery = strings.Join(keys, "&") + "="
	} else {
		u.RawQuery = ""
	}

	return u.String(), nil
}

func extractParams(line string) (string, error) {
	u, err := url.Parse(line)
	if err != nil {
		return "", err
	}
	
	q := u.Query()
	if len(q) == 0 {
		return "", fmt.Errorf("no parameters")
	}
	
	// Get unique parameter names, sorted
	params := make([]string, 0, len(q))
	for k := range q {
		params = append(params, k)
	}
	sort.Strings(params)
	
	return strings.Join(params, ","), nil
}

func normalizeLine(line string, ignoredSet, allowedDomains, blockedDomains map[string]struct{}) (string, error) {
	if *trimSpaces {
		line = strings.TrimSpace(line)
	}
	if line == "" {
		return "", fmt.Errorf("empty line")
	}

	switch *mode {
	case "raw":
		if !*caseSensitive {
			return strings.ToLower(line), nil
		}
		return line, nil
		
	case "host":
		u, err := url.Parse(line)
		if err != nil {
			if !*caseSensitive {
				return strings.ToLower(line), nil
			}
			return line, nil
		}
		h := u.Host
		if !*keepWWW && strings.HasPrefix(h, "www.") {
			h = strings.TrimPrefix(h, "www.")
		}
		if !*caseSensitive {
			h = strings.ToLower(h)
		}
		return h, nil
		
	case "path":
		u, err := url.Parse(line)
		if err != nil {
			if !*caseSensitive {
				return strings.ToLower(line), nil
			}
			return line, nil
		}
		
		host := u.Host
		if !*keepWWW && strings.HasPrefix(host, "www.") {
			host = strings.TrimPrefix(host, "www.")
		}
		if !*caseSensitive {
			host = strings.ToLower(host)
		}
		
		path := normalizePath(u.Path)
		if *fuzzyMode {
			path = fuzzyPath(path)
		}
		
		result := host + path
		
		// Optionally include normalized query
		if *pathIncludeQuery && u.RawQuery != "" {
			q := u.Query()
			for p := range ignoredSet {
				q.Del(p)
			}
			if *sortParams {
				result += "?" + buildSortedQuery(q)
			} else {
				result += "?" + q.Encode()
			}
		}
		
		return result, nil
		
	case "params":
		return extractParams(line)
		
	case "url":
		return normalizeURL(line, ignoredSet, allowedDomains, blockedDomains)
		
	default:
		return "", fmt.Errorf("unknown mode: %s", *mode)
	}
}

func parseSet(s string) map[string]struct{} {
	m := map[string]struct{}{}
	if s == "" {
		return m
	}
	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		m[strings.ToLower(item)] = struct{}{}
	}
	return m
}

func outputText(entries []URLEntry) {
	for _, entry := range entries {
		if *printCounts {
			fmt.Printf("%d %s\n", entry.Count, entry.URL)
		} else {
			fmt.Println(entry.URL)
		}
	}
}

func outputJSON(entries []URLEntry) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(entries)
}

func outputCSV(entries []URLEntry) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	
	// Write header
	if err := w.Write([]string{"url", "count"}); err != nil {
		return err
	}
	
	// Write data
	for _, entry := range entries {
		if err := w.Write([]string{entry.URL, fmt.Sprintf("%d", entry.Count)}); err != nil {
			return err
		}
	}
	
	return nil
}

func printStats() {
	fmt.Fprintln(os.Stderr, "\n=== Statistics ===")
	fmt.Fprintf(os.Stderr, "Total URLs processed: %d\n", stats.TotalProcessed)
	fmt.Fprintf(os.Stderr, "Unique URLs:          %d\n", stats.UniqueURLs)
	fmt.Fprintf(os.Stderr, "Duplicates removed:   %d\n", stats.Duplicates)
	fmt.Fprintf(os.Stderr, "Parse errors:         %d\n", stats.ParseErrors)
	fmt.Fprintf(os.Stderr, "Filtered out:         %d\n", stats.Filtered)
	fmt.Fprintln(os.Stderr, "==================")
}

func main() {
	flag.Parse()
	
	// Parse sets
	ignoredSet := parseSet(*ignoreParams)
	allowedDomains := parseSet(*allowDomains)
	blockedDomains := parseSet(*blockDomains)

	// Maps to track first-seen URLs with parameter values
	seen := map[string]string{} // dedup key -> first full URL with values
	counts := map[string]int{}
	order := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		stats.TotalProcessed++
		
		if *trimSpaces && strings.TrimSpace(line) == "" {
			continue
		}
		
		// Create dedup key (without parameter values for comparison)
		key, err := createDedupKey(line, ignoredSet, allowedDomains, blockedDomains)
		if err != nil {
			if *verbose {
				fmt.Fprintf(os.Stderr, "Line %d: %v - %s\n", lineNum, err, line)
			}
			if strings.Contains(err.Error(), "parse error") {
				stats.ParseErrors++
			} else if strings.Contains(err.Error(), "ignored extension") || 
			          strings.Contains(err.Error(), "blacklist") ||
			          strings.Contains(err.Error(), "whitelist") {
				stats.Filtered++
			}
			continue
		}
		
		// Get normalized URL with values preserved
		normalizedURL, err := normalizeURL(line, ignoredSet, allowedDomains, blockedDomains)
		if err != nil {
			continue
		}
		
		// If this key hasn't been seen, store the first URL with its values
		if _, ok := seen[key]; !ok {
			seen[key] = normalizedURL
			order = append(order, key)
			stats.UniqueURLs++
		} else {
			stats.Duplicates++
		}
		counts[key]++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading stdin:", err)
		os.Exit(1)
	}

	// Prepare entries with first-seen URLs (which have parameter values)
	entries := make([]URLEntry, len(order))
	for i, k := range order {
		entries[i] = URLEntry{
			URL:   seen[k],
			Count: counts[k],
		}
	}

	// Output results
	switch *outputFormat {
	case "json":
		if err := outputJSON(entries); err != nil {
			fmt.Fprintln(os.Stderr, "error writing JSON:", err)
			os.Exit(1)
		}
	case "csv":
		if err := outputCSV(entries); err != nil {
			fmt.Fprintln(os.Stderr, "error writing CSV:", err)
			os.Exit(1)
		}
	case "text":
		outputText(entries)
	default:
		fmt.Fprintf(os.Stderr, "unknown output format: %s\n", *outputFormat)
		os.Exit(1)
	}

	// Print statistics if requested
	if *showStats {
		printStats()
	}
}

// dedupStandard implements the original deduplication logic
