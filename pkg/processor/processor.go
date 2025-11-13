package processor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/lcalzada-xor/dupdurl/pkg/deduplicator"
	"github.com/lcalzada-xor/dupdurl/pkg/normalizer"
	"github.com/lcalzada-xor/dupdurl/pkg/stats"
)

const (
	defaultBufferSize = 64 * 1024
	maxLineLength     = 10 * 1024 * 1024
)

// Config holds processor configuration
type Config struct {
	Normalizer *normalizer.Config
	Workers    int
	BatchSize  int
	Verbose    bool
}

// NewConfig creates a default processor configuration
func NewConfig() *Config {
	return &Config{
		Normalizer: normalizer.NewConfig(),
		Workers:    runtime.NumCPU(),
		BatchSize:  1000,
		Verbose:    false,
	}
}

// Processor handles the main URL processing pipeline
type Processor struct {
	config *Config
	stats  *stats.Statistics
	dedup  *deduplicator.Deduplicator
}

// New creates a new Processor instance
func New(config *Config) *Processor {
	st := stats.NewStatistics()
	return &Processor{
		config: config,
		stats:  st,
		dedup:  deduplicator.New(st),
	}
}

// Process reads URLs from input and returns deduplicated entries
func (p *Processor) Process(input io.Reader) ([]deduplicator.Entry, error) {
	if p.config.Workers > 1 {
		return p.processParallel(input)
	}
	return p.processSequential(input)
}

// processSequential processes URLs sequentially (original behavior)
func (p *Processor) processSequential(input io.Reader) ([]deduplicator.Entry, error) {
	scanner := bufio.NewScanner(input)
	buf := make([]byte, 0, defaultBufferSize)
	scanner.Buffer(buf, maxLineLength)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		p.stats.TotalProcessed++

		if p.config.Normalizer.TrimSpaces && strings.TrimSpace(line) == "" {
			continue
		}

		// Create dedup key (without parameter values for comparison)
		key, err := p.config.Normalizer.CreateDedupKey(line)
		if err != nil {
			p.handleError(lineNum, line, err)
			continue
		}

		// Get normalized URL with values preserved
		normalizedURL, err := p.config.Normalizer.NormalizeURL(line)
		if err != nil {
			continue
		}

		// Add to deduplicator
		p.dedup.Add(key, normalizedURL)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	p.stats.Finish()
	return p.dedup.GetEntries(), nil
}

// processedURL represents a URL that has been processed
type processedURL struct {
	lineNum       int
	originalLine  string
	dedupKey      string
	normalizedURL string
	err           error
}

// processParallel processes URLs in parallel using worker pool
func (p *Processor) processParallel(input io.Reader) ([]deduplicator.Entry, error) {
	jobs := make(chan string, p.config.BatchSize)
	results := make(chan processedURL, p.config.BatchSize)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < p.config.Workers; i++ {
		wg.Add(1)
		go p.worker(&wg, jobs, results)
	}

	// Start result collector
	done := make(chan struct{})
	go p.collector(results, done)

	// Read and send jobs
	scanner := bufio.NewScanner(input)
	buf := make([]byte, 0, defaultBufferSize)
	scanner.Buffer(buf, maxLineLength)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		p.stats.TotalProcessed++

		if p.config.Normalizer.TrimSpaces && strings.TrimSpace(line) == "" {
			continue
		}

		jobs <- line
	}

	close(jobs)
	wg.Wait()
	close(results)
	<-done

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	p.stats.Finish()
	return p.dedup.GetEntries(), nil
}

// worker processes URLs from the jobs channel
func (p *Processor) worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- processedURL) {
	defer wg.Done()

	lineNum := 0
	for line := range jobs {
		lineNum++

		// Create dedup key
		key, err := p.config.Normalizer.CreateDedupKey(line)
		if err != nil {
			results <- processedURL{lineNum: lineNum, originalLine: line, err: err}
			continue
		}

		// Get normalized URL
		normalizedURL, err := p.config.Normalizer.NormalizeURL(line)
		if err != nil {
			results <- processedURL{lineNum: lineNum, originalLine: line, err: err}
			continue
		}

		results <- processedURL{
			lineNum:       lineNum,
			originalLine:  line,
			dedupKey:      key,
			normalizedURL: normalizedURL,
		}
	}
}

// collector collects results from workers
func (p *Processor) collector(results <-chan processedURL, done chan<- struct{}) {
	// Need mutex for parallel access to deduplicator
	var mu sync.Mutex

	for result := range results {
		if result.err != nil {
			p.handleError(result.lineNum, result.originalLine, result.err)
			continue
		}

		mu.Lock()
		p.dedup.Add(result.dedupKey, result.normalizedURL)
		mu.Unlock()
	}

	done <- struct{}{}
}

// handleError handles processing errors
func (p *Processor) handleError(lineNum int, line string, err error) {
	if p.config.Verbose && line != "" {
		fmt.Fprintf(os.Stderr, "Line %d: %v - %s\n", lineNum, err, line)
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "parse error") {
		p.stats.ParseErrors++
	} else if strings.Contains(errMsg, "ignored extension") ||
		strings.Contains(errMsg, "blacklist") ||
		strings.Contains(errMsg, "whitelist") ||
		strings.Contains(errMsg, "domain") {
		p.stats.Filtered++
	}
}

// GetStatistics returns the processor statistics
func (p *Processor) GetStatistics() *stats.Statistics {
	return p.stats
}
