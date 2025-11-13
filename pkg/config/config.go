package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// File represents the complete config file structure
type File struct {
	// Core options
	Mode             string   `yaml:"mode"`
	IgnoreParams     []string `yaml:"ignore-params"`
	SortParams       bool     `yaml:"sort-params"`
	IgnoreFragment   bool     `yaml:"ignore-fragment"`
	CaseSensitive    bool     `yaml:"case-sensitive"`
	KeepWWW          bool     `yaml:"keep-www"`
	KeepScheme       bool     `yaml:"keep-scheme"`
	TrimSpaces       bool     `yaml:"trim-spaces"`

	// Output options
	PrintCounts      bool   `yaml:"print-counts"`
	OutputFormat     string `yaml:"output-format"`
	ShowStats        bool   `yaml:"show-stats"`
	ShowStatsDetailed bool  `yaml:"show-stats-detailed"`
	Verbose          bool   `yaml:"verbose"`

	// Advanced normalization
	FuzzyMode        bool     `yaml:"fuzzy"`
	FuzzyPatterns    []string `yaml:"fuzzy-patterns"`
	PathIncludeQuery bool     `yaml:"path-include-query"`
	IgnoreExtensions []string `yaml:"ignore-extensions"`

	// Filtering
	AllowDomains []string `yaml:"allow-domains"`
	BlockDomains []string `yaml:"block-domains"`

	// Performance
	Workers   int  `yaml:"workers"`
	BatchSize int  `yaml:"batch-size"`
	Streaming bool `yaml:"streaming"`

	// Streaming options
	StreamingFlushInterval string `yaml:"streaming-flush-interval"`
	StreamingMaxBuffer     int    `yaml:"streaming-max-buffer"`

	// Profiles
	Profiles map[string]Profile `yaml:"profiles"`
}

// Profile represents a named configuration profile
type Profile struct {
	Mode             string   `yaml:"mode"`
	FuzzyMode        bool     `yaml:"fuzzy"`
	FuzzyPatterns    []string `yaml:"fuzzy-patterns"`
	IgnoreParams     []string `yaml:"ignore-params"`
	IgnoreExtensions []string `yaml:"ignore-extensions"`
	AllowDomains     []string `yaml:"allow-domains"`
	BlockDomains     []string `yaml:"block-domains"`
	Workers          int      `yaml:"workers"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *File {
	return &File{
		Mode:           "url",
		IgnoreFragment: true,
		TrimSpaces:     true,
		OutputFormat:   "text",
		Workers:        1,
		BatchSize:      1000,
		FuzzyPatterns:  []string{"numeric"},
		StreamingFlushInterval: "5s",
		StreamingMaxBuffer:     10000,
		Profiles: map[string]Profile{
			"aggressive": {
				Mode:          "url",
				FuzzyMode:     true,
				FuzzyPatterns: []string{"numeric", "uuid", "hash", "token"},
				IgnoreParams:  []string{"utm_source", "utm_medium", "utm_campaign", "fbclid", "gclid"},
				Workers:       4,
			},
			"conservative": {
				Mode:      "url",
				FuzzyMode: false,
				Workers:   1,
			},
			"bubbounty": {
				Mode:             "url",
				FuzzyMode:        true,
				FuzzyPatterns:    []string{"numeric", "uuid"},
				IgnoreExtensions: []string{"jpg", "jpeg", "png", "gif", "css", "js", "woff", "woff2", "ttf", "svg", "ico"},
				IgnoreParams:     []string{"utm_source", "utm_medium", "utm_campaign", "ref", "fbclid", "gclid"},
				Workers:          4,
			},
		},
	}
}

// Load loads configuration from a file
func Load(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// LoadWithProfile loads configuration and applies a profile
func LoadWithProfile(path, profileName string) (*File, error) {
	config, err := Load(path)
	if err != nil {
		return nil, err
	}

	if profileName != "" {
		if err := config.ApplyProfile(profileName); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// ApplyProfile applies a named profile to the configuration
func (c *File) ApplyProfile(name string) error {
	profile, ok := c.Profiles[name]
	if !ok {
		return fmt.Errorf("profile not found: %s", name)
	}

	// Apply profile settings (profile overrides base config)
	if profile.Mode != "" {
		c.Mode = profile.Mode
	}
	if profile.FuzzyMode {
		c.FuzzyMode = profile.FuzzyMode
	}
	if len(profile.FuzzyPatterns) > 0 {
		c.FuzzyPatterns = profile.FuzzyPatterns
	}
	if len(profile.IgnoreParams) > 0 {
		c.IgnoreParams = profile.IgnoreParams
	}
	if len(profile.IgnoreExtensions) > 0 {
		c.IgnoreExtensions = profile.IgnoreExtensions
	}
	if len(profile.AllowDomains) > 0 {
		c.AllowDomains = profile.AllowDomains
	}
	if len(profile.BlockDomains) > 0 {
		c.BlockDomains = profile.BlockDomains
	}
	if profile.Workers > 0 {
		c.Workers = profile.Workers
	}

	return nil
}

// Save saves configuration to a file
func (c *File) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultConfigPath returns the default config file path
func GetDefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "dupdurl", "config.yml")
}

// LoadOrDefault tries to load config from default path, returns default if not found
func LoadOrDefault() *File {
	path := GetDefaultConfigPath()
	if path == "" {
		return DefaultConfig()
	}

	config, err := Load(path)
	if err != nil {
		return DefaultConfig()
	}

	return config
}
