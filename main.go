// dupdurl - Advanced URL deduplication tool for bug bounty pipelines
//
// A fast, powerful URL deduplication tool designed for security researchers
// and bug bounty hunters. Features fuzzy matching, parameter filtering,
// parallel processing, and multi-format output.
//
// Usage examples:
//   cat urls.txt | dupdurl                    # print unique URLs
//   cat urls.txt | dupdurl -fuzzy             # with fuzzy ID matching
//   cat urls.txt | dupdurl -mode=path         # dedupe by path only
//   cat urls.txt | dupdurl -workers=4         # parallel processing
//   cat urls.txt | dupdurl -output=json       # JSON output
//   cat urls.txt | dupdurl -storage=sqlite    # use SQLite for massive datasets
//
// Build: go build -o dupdurl

package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/lcalzada-xor/dupdurl/pkg/config"
	"github.com/lcalzada-xor/dupdurl/pkg/deduplicator"
	"github.com/lcalzada-xor/dupdurl/pkg/diff"
	"github.com/lcalzada-xor/dupdurl/pkg/normalizer"
	"github.com/lcalzada-xor/dupdurl/pkg/output"
	"github.com/lcalzada-xor/dupdurl/pkg/processor"
	"github.com/lcalzada-xor/dupdurl/pkg/scope"
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

	// Config file
	ConfigFile string
	Profile    string
	SaveConfig string

	// Diff mode
	DiffBaseline string
	SaveBaseline string

	// Streaming mode
	Streaming              bool
	StreamingFlushInterval string
	StreamingMaxBuffer     int

	// Scope checking
	ScopeFile      string
	OutOfScope     bool
	ScopeStats     bool
}

// ParseFlags parses command-line flags and returns configuration
func ParseFlags() *CLIConfig {
	config := &CLIConfig{}

	// Override default Usage to show custom help
	flag.Usage = printUsage

	// === CORE NORMALIZATION OPTIONS ===
	flag.StringVar(&config.Mode, "mode", "url", "")
	flag.StringVar(&config.Mode, "m", "url", "")

	flag.BoolVar(&config.FuzzyMode, "fuzzy", false, "")
	flag.BoolVar(&config.FuzzyMode, "f", false, "")

	flag.StringVar(&config.FuzzyPatterns, "fuzzy-patterns", "numeric", "")
	flag.StringVar(&config.FuzzyPatterns, "fp", "numeric", "")

	flag.BoolVar(&config.IgnoreFragment, "ignore-fragment", true, "")
	flag.BoolVar(&config.CaseSensitive, "case-sensitive", false, "")
	flag.BoolVar(&config.KeepWWW, "keep-www", false, "")
	flag.BoolVar(&config.KeepScheme, "keep-scheme", false, "")
	flag.BoolVar(&config.TrimSpaces, "trim", true, "")
	flag.BoolVar(&config.TrimSpaces, "t", true, "")

	// === PARAMETER & QUERY HANDLING ===
	flag.StringVar(&config.IgnoreParams, "ignore-params", "", "")
	flag.StringVar(&config.IgnoreParams, "ip", "", "")

	flag.BoolVar(&config.SortParams, "sort-params", false, "")
	flag.BoolVar(&config.SortParams, "sp", false, "")

	flag.BoolVar(&config.PathIncludeQuery, "path-include-query", false, "")

	// === FILTERING OPTIONS ===
	flag.StringVar(&config.IgnoreExtensions, "ignore-extensions", "", "")
	flag.StringVar(&config.IgnoreExtensions, "ie", "", "")

	flag.StringVar(&config.AllowDomains, "allow-domains", "", "")
	flag.StringVar(&config.AllowDomains, "ad", "", "")

	flag.StringVar(&config.BlockDomains, "block-domains", "", "")
	flag.StringVar(&config.BlockDomains, "bd", "", "")

	// === OUTPUT OPTIONS ===
	flag.StringVar(&config.OutputFormat, "output", "text", "")
	flag.StringVar(&config.OutputFormat, "o", "text", "")

	flag.BoolVar(&config.PrintCounts, "counts", false, "")
	flag.BoolVar(&config.PrintCounts, "c", false, "")

	flag.BoolVar(&config.ShowStats, "stats", false, "")
	flag.BoolVar(&config.ShowStats, "s", false, "")

	flag.BoolVar(&config.ShowStatsDetailed, "stats-detailed", false, "")
	flag.BoolVar(&config.ShowStatsDetailed, "sd", false, "")

	flag.BoolVar(&config.Verbose, "verbose", false, "")
	flag.BoolVar(&config.Verbose, "v", false, "")

	// === PERFORMANCE OPTIONS ===
	flag.IntVar(&config.Workers, "workers", 1, "")
	flag.IntVar(&config.Workers, "w", 1, "")

	flag.IntVar(&config.BatchSize, "batch-size", 1000, "")

	// === STREAMING MODE ===
	flag.BoolVar(&config.Streaming, "stream", false, "")
	flag.StringVar(&config.StreamingFlushInterval, "stream-interval", "5s", "")
	flag.IntVar(&config.StreamingMaxBuffer, "stream-buffer", 10000, "")

	// === DIFF MODE ===
	flag.StringVar(&config.DiffBaseline, "diff", "", "")
	flag.StringVar(&config.DiffBaseline, "d", "", "")

	flag.StringVar(&config.SaveBaseline, "save-baseline", "", "")
	flag.StringVar(&config.SaveBaseline, "sb", "", "")

	// === CONFIG FILE ===
	flag.StringVar(&config.ConfigFile, "config", "", "")
	flag.StringVar(&config.Profile, "profile", "", "")
	flag.StringVar(&config.Profile, "p", "", "")
	flag.StringVar(&config.SaveConfig, "save-config", "", "")

	// === STORAGE OPTIONS ===
	flag.StringVar(&config.StorageBackend, "storage", "memory", "")
	flag.StringVar(&config.DBPath, "db-path", ":memory:", "")

	// === SCOPE CHECKING ===
	flag.StringVar(&config.ScopeFile, "scope", "", "")
	flag.StringVar(&config.ScopeFile, "S", "", "")
	flag.BoolVar(&config.OutOfScope, "out-of-scope", false, "")
	flag.BoolVar(&config.ScopeStats, "scope-stats", false, "")

	flag.Parse()
	return config
}

// printUsage prints a professional, categorized help message
func printUsage() {
	fmt.Fprintf(os.Stderr, `dupdurl v2.1 - URL Deduplication Tool for Bug Bounty & Pentesting

USAGE:
  dupdurl [OPTIONS] < urls.txt
  cat urls.txt | dupdurl [OPTIONS]
  waybackurls target.com | dupdurl --fuzzy --stats

CORE NORMALIZATION OPTIONS:
  -m, --mode <mode>              Normalization mode (default: url)
                                 Modes: url, path, host, params, raw
  -f, --fuzzy                    Enable fuzzy matching for IDs
  -fp, --fuzzy-patterns <list>   Fuzzy patterns (default: numeric)
                                 Patterns: numeric, uuid, hash, token (comma-separated)
  --ignore-fragment              Remove URL fragments (#...) (default: true)
  --case-sensitive               Consider case when comparing paths/hosts
  --keep-www                     Don't strip leading www. from host
  --keep-scheme                  Distinguish between http:// and https://
  -t, --trim                     Trim surrounding spaces (default: true)

PARAMETER & QUERY HANDLING:
  -ip, --ignore-params <list>    Comma-separated query params to remove
                                 Example: utm_source,utm_medium,fbclid
  -sp, --sort-params             Sort query parameters alphabetically
  --path-include-query           In path mode, include normalized query string

FILTERING OPTIONS:
  -ie, --ignore-extensions <ext> Comma-separated extensions to skip
                                 Example: jpg,png,css,js,woff
  -ad, --allow-domains <list>    Only process these domains (whitelist)
  -bd, --block-domains <list>    Skip these domains (blacklist)

OUTPUT OPTIONS:
  -o, --output <format>          Output format (default: text)
                                 Formats: text, json, csv
  -c, --counts                   Print occurrence counts for each URL
  -s, --stats                    Print basic statistics at the end
  -sd, --stats-detailed          Print detailed statistics with analysis
  -v, --verbose                  Show warnings and parse errors

PERFORMANCE OPTIONS:
  -w, --workers <n>              Number of parallel workers (default: 1)
                                 Set to 0 for NumCPU
  --batch-size <n>               Batch size for parallel processing (default: 1000)

STREAMING MODE (NEW in v2.1):
  --stream                       Process infinite datasets with periodic flush
  --stream-interval <duration>   Flush interval (default: 5s)
                                 Examples: 1s, 30s, 1m
  --stream-buffer <n>            Max buffer size before forced flush (default: 10000)

DIFF MODE (NEW in v2.1):
  -d, --diff <baseline.json>     Compare against baseline JSON file
  -sb, --save-baseline <file>    Save results as baseline JSON file

CONFIG FILE (NEW in v2.1):
  --config <path>                Path to config file
                                 Default: ~/.config/dupdurl/config.yml
  -p, --profile <name>           Use predefined profile
                                 Profiles: bugbounty, aggressive, conservative
  --save-config <path>           Save current settings to config file

SCOPE CHECKING (NEW in v2.1):
  -S, --scope <file>             Scope file with domain patterns
                                 Supports wildcards: *.example.com
                                 Exclude with !: !staging.example.com
  --out-of-scope                 Show only out-of-scope URLs
  --scope-stats                  Show in/out scope statistics

STORAGE OPTIONS:
  --storage <backend>            Storage backend (default: memory)
                                 Backends: memory, sqlite
  --db-path <path>               SQLite database path (default: :memory:)

EXAMPLES:
  Basic deduplication:
    cat urls.txt | dupdurl

  Bug bounty workflow:
    waybackurls target.com | dupdurl --fuzzy --ignore-params=utm_source,fbclid --stats

  With parallel processing:
    cat urls.txt | dupdurl --workers=4 --output=json

  Streaming mode for live logs:
    tail -f access.log | dupdurl --stream --stream-interval=10s --fuzzy

  Diff mode for change tracking:
    waybackurls target.com | dupdurl --save-baseline=day1.json
    waybackurls target.com | dupdurl --diff=day1.json

  Using config profiles:
    cat urls.txt | dupdurl --profile=bugbounty

  Scope checking:
    echo "*.example.com" > scope.txt
    echo "!dev.example.com" >> scope.txt
    cat urls.txt | dupdurl --scope=scope.txt --scope-stats

MORE INFO:
  Documentation: https://github.com/lcalzada-xor/dupdurl
  Report bugs:   https://github.com/lcalzada-xor/dupdurl/issues

`)
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

func main() {
	// Parse command-line flags
	cliConfig := ParseFlags()

	// Load config file if specified (or use default location)
	var fileConfig *config.File
	if cliConfig.ConfigFile != "" {
		var err error
		fileConfig, err = config.Load(cliConfig.ConfigFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Try to load from default location
		fileConfig = config.LoadOrDefault()
	}

	// Apply profile if specified
	if cliConfig.Profile != "" {
		if err := fileConfig.ApplyProfile(cliConfig.Profile); err != nil {
			fmt.Fprintf(os.Stderr, "Error applying profile: %v\n", err)
			os.Exit(1)
		}
	}

	// Merge file config with CLI flags (CLI flags take precedence)
	mergeConfigs(cliConfig, fileConfig)

	// Save config if requested
	if cliConfig.SaveConfig != "" {
		if err := fileConfig.Save(cliConfig.SaveConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Config saved to %s\n", cliConfig.SaveConfig)
		return
	}

	// Validate configuration
	if err := cliConfig.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Use -h or --help for usage information\n")
		os.Exit(1)
	}

	// Auto-detect number of workers if set to 0
	if cliConfig.Workers == 0 {
		cliConfig.Workers = runtime.NumCPU()
	}

	// Load scope checker if specified
	var scopeChecker *scope.Checker
	if cliConfig.ScopeFile != "" {
		scopeChecker = scope.NewChecker()
		if err := scopeChecker.LoadFromFile(cliConfig.ScopeFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading scope file: %v\n", err)
			os.Exit(1)
		}
		if cliConfig.Verbose {
			stats := scopeChecker.GetStats()
			fmt.Fprintf(os.Stderr, "Scope loaded: %d includes, %d excludes\n",
				stats.IncludePatterns, stats.ExcludePatterns)
		}
	}

	// Check if we're in diff mode
	var differ *diff.Differ
	if cliConfig.DiffBaseline != "" {
		differ = diff.NewDiffer()
		if err := differ.LoadBaseline(cliConfig.DiffBaseline); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading baseline: %v\n", err)
			os.Exit(1)
		}
	}

	// Get output formatter
	formatter, err := output.GetFormatter(cliConfig.OutputFormat, cliConfig.PrintCounts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating formatter: %v\n", err)
		os.Exit(1)
	}

	var entries []deduplicator.Entry

	// Choose processing mode: streaming or batch
	if cliConfig.Streaming {
		// Streaming mode
		streamConfig := processor.NewStreamingConfig()
		streamConfig.Normalizer = cliConfig.ToNormalizerConfig()
		streamConfig.Workers = cliConfig.Workers
		streamConfig.Verbose = cliConfig.Verbose
		streamConfig.Output = formatter
		streamConfig.OutputWriter = os.Stdout

		// Parse flush interval
		if cliConfig.StreamingFlushInterval != "" {
			interval, err := time.ParseDuration(cliConfig.StreamingFlushInterval)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid flush interval: %v\n", err)
				os.Exit(1)
			}
			streamConfig.FlushInterval = interval
		}

		if cliConfig.StreamingMaxBuffer > 0 {
			streamConfig.MaxBuffer = cliConfig.StreamingMaxBuffer
		}

		streamProc := processor.NewStreaming(streamConfig)
		if err := streamProc.ProcessStreaming(os.Stdin); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing URLs: %v\n", err)
			os.Exit(1)
		}

		// Print statistics if requested
		stats := streamProc.GetStatistics()
		if cliConfig.ShowStatsDetailed {
			stats.PrintDetailed(os.Stderr)
		} else if cliConfig.ShowStats {
			stats.Print(os.Stderr)
		}

		return
	}

	// Batch mode (original behavior)
	procConfig := cliConfig.ToProcessorConfig()
	proc := processor.New(procConfig)

	entries, err = proc.Process(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing URLs: %v\n", err)
		os.Exit(1)
	}

	// Apply scope filtering if specified
	if scopeChecker != nil {
		// Count stats BEFORE filtering
		if cliConfig.ScopeStats {
			inScope, outScope := countScopeStats(entries, scopeChecker)
			fmt.Fprintf(os.Stderr, "\n=== Scope Statistics ===\n")
			fmt.Fprintf(os.Stderr, "In scope:     %d URLs\n", inScope)
			fmt.Fprintf(os.Stderr, "Out of scope: %d URLs\n", outScope)
			fmt.Fprintf(os.Stderr, "========================\n\n")
		}

		// Then filter
		entries = filterByScope(entries, scopeChecker, cliConfig.OutOfScope)
	}

	// Save baseline if requested
	if cliConfig.SaveBaseline != "" {
		if err := diff.SaveBaseline(entries, cliConfig.SaveBaseline); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving baseline: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Baseline saved to %s\n", cliConfig.SaveBaseline)
	}

	// Diff mode
	if differ != nil {
		report := differ.Compare(entries)
		report.PrintReport(os.Stderr)
		fmt.Fprintf(os.Stderr, "\nSummary: %s\n", report.Summary())
		return
	}

	// Output results
	if err := formatter.Format(entries, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	// Print statistics if requested
	stats := proc.GetStatistics()
	if cliConfig.ShowStatsDetailed {
		stats.PrintDetailed(os.Stderr)
	} else if cliConfig.ShowStats {
		stats.Print(os.Stderr)
	}
}

// mergeConfigs merges file config with CLI config (CLI takes precedence)
func mergeConfigs(cli *CLIConfig, file *config.File) {
	// Only apply file config if CLI flag wasn't explicitly set
	// This is simplified - in production you'd track which flags were actually set
	if cli.Mode == "url" && file.Mode != "" {
		cli.Mode = file.Mode
	}
	if !cli.FuzzyMode && file.FuzzyMode {
		cli.FuzzyMode = file.FuzzyMode
	}
	if cli.Workers == 1 && file.Workers > 0 {
		cli.Workers = file.Workers
	}
	// Add more field merging as needed
}

// filterByScope filters entries based on scope checker
func filterByScope(entries []deduplicator.Entry, checker *scope.Checker, showOutOfScope bool) []deduplicator.Entry {
	if checker == nil {
		return entries
	}

	filtered := make([]deduplicator.Entry, 0, len(entries))
	for _, entry := range entries {
		// Parse URL to extract host
		u, err := url.Parse(entry.URL)
		if err != nil {
			// If can't parse, skip it
			continue
		}

		inScope := checker.IsInScope(u.Host)

		// Include based on mode
		if showOutOfScope {
			// Show only out-of-scope URLs
			if !inScope {
				filtered = append(filtered, entry)
			}
		} else {
			// Show only in-scope URLs (default)
			if inScope {
				filtered = append(filtered, entry)
			}
		}
	}

	return filtered
}

// countScopeStats counts in-scope and out-of-scope URLs
func countScopeStats(entries []deduplicator.Entry, checker *scope.Checker) (inScope, outScope int) {
	for _, entry := range entries {
		u, err := url.Parse(entry.URL)
		if err != nil {
			continue
		}

		if checker.IsInScope(u.Host) {
			inScope++
		} else {
			outScope++
		}
	}
	return
}
