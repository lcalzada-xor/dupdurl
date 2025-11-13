package benchmark

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lcalzada-xor/dupdurl/pkg/normalizer"
	"github.com/lcalzada-xor/dupdurl/pkg/processor"
)

func BenchmarkNormalizePath(b *testing.B) {
	paths := []string{
		"/api/users/profile",
		"/api//users///profile",
		"/path/to/resource/",
		"api/users",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range paths {
			normalizer.NormalizePath(p)
		}
	}
}

func BenchmarkFuzzyPath(b *testing.B) {
	path := "/api/users/123/posts/456/comments/789"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		normalizer.FuzzyPath(path)
	}
}

func BenchmarkNormalizeURL(b *testing.B) {
	config := normalizer.NewConfig()
	url := "https://www.example.com/api/users?sort=name&filter=active#section"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.NormalizeURL(url)
	}
}

func BenchmarkProcessSequential(b *testing.B) {
	// Generate test data
	var input strings.Builder
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(&input, "https://example.com/page%d\n", i%100)
	}
	inputData := input.String()

	config := processor.NewConfig()
	config.Workers = 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proc := processor.New(config)
		proc.Process(strings.NewReader(inputData))
	}
}

func BenchmarkProcessParallel(b *testing.B) {
	// Generate test data
	var input strings.Builder
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(&input, "https://example.com/page%d\n", i%100)
	}
	inputData := input.String()

	config := processor.NewConfig()
	config.Workers = 4
	config.BatchSize = 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proc := processor.New(config)
		proc.Process(strings.NewReader(inputData))
	}
}

func BenchmarkLargeDataset(b *testing.B) {
	// Generate large dataset
	var input strings.Builder
	for i := 0; i < 100000; i++ {
		fmt.Fprintf(&input, "https://example.com/api/users/%d/profile?sort=date&filter=%d\n", i%1000, i%10)
	}
	inputData := input.String()

	config := processor.NewConfig()
	config.Workers = 4
	config.Normalizer.FuzzyMode = true

	b.ResetTimer()
	b.Run("100k URLs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			proc := processor.New(config)
			proc.Process(strings.NewReader(inputData))
		}
	})
}

func BenchmarkBuildSortedQuery(b *testing.B) {
	query := map[string][]string{
		"sort":   {"name", "date", "id"},
		"filter": {"active", "pending"},
		"page":   {"1"},
		"limit":  {"100"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		normalizer.BuildSortedQuery(query)
	}
}
