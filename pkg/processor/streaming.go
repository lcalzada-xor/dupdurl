package processor

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/lcalzada-xor/dupdurl/pkg/deduplicator"
	"github.com/lcalzada-xor/dupdurl/pkg/output"
	"github.com/lcalzada-xor/dupdurl/pkg/stats"
)

// StreamingConfig holds streaming processor configuration
type StreamingConfig struct {
	*Config
	FlushInterval time.Duration // Flush every N seconds
	MaxBuffer     int           // Max entries before forced flush
	Output        output.Formatter
	OutputWriter  io.Writer
}

// NewStreamingConfig creates a default streaming configuration
func NewStreamingConfig() *StreamingConfig {
	return &StreamingConfig{
		Config:        NewConfig(),
		FlushInterval: 5 * time.Second,
		MaxBuffer:     10000,
	}
}

// StreamingProcessor handles streaming URL processing with periodic flushes
type StreamingProcessor struct {
	config *StreamingConfig
	stats  *stats.Statistics
	mu     sync.Mutex
}

// NewStreaming creates a new StreamingProcessor instance
func NewStreaming(config *StreamingConfig) *StreamingProcessor {
	return &StreamingProcessor{
		config: config,
		stats:  stats.NewStatistics(),
	}
}

// ProcessStreaming processes URLs in streaming mode with periodic flushes
// This allows processing infinite datasets without loading everything in memory
func (sp *StreamingProcessor) ProcessStreaming(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	buf := make([]byte, 0, defaultBufferSize)
	scanner.Buffer(buf, maxLineLength)

	// Create temporary deduplicator for current window
	dedup := deduplicator.New(sp.stats)

	// Setup periodic flush ticker
	ticker := time.NewTicker(sp.config.FlushInterval)
	defer ticker.Stop()

	// Channel for flush signals
	flushChan := make(chan struct{}, 1)
	done := make(chan struct{})

	// Goroutine to handle periodic flushes
	go func() {
		for {
			select {
			case <-ticker.C:
				flushChan <- struct{}{}
			case <-done:
				return
			}
		}
	}()

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		sp.stats.TotalProcessed++

		if sp.config.Normalizer.TrimSpaces && strings.TrimSpace(line) == "" {
			continue
		}

		// Create dedup key
		key, err := sp.config.Normalizer.CreateDedupKey(line)
		if err != nil {
			sp.handleError(lineNum, line, err)
			continue
		}

		// Get normalized URL
		normalizedURL, err := sp.config.Normalizer.NormalizeURL(line)
		if err != nil {
			continue
		}

		// Add to current window
		dedup.Add(key, normalizedURL)

		// Check if we need to flush due to buffer size
		if dedup.Count() >= sp.config.MaxBuffer {
			if err := sp.flush(dedup); err != nil {
				return err
			}
			dedup = deduplicator.New(sp.stats) // Reset window
		}

		// Check for periodic flush signal (non-blocking)
		select {
		case <-flushChan:
			if dedup.Count() > 0 {
				if err := sp.flush(dedup); err != nil {
					return err
				}
				dedup = deduplicator.New(sp.stats) // Reset window
			}
		default:
			// Continue processing
		}
	}

	// Final flush of remaining entries
	close(done)
	if dedup.Count() > 0 {
		if err := sp.flush(dedup); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}

// flush writes current buffer to output
func (sp *StreamingProcessor) flush(dedup *deduplicator.Deduplicator) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	entries := dedup.GetEntries()
	if len(entries) == 0 {
		return nil
	}

	if sp.config.Output != nil && sp.config.OutputWriter != nil {
		return sp.config.Output.Format(entries, sp.config.OutputWriter)
	}

	return nil
}

// handleError handles processing errors in streaming mode
func (sp *StreamingProcessor) handleError(lineNum int, line string, err error) {
	if sp.config.Verbose && line != "" {
		fmt.Fprintf(sp.config.OutputWriter, "Line %d: %v - %s\n", lineNum, err, line)
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "parse error") {
		sp.stats.ParseErrors++
	} else if strings.Contains(errMsg, "ignored extension") ||
		strings.Contains(errMsg, "blacklist") ||
		strings.Contains(errMsg, "whitelist") ||
		strings.Contains(errMsg, "domain") {
		sp.stats.Filtered++
	}
}

// GetStatistics returns the streaming processor statistics
func (sp *StreamingProcessor) GetStatistics() *stats.Statistics {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	return sp.stats
}
