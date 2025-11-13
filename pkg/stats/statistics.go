package stats

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Statistics tracks processing metrics
type Statistics struct {
	TotalProcessed int
	UniqueURLs     int
	Duplicates     int
	ParseErrors    int
	Filtered       int
	StartTime      time.Time
	EndTime        time.Time

	// Enhanced statistics
	TopDomains     map[string]int
	ParamFrequency map[string]int
	ExtensionCount map[string]int
	totalParams    int
}

// NewStatistics creates a new Statistics instance
func NewStatistics() *Statistics {
	return &Statistics{
		StartTime:      time.Now(),
		TopDomains:     make(map[string]int),
		ParamFrequency: make(map[string]int),
		ExtensionCount: make(map[string]int),
	}
}

// Finish marks the end of processing
func (s *Statistics) Finish() {
	s.EndTime = time.Now()
}

// ProcessingTime returns the total processing duration
func (s *Statistics) ProcessingTime() time.Duration {
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// AvgQueryParams returns the average number of query parameters per URL
func (s *Statistics) AvgQueryParams() float64 {
	if s.UniqueURLs == 0 {
		return 0
	}
	return float64(s.totalParams) / float64(s.UniqueURLs)
}

// RecordDomain records a domain occurrence
func (s *Statistics) RecordDomain(domain string) {
	s.TopDomains[domain]++
}

// RecordParam records a parameter occurrence
func (s *Statistics) RecordParam(param string) {
	s.ParamFrequency[param]++
	s.totalParams++
}

// RecordExtension records an extension occurrence
func (s *Statistics) RecordExtension(ext string) {
	s.ExtensionCount[ext]++
}

// Print outputs basic statistics to the given writer
func (s *Statistics) Print(w io.Writer) {
	fmt.Fprintln(w, "\n=== Statistics ===")
	fmt.Fprintf(w, "Total URLs processed: %d\n", s.TotalProcessed)
	fmt.Fprintf(w, "Unique URLs:          %d\n", s.UniqueURLs)
	fmt.Fprintf(w, "Duplicates removed:   %d\n", s.Duplicates)
	fmt.Fprintf(w, "Parse errors:         %d\n", s.ParseErrors)
	fmt.Fprintf(w, "Filtered out:         %d\n", s.Filtered)
	fmt.Fprintf(w, "Processing time:      %v\n", s.ProcessingTime())
	fmt.Fprintln(w, "==================")
}

// PrintDetailed outputs detailed statistics to the given writer
func (s *Statistics) PrintDetailed(w io.Writer) {
	s.Print(w)

	// Top domains
	if len(s.TopDomains) > 0 {
		fmt.Fprintln(w, "\n=== Top Domains ===")
		topDomains := s.getTopN(s.TopDomains, 10)
		for i, kv := range topDomains {
			fmt.Fprintf(w, "%d. %s: %d\n", i+1, kv.Key, kv.Value)
		}
	}

	// Top parameters
	if len(s.ParamFrequency) > 0 {
		fmt.Fprintln(w, "\n=== Top Parameters ===")
		topParams := s.getTopN(s.ParamFrequency, 10)
		for i, kv := range topParams {
			fmt.Fprintf(w, "%d. %s: %d\n", i+1, kv.Key, kv.Value)
		}
		fmt.Fprintf(w, "Average params per URL: %.2f\n", s.AvgQueryParams())
	}

	// Extensions
	if len(s.ExtensionCount) > 0 {
		fmt.Fprintln(w, "\n=== File Extensions ===")
		topExts := s.getTopN(s.ExtensionCount, 10)
		for i, kv := range topExts {
			fmt.Fprintf(w, "%d. .%s: %d\n", i+1, kv.Key, kv.Value)
		}
	}
}

// KeyValue represents a key-value pair for sorting
type KeyValue struct {
	Key   string
	Value int
}

// getTopN returns the top N items from a map by value
func (s *Statistics) getTopN(m map[string]int, n int) []KeyValue {
	pairs := make([]KeyValue, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, KeyValue{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	if len(pairs) > n {
		pairs = pairs[:n]
	}

	return pairs
}

// ToJSON returns statistics as a JSON-compatible map
func (s *Statistics) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"total_processed":    s.TotalProcessed,
		"unique_urls":        s.UniqueURLs,
		"duplicates":         s.Duplicates,
		"parse_errors":       s.ParseErrors,
		"filtered":           s.Filtered,
		"processing_time_ms": s.ProcessingTime().Milliseconds(),
		"avg_query_params":   s.AvgQueryParams(),
		"top_domains":        s.getTopN(s.TopDomains, 10),
		"top_parameters":     s.getTopN(s.ParamFrequency, 10),
		"extensions":         s.getTopN(s.ExtensionCount, 10),
	}
}
