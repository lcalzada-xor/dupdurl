package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lcalzada-xor/dupdurl/pkg/normalizer"
	"github.com/lcalzada-xor/dupdurl/pkg/processor"
)

// CLIConfig holds all command-line flags
type CLIConfig struct {
	// Core options
	Mode             string
	IgnoreParams     string
	SortParams       bool
	IgnoreFragment   bool
	CaseSensitive    bool
	KeepWWW          bool
	KeepScheme       bool
	TrimSpaces       bool

	// Output options
	PrintCounts      bool
	OutputFormat     string
	ShowStats        bool
	ShowStatsDetailed bool
	Verbose          bool

	// Advanced normalization
	FuzzyMode        bool
	FuzzyPatterns    string
	PathIncludeQuery bool
	IgnoreExtensions string

	// Filtering
	AllowDomains     string
	BlockDomains     string

	// Performance
	Workers          int
	BatchSize        int

	// Storage
	StorageBackend   string
	DBPath           string
}

// ParseFlags parses command-line flags and returns configuration
func ParseFlags() *CLIConfig {
	config := &CLIConfig{}

	// Core options
	flag.StringVar(&config.Mode, "mode", "url", "normalization mode: url|path|host|raw|params")
	flag.StringVar(&config.IgnoreParams, "ignore-params", "", "comma-separated query params to remove")
	flag.BoolVar(&config.SortParams, "sort-params", false, "sort query parameters alphabetically")
	flag.BoolVar(&config.IgnoreFragment, "ignore-fragment", true, "remove URL fragment (#...)")
	flag.BoolVar(&config.CaseSensitive, "case-sensitive", false, "consider case when comparing paths/hosts")
	flag.BoolVar(&config.KeepWWW, "keep-www", false, "don't strip leading www. from host")
	flag.BoolVar(&config.KeepScheme, "keep-scheme", false, "distinguish between http:// and https://")
	flag.BoolVar(&config.TrimSpaces, "trim", true, "trim surrounding spaces")

	// Output options
	flag.BoolVar(&config.PrintCounts, "counts", false, "print counts before each unique entry")
	flag.StringVar(&config.OutputFormat, "output", "text", "output format: text|json|csv")
	flag.BoolVar(&config.ShowStats, "stats", false, "print statistics at the end")
	flag.BoolVar(&config.ShowStatsDetailed, "stats-detailed", false, "print detailed statistics")
	flag.BoolVar(&config.Verbose, "verbose", false, "verbose mode: show warnings and parse errors")

	// Advanced normalization
	flag.BoolVar(&config.FuzzyMode, "fuzzy", false, "enable fuzzy matching for IDs")
	flag.StringVar(&config.FuzzyPatterns, "fuzzy-patterns", "numeric", "fuzzy patterns: numeric|uuid|hash|token (comma-separated)")
	flag.BoolVar(&config.PathIncludeQuery, "path-include-query", false, "in path mode, include normalized query string")
	flag.StringVar(&config.IgnoreExtensions, "ignore-extensions", "", "comma-separated extensions to skip (e.g., jpg,png,css)")

	// Filtering
	flag.StringVar(&config.AllowDomains, "allow-domains", "", "comma-separated list of allowed domains (whitelist)")
	flag.StringVar(&config.BlockDomains, "block-domains", "", "comma-separated list of blocked domains (blacklist)")

	// Performance
	flag.IntVar(&config.Workers, "workers", 1, "number of parallel workers (0 = NumCPU)")
	flag.IntVar(&config.BatchSize, "batch-size", 1000, "batch size for parallel processing")

	// Storage
	flag.StringVar(&config.StorageBackend, "storage", "memory", "storage backend: memory|sqlite")
	flag.StringVar(&config.DBPath, "db-path", ":memory:", "SQLite database path")

	flag.Parse()
	return config
}

// Validate checks if the configuration is valid
func (c *CLIConfig) Validate() error {
	// Validate mode
	validModes := []string{"url", "path", "host", "raw", "params"}
	if !contains(validModes, c.Mode) {
		return fmt.Errorf("invalid mode: %s (valid: %s)", c.Mode, strings.Join(validModes, ", "))
	}

	// Validate output format
	validFormats := []string{"text", "json", "csv"}
	if !contains(validFormats, c.OutputFormat) {
		return fmt.Errorf("invalid output format: %s (valid: %s)", c.OutputFormat, strings.Join(validFormats, ", "))
	}

	// Validate storage backend
	validBackends := []string{"memory", "sqlite"}
	if !contains(validBackends, c.StorageBackend) {
		return fmt.Errorf("invalid storage backend: %s (valid: %s)", c.StorageBackend, strings.Join(validBackends, ", "))
	}

	// Validate workers
	if c.Workers < 0 {
		return fmt.Errorf("workers must be >= 0")
	}

	// Validate batch size
	if c.BatchSize < 1 {
		return fmt.Errorf("batch-size must be >= 1")
	}

	return nil
}

// ToNormalizerConfig converts CLI config to normalizer config
func (c *CLIConfig) ToNormalizerConfig() *normalizer.Config {
	config := normalizer.NewConfig()

	config.Mode = c.Mode
	config.IgnoreParams = normalizer.ParseSet(c.IgnoreParams)
	config.SortParams = c.SortParams
	config.IgnoreFragment = c.IgnoreFragment
	config.CaseSensitive = c.CaseSensitive
	config.KeepWWW = c.KeepWWW
	config.KeepScheme = c.KeepScheme
	config.TrimSpaces = c.TrimSpaces
	config.FuzzyMode = c.FuzzyMode
	config.PathIncludeQuery = c.PathIncludeQuery
	config.AllowDomains = normalizer.ParseSet(c.AllowDomains)
	config.BlockDomains = normalizer.ParseSet(c.BlockDomains)
	config.IgnoreExtensions = normalizer.ParseSet(c.IgnoreExtensions)

	// Configure fuzzy patterns
	if c.FuzzyMode && c.FuzzyPatterns != "" {
		patterns := strings.Split(c.FuzzyPatterns, ",")
		normalizer.EnablePatterns(config.FuzzyPatterns, patterns)
	}

	return config
}

// ToProcessorConfig converts CLI config to processor config
func (c *CLIConfig) ToProcessorConfig() *processor.Config {
	config := processor.NewConfig()

	config.Normalizer = c.ToNormalizerConfig()
	config.Workers = c.Workers
	config.BatchSize = c.BatchSize
	config.Verbose = c.Verbose

	return config
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
