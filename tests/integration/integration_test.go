package integration

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lcalzada-xor/dupdurl/pkg/deduplicator"
	"github.com/lcalzada-xor/dupdurl/pkg/normalizer"
	"github.com/lcalzada-xor/dupdurl/pkg/output"
	"github.com/lcalzada-xor/dupdurl/pkg/processor"
)

func TestEndToEndBasic(t *testing.T) {
	input := `https://www.example.com/page1
https://example.com/page1
https://example.com/page2
`

	config := processor.NewConfig()
	config.Normalizer = normalizer.NewConfig()
	config.Workers = 1

	proc := processor.New(config)
	entries, err := proc.Process(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("Expected 2 unique URLs, got %d", len(entries))
	}

	stats := proc.GetStatistics()
	if stats.TotalProcessed != 3 {
		t.Errorf("TotalProcessed = %d; want 3", stats.TotalProcessed)
	}
	if stats.UniqueURLs != 2 {
		t.Errorf("UniqueURLs = %d; want 2", stats.UniqueURLs)
	}
	if stats.Duplicates != 1 {
		t.Errorf("Duplicates = %d; want 1", stats.Duplicates)
	}
}

func TestEndToEndFuzzyMode(t *testing.T) {
	input := `https://example.com/api/users/123/profile
https://example.com/api/users/456/profile
https://example.com/api/users/789/profile
https://example.com/api/posts/111/comments
`

	config := processor.NewConfig()
	config.Normalizer = normalizer.NewConfig()
	config.Normalizer.FuzzyMode = true
	config.Workers = 1

	proc := processor.New(config)
	entries, err := proc.Process(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	// With fuzzy mode, all user profile URLs should deduplicate to one
	if len(entries) != 2 {
		t.Errorf("Expected 2 unique URLs with fuzzy mode, got %d", len(entries))
	}

	// First entry should have count of 3
	if entries[0].Count != 3 {
		t.Errorf("First entry count = %d; want 3", entries[0].Count)
	}
}

func TestEndToEndParallelProcessing(t *testing.T) {
	// Create a large input to test parallel processing
	var input strings.Builder
	for i := 0; i < 1000; i++ {
		input.WriteString("https://example.com/page")
		input.WriteString(strings.Repeat("1", i%10))
		input.WriteString("\n")
	}

	config := processor.NewConfig()
	config.Normalizer = normalizer.NewConfig()
	config.Workers = 4
	config.BatchSize = 100

	proc := processor.New(config)
	entries, err := proc.Process(strings.NewReader(input.String()))

	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	if len(entries) == 0 {
		t.Error("Expected non-zero unique URLs")
	}

	stats := proc.GetStatistics()
	if stats.TotalProcessed != 1000 {
		t.Errorf("TotalProcessed = %d; want 1000", stats.TotalProcessed)
	}
}

func TestEndToEndIgnoreParams(t *testing.T) {
	input := `https://example.com/page?utm_source=google&id=123
https://example.com/page?utm_source=facebook&id=123
https://example.com/page?id=123
`

	config := processor.NewConfig()
	config.Normalizer = normalizer.NewConfig()
	config.Normalizer.IgnoreParams = normalizer.ParseSet("utm_source")
	config.Workers = 1

	proc := processor.New(config)
	entries, err := proc.Process(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	// All should deduplicate to one URL (ignoring utm_source)
	if len(entries) != 1 {
		t.Errorf("Expected 1 unique URL, got %d", len(entries))
	}
}

func TestOutputFormatters(t *testing.T) {
	entries := []deduplicator.Entry{
		{URL: "https://example.com/page1", Count: 2},
		{URL: "https://example.com/page2", Count: 1},
	}

	tests := []struct {
		name   string
		format string
		want   []string
	}{
		{
			name:   "text format",
			format: "text",
			want:   []string{"https://example.com/page1", "https://example.com/page2"},
		},
		{
			name:   "json format",
			format: "json",
			want:   []string{`"url"`, `"count"`, "example.com/page1"},
		},
		{
			name:   "csv format",
			format: "csv",
			want:   []string{"url,count", "example.com/page1,2", "example.com/page2,1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := output.GetFormatter(tt.format, false)
			if err != nil {
				t.Fatalf("GetFormatter() error = %v", err)
			}

			var buf bytes.Buffer
			if err := formatter.Format(entries, &buf); err != nil {
				t.Fatalf("Format() error = %v", err)
			}

			output := buf.String()
			for _, want := range tt.want {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing expected content %q\nGot: %s", want, output)
				}
			}
		})
	}
}

func TestEndToEndExtensionFilter(t *testing.T) {
	input := `https://example.com/image.jpg
https://example.com/style.css
https://example.com/page.html
https://example.com/api/data
`

	config := processor.NewConfig()
	config.Normalizer = normalizer.NewConfig()
	config.Normalizer.IgnoreExtensions = normalizer.ParseSet("jpg,css")
	config.Workers = 1

	proc := processor.New(config)
	entries, err := proc.Process(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	// Only .html and api/data should remain
	if len(entries) != 2 {
		t.Errorf("Expected 2 unique URLs after filtering, got %d", len(entries))
	}

	// Verify the correct URLs are present
	hasHTML := false
	hasData := false
	for _, entry := range entries {
		if strings.Contains(entry.URL, "page.html") {
			hasHTML = true
		}
		if strings.Contains(entry.URL, "api/data") {
			hasData = true
		}
		if strings.Contains(entry.URL, ".jpg") || strings.Contains(entry.URL, ".css") {
			t.Errorf("Filtered extension found in results: %s", entry.URL)
		}
	}

	if !hasHTML {
		t.Error("Expected page.html in results")
	}
	if !hasData {
		t.Error("Expected api/data in results")
	}
}

func TestEndToEndDomainFilter(t *testing.T) {
	input := `https://example.com/page1
https://test.com/page1
https://allowed.com/page1
`

	config := processor.NewConfig()
	config.Normalizer = normalizer.NewConfig()
	config.Normalizer.AllowDomains = normalizer.ParseSet("example.com,allowed.com")
	config.Workers = 1

	proc := processor.New(config)
	entries, err := proc.Process(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	// Only example.com and allowed.com should pass
	if len(entries) != 2 {
		t.Errorf("Expected 2 unique URLs after domain filtering, got %d", len(entries))
	}

	// Verify correct domains are present
	hasExample := false
	hasAllowed := false
	for _, entry := range entries {
		if strings.Contains(entry.URL, "example.com") {
			hasExample = true
		}
		if strings.Contains(entry.URL, "allowed.com") {
			hasAllowed = true
		}
		if strings.Contains(entry.URL, "test.com") {
			t.Errorf("Filtered domain found in results: %s", entry.URL)
		}
	}

	if !hasExample {
		t.Error("Expected example.com in results")
	}
	if !hasAllowed {
		t.Error("Expected allowed.com in results")
	}
}
