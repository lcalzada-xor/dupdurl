package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/lcalzada-xor/dupdurl/pkg/deduplicator"
)

// Formatter defines the interface for output formatters
type Formatter interface {
	Format(entries []deduplicator.Entry, w io.Writer) error
}

// TextFormatter outputs URLs as plain text
type TextFormatter struct {
	PrintCounts bool
}

// Format writes entries as plain text
func (f *TextFormatter) Format(entries []deduplicator.Entry, w io.Writer) error {
	for _, entry := range entries {
		if f.PrintCounts {
			fmt.Fprintf(w, "%d %s\n", entry.Count, entry.URL)
		} else {
			fmt.Fprintln(w, entry.URL)
		}
	}
	return nil
}

// JSONFormatter outputs URLs as JSON
type JSONFormatter struct{}

// Format writes entries as JSON
func (f *JSONFormatter) Format(entries []deduplicator.Entry, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(entries)
}

// CSVFormatter outputs URLs as CSV
type CSVFormatter struct{}

// Format writes entries as CSV
func (f *CSVFormatter) Format(entries []deduplicator.Entry, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"url", "count"}); err != nil {
		return err
	}

	// Write data
	for _, entry := range entries {
		if err := writer.Write([]string{entry.URL, fmt.Sprintf("%d", entry.Count)}); err != nil {
			return err
		}
	}

	return nil
}

// GetFormatter returns the appropriate formatter based on format string
func GetFormatter(format string, printCounts bool) (Formatter, error) {
	switch format {
	case "text":
		return &TextFormatter{PrintCounts: printCounts}, nil
	case "json":
		return &JSONFormatter{}, nil
	case "csv":
		return &CSVFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown output format: %s", format)
	}
}
