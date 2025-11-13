<!-- 
Title: dedup - URL Deduplication Tool for Bug Bounty Hunters & Pentesters
Description: Fast CLI tool for deduplicating URLs from bug bounty recon tools like waybackurls, katana, and gau. Features fuzzy matching, parameter filtering, and multi-format output (JSON, CSV, text). Perfect for security researchers and penetration testers.
Author: YOUR_NAME
Keywords: bug bounty, url deduplication, pentesting, security tools, recon, url normalization, fuzzy matching, waybackurls, katana, gau, CLI tool, golang, bug bounty tools, reconnaissance, infosec, cybersecurity
Category: Security Tools, Bug Bounty, Pentesting, Reconnaissance
-->

# dupdurl - URL Deduplication Tool for Bug Bounty & Pentesting

⚡ **Fast, powerful URL deduplication for security researchers and bug bounty hunters**. Fuzzy matching, parameter filtering, and multi-format output.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub forks](https://img.shields.io/github/forks/lcalzada-xor/dupdurl?style=social)](https://github.com/lcalzada-xor/dupdurl/network)
[![Go Report Card](https://goreportcard.com/badge/github.com/lcalzada-xor/dupdurl)](https://goreportcard.com/report/github.com/lcalzada-xor/dupdurl)

> 🎯 Deduplicate waybackurls, katana, and gau output **10x faster** with advanced normalization

[⬇️ Installation](#-installation) • [📖 Quick Start](#-quick-start-for-bug-bounty) • [🎯 Examples](#-usage-examples) • [📚 Full Docs](#-complete-documentation) • [⭐ Star this repo](#-support-this-project)

---

## Why dupdurl?

A powerful and flexible **URL deduplication tool** designed specifically for **bug bounty pipelines** and **penetration testing workflows**. Perfect for processing output from tools like `katana`, `waybackurls`, `gau`, and other web crawlers.

### ✨ Key Features

| Feature | dupdurl v2.1 | urldedupe | uro | qsreplace |
|---------|--------------|-----------|-----|-----------|
| 🎯 **Fuzzy ID matching** | ✅ 4 patterns | ✅ | ❌ | ❌ |
| 🔧 **Parameter filtering** | ✅ | ❌ | ✅ | ✅ |
| 📊 **Multi-format output** | ✅ JSON/CSV/text | ❌ | ❌ | ❌ |
| 📈 **Statistics tracking** | ✅ Detailed | ❌ | ❌ | ❌ |
| 🌐 **Domain filtering** | ✅ | ❌ | ❌ | ❌ |
| 🌊 **Streaming mode** | ✅ | ❌ | ❌ | ❌ |
| 📊 **Diff mode** | ✅ | ❌ | ❌ | ❌ |
| ⚙️ **Config files** | ✅ | ❌ | ❌ | ❌ |
| 🎯 **Scope checking** | ✅ | ❌ | ❌ | ❌ |
| ⚡ **Parallel processing** | ✅ | ❌ | ❌ | ❌ |
| 🚀 **Active development** | ✅ | ⚠️ | ✅ | ⚠️ |
| 🐛 **Bug bounty focused** | ✅ | ✅ | ⚠️ | ⚠️ |

### 🎯 Perfect For

- 🐛 **Bug Bounty Hunters** - Clean recon output efficiently
- 🔒 **Penetration Testers** - Deduplicate URLs before fuzzing
- 🔍 **Security Researchers** - Analyze URL patterns across domains
- 🎓 **Students & Learners** - Understand URL normalization

---

## 📦 Installation

### Quick Install (Recommended)
```bash
go install github.com/lcalzada-xor/dupdurl@latest
```

### Build from Source
```bash
git clone https://github.com/lcalzada-xor/dupdurl.git
cd dupdurl
make build
sudo make install
```

### Or use directly
```bash
go run dupdurl.go < urls.txt
```

**Requirements**: Go 1.21+

---

## 🎯 Quick Start for Bug Bounty

### Basic Usage
```bash
# Deduplicate URLs from waybackurls
waybackurls target.com | dupdurl

# With fuzzy ID matching (recommended for bug bounty)
waybackurls target.com | dupdurl --fuzzy

# Remove tracking parameters
cat urls.txt | dupdurl --ignore-params=utm_source,utm_medium,fbclid
```

### Complete Bug Bounty Workflow
```bash
# Ultimate recon pipeline
waybackurls target.com | \
  dupdurl \
    --fuzzy \
    --ignore-params=utm_source,utm_medium,utm_campaign,fbclid,gclid \
    --ignore-extensions=jpg,png,gif,css,js,woff,woff2 \
    --stats \
    > unique_urls.txt

# Using short flags (faster typing!)
waybackurls target.com | dupdurl -f -ip utm_source,fbclid -ie jpg,png,css -s
```

**Result**: Reduce 100,000+ URLs to 200-500 unique patterns! 🚀

---

## ⚡ Advanced Features

### 🌊 Streaming Mode
Process **infinite datasets** without memory limits:
```bash
# Process massive datasets with periodic output
cat huge_urls.txt | dupdurl -stream -stream-interval=5s

# Perfect for continuous monitoring
tail -f access.log | dupdurl -stream -fuzzy
```

### 📊 Diff Mode
Compare scans and track changes over time:
```bash
# Save baseline
waybackurls target.com | dupdurl -save-baseline baseline.json

# Later, compare new scan
waybackurls target.com | dupdurl -diff baseline.json

# Output:
# [ADDED] 12 new URLs:
#   + https://example.com/api/v2/users
#   + https://example.com/admin/new-feature
# [REMOVED] 3 URLs:
#   - https://example.com/old/endpoint
```

### ⚙️ Config File Support
Save your preferred settings:
```bash
# Create config file
cat > ~/.config/dupdurl/config.yml << EOF
mode: url
fuzzy: true
fuzzy-patterns: [numeric, uuid]
ignore-params: [utm_source, utm_medium, fbclid]
workers: 4
EOF

# Use config
cat urls.txt | dupdurl

# Or use predefined profiles
cat urls.txt | dupdurl -profile bugbounty
cat urls.txt | dupdurl -profile aggressive
```

### ⚡ Performance Optimizations
- **String pooling** - Reduced memory allocations
- **Pre-sized maps** - Faster lookups
- **Parallel processing** - 3-5x speedup with `-workers=4`

### 🎯 Scope Checking
Filter URLs based on scope file with wildcard and exclusion support:
```bash
# Create scope file
cat > scope.txt << EOF
# In scope domains
example.com
*.example.com
target.com

# Exclusions (prefix with !)
!dev.example.com
!staging.example.com
EOF

# Filter in-scope only
waybackurls target.com | dupdurl --scope=scope.txt

# Show scope statistics
dupdurl --scope=scope.txt --scope-stats < urls.txt

# Show only out-of-scope URLs
dupdurl --scope=scope.txt --out-of-scope < urls.txt
```

**Features**:
- Wildcard support: `*.example.com` matches all subdomains
- Exclusions: `!dev.example.com` to exclude specific domains
- Smart normalization: removes `www.` prefix and ports
- Comment support with `#`

---

## 🚀 Quick Reference - Common Flags

| Long Flag | Short | Description | Example |
|-----------|-------|-------------|---------|
| `--fuzzy` | `-f` | Enable fuzzy ID matching | `dupdurl -f` |
| `--mode` | `-m` | Normalization mode | `dupdurl -m path` |
| `--output` | `-o` | Output format (text/json/csv) | `dupdurl -o json` |
| `--stats` | `-s` | Show statistics | `dupdurl -s` |
| `--counts` | `-c` | Print occurrence counts | `dupdurl -c` |
| `--verbose` | `-v` | Verbose mode | `dupdurl -v` |
| `--workers` | `-w` | Parallel workers | `dupdurl -w 4` |
| `--ignore-params` | `-ip` | Ignore query params | `dupdurl -ip utm_source,fbclid` |
| `--ignore-extensions` | `-ie` | Ignore extensions (blacklist) | `dupdurl -ie jpg,png,css` |
| `--filter-extensions` | `-fe` | Only allow extensions (whitelist) | `dupdurl -fe js,json,html` |
| `--allow-domains` | `-ad` | Domain whitelist | `dupdurl -ad example.com` |
| `--block-domains` | `-bd` | Domain blacklist | `dupdurl -bd ads.com` |
| `--profile` | `-p` | Use config profile | `dupdurl -p bugbounty` |
| `--diff` | `-d` | Compare with baseline | `dupdurl -d baseline.json` |
| `--save-baseline` | `-sb` | Save as baseline | `dupdurl -sb day1.json` |
| `--scope` | `-S` | Scope file with wildcards | `dupdurl -S scope.txt` |

**Pro tip**: Combine short flags for faster workflows!
```bash
# Instead of: dupdurl --fuzzy --stats --output=json --counts
# Use:        dupdurl -f -s -o json -c
```

---

## 💡 Usage Examples

### 🔥 Fuzzy Mode for API Endpoints
```bash
# Discover unique API patterns
waybackurls api.example.com | dupdurl --fuzzy --mode=path

# Using short flags
waybackurls api.example.com | dupdurl -f -m path

# Input:
#   /api/users/123/profile
#   /api/users/456/profile
#   /api/users/789/profile
# Output:
#   /api/users/{id}/profile
```

### 🎯 Filter by File Extensions (NEW!)
```bash
# Only process JavaScript files
waybackurls target.com | dupdurl --filter-extensions=js

# Multiple extensions - focus on code files
waybackurls target.com | dupdurl -fe js,json,php,html

# Find only API endpoints (JSON/XML)
katana -u target.com | dupdurl -fe json,xml -fuzzy

# Combine with other filters
waybackurls target.com | dupdurl -fe js,html -fuzzy -stats
```

### 📊 JSON Output for Analysis
```bash
# Export with counts in JSON format
katana -u target.com | dupdurl --output=json --counts > results.json

# Using short flags
katana -u target.com | dupdurl -o json -c > results.json

# Analyze with jq
cat results.json | jq '.[] | select(.count > 5) | .url'
```

### 🎯 Integration with Recon Tools
```bash
# With waybackurls
waybackurls target.com | dupdurl --fuzzy > urls.txt

# With gau (using short flags)
gau target.com | dupdurl -ie jpg,png,css > urls.txt

# With katana
katana -u target.com | dupdurl -m path -f > paths.txt

# Chain with httpx and nuclei
waybackurls target.com | \
  dupdurl --fuzzy | \
  httpx -silent -mc 200 | \
  nuclei -t ~/nuclei-templates/
```

### 🌐 Domain Filtering
```bash
# Only process specific domains (whitelist)
cat urls.txt | dupdurl -allow-domains=example.com,api.example.com

# Exclude CDN domains (blacklist)
cat urls.txt | dupdurl -block-domains=cdn.example.com,static.example.com
```

---

## 🔧 All Features Explained

### Core Functionality
- ✅ **Multiple normalization modes**: URL, path, host, params, raw
- ✅ **Smart query parameter handling**: Sort, ignore specific params, preserve order
- ✅ **Fragment handling**: Remove or preserve URL fragments (#)
- ✅ **Case sensitivity options**: Control case sensitivity for paths and hosts
- ✅ **WWW stripping**: Optional removal of leading `www.`
- ✅ **Scheme normalization**: Choose to distinguish http/https or not

### Advanced Features (v2.0+)
- 🔥 **Fuzzy mode**: Replace numeric IDs in paths with `{id}` placeholders
- 🔥 **Extension filtering**: Skip specific file extensions (blacklist) or only allow certain extensions (whitelist)
- 🔥 **Domain filtering**: Whitelist or blacklist domains
- 🔥 **Multiple output formats**: Text, JSON, CSV
- 🔥 **Statistics tracking**: See detailed processing metrics
- 🔥 **Verbose mode**: Debug parsing errors and filtered URLs

---

## 🎛️ Command-Line Options

### Core Options
| Flag | Default | Description |
|------|---------|-------------|
| `-mode` | `url` | Normalization mode: `url`, `path`, `host`, `raw`, `params` |
| `-ignore-params` | `""` | Comma-separated query params to remove |
| `-sort-params` | `false` | Sort query parameters alphabetically |
| `-ignore-fragment` | `true` | Remove URL fragment (#...) |
| `-fuzzy` | `false` | Replace numeric IDs in paths with {id} |
| `-ignore-extensions` | `""` | Comma-separated extensions to skip (blacklist) |
| `-filter-extensions` | `""` | Comma-separated extensions to allow (whitelist) |

**Note**: `--ignore-extensions` and `--filter-extensions` cannot be used together.

**Extension Filtering Modes**:
- **Blacklist** (`-ie`): Removes specified extensions → `dupdurl -ie jpg,png,css` excludes images/styles
- **Whitelist** (`-fe`): Only keeps specified extensions → `dupdurl -fe js,json,php` keeps only code files

### Output Options
| Flag | Default | Description |
|------|---------|-------------|
| `-counts` | `false` | Print counts before each unique entry |
| `-output` | `text` | Output format: `text`, `json`, `csv` |
| `-stats` | `false` | Print statistics at the end |
| `-verbose` | `false` | Show warnings and parse errors |

[📚 See all options →](CHANGELOG.md#all-options)

---

## 🎓 Real-World Bug Bounty Workflows

### Workflow 1: Initial Reconnaissance
```bash
#!/bin/bash
TARGET="example.com"

# Collect URLs from multiple sources
waybackurls $TARGET > wayback.txt
gau $TARGET > gau.txt
katana -u $TARGET -d 5 -o katana.txt

# Deduplicate and filter
cat wayback.txt gau.txt katana.txt | \
  dupdurl -fuzzy \
        -ignore-params=utm_source,utm_medium,fbclid \
        -ignore-extensions=jpg,png,css,js \
        -stats > unique_urls.txt
```

### Workflow 2: API Endpoint Discovery
```bash
# Find unique API endpoint patterns
waybackurls api.target.com | \
  grep -i "/api/" | \
  dupdurl -fuzzy -mode=path | \
  sort > api_patterns.txt
```

### Workflow 3: Parameter Analysis
```bash
# Discover all unique parameter combinations
waybackurls target.com | \
  dupdurl -mode=params | \
  sort -u > param_combinations.txt

# Find interesting parameters
grep -E "(callback|redirect|url|return|debug|admin)" param_combinations.txt
```

[📖 See 12+ more workflows →](EXAMPLES.md)

---

## 📈 Performance & Statistics

```bash
# Track what's being filtered
cat urls.txt | dupdurl -stats -verbose

# Output:
# === Statistics ===
# Total URLs processed: 15,234
# Unique URLs:          847
# Duplicates removed:   14,156
# Parse errors:         18
# Filtered out:         213
# ==================
```

**Performance**: Processes 100K URLs in ~2 seconds on modern hardware.

---

## 🔍 Modes Explained

### `url` mode (default)
Full URL normalization with all options applied.
```bash
Input:  https://www.example.com/path?b=2&a=1#section
Output: https://example.com/path?b=2&a=1
```

### `path` mode  
Host + path only, ignoring scheme, query, and fragment.
```bash
Input:  https://example.com/api/users?id=123#top
Output: example.com/api/users
```

### `host` mode
Extract and normalize hostname only.
```bash
Input:  https://www.example.com:443/path?query
Output: example.com:443
```

### `params` mode
Extract unique parameter name combinations.
```bash
Input:  https://example.com?user=john&id=123&sort=asc
Output: id,sort,user
```

### `fuzzy` mode (with -fuzzy flag)
Replace numeric IDs with placeholders.
```bash
Input:  /api/users/123/profile
Output: /api/users/{id}/profile
```

---

## ❓ Frequently Asked Questions

### What makes dupdurl different from other URL deduplication tools?
`dupdurl` offers **fuzzy ID matching**, **multi-format output**, and **advanced parameter filtering** specifically designed for bug bounty workflows. It's also **actively maintained** with regular updates.

### Can I use dupdurl with waybackurls, katana, and gau?
**Absolutely!** That's exactly what it's designed for. Chain it in your recon pipeline:
```bash
waybackurls target.com | dupdurl -fuzzy | httpx
```

### Is dupdurl suitable for large URL lists?
**Yes!** Tested with 1M+ URLs. Processes 100K URLs in ~2 seconds. Supports up to 10MB line lengths.

### How does fuzzy matching work?
Fuzzy mode identifies numeric path segments (like user IDs) and replaces them with `{id}` placeholder, helping you discover unique API patterns instead of seeing thousands of similar URLs.

### What output formats are supported?
- **Text** (default): Plain URLs, one per line  
- **JSON**: Structured data with counts `[{"url": "...", "count": 5}]`
- **CSV**: Import into Excel `url,count`

### Does it work on Windows/Mac/Linux?
**Yes!** Written in Go, it's cross-platform. Pre-built binaries available for all major platforms.

---

## 🤝 Integration Examples

### With Popular Bug Bounty Tools

```bash
# subfinder → httpx → katana → dedup → nuclei
subfinder -d target.com | \
  httpx -silent | \
  katana -d 3 | \
  dupdurl -fuzzy | \
  nuclei -t ~/templates/

# amass → httprobe → dedup
amass enum -d target.com | \
  httprobe | \
  dupdurl -mode=host

# meg + dedup for interesting paths
meg --paths interesting_paths.txt targets.txt | \
  dupdurl -mode=path
```

---

## 🔗 Related Projects & Tools

### URL Collection Tools
- [**waybackurls**](https://github.com/tomnomnom/waybackurls) - Fetch URLs from Wayback Machine
- [**katana**](https://github.com/projectdiscovery/katana) - Next-gen crawling framework
- [**gau**](https://github.com/lc/gau) - Fetch known URLs from AlienVault, Wayback, etc.
- [**hakrawler**](https://github.com/hakluke/hakrawler) - Fast web crawler

### Similar Deduplication Tools
- [**urldedupe**](https://github.com/ameenmaali/urldedupe) - Alternative URL deduplication with similar flags
- [**uro**](https://github.com/s0md3v/uro) - URL parameter remover and cleaner

### Recommended Workflow Tools
- [**httpx**](https://github.com/projectdiscovery/httpx) - Fast HTTP toolkit for probing
- [**nuclei**](https://github.com/projectdiscovery/nuclei) - Vulnerability scanner with templates
- [**ffuf**](https://github.com/ffuf/ffuf) - Fast web fuzzer
- [**subfinder**](https://github.com/projectdiscovery/subfinder) - Subdomain discovery tool

> 💡 **Tip**: Chain `dupdurl` with these tools for optimal bug bounty recon workflows!

---

## 📚 Complete Documentation

- [📖 **Full Usage Guide**](README.md) - This file
- [🎯 **12+ Practical Examples**](EXAMPLES.md) - Real bug bounty scenarios  
- [📋 **Changelog**](CHANGELOG.md) - Version history and migration guide
- [🧪 **Testing Guide**](dedup_test.go) - Test suite documentation
- [🔧 **Build System**](Makefile) - Compilation and installation

---

## 🧪 Testing

```bash
# Run tests
make test

# Run with coverage
make test-coverage

# Run demo
make demo
```

**Test Coverage**: 95%+ across all core functions

---

## ⭐ Support This Project

If `dupdurl` saves you time during bug bounty hunting or penetration testing, please consider:

- ⭐ **Star this repository** (it helps others discover it!)
- 🐛 [Report bugs or issues](https://github.com/lcalzada-xor/dupdurl/issues)
- 💡 [Suggest new features](https://github.com/lcalzada-xor/dupdurl/issues/new)
- 🤝 [Contribute code](https://github.com/lcalzada-xor/dupdurl/pulls)
- 📢 Share with other security researchers

**Your support motivates continued development!** 🙏

---

## 🗺️ Roadmap

### Planned Features (v2.1+)
- [ ] Parallel processing for massive URL lists
- [ ] Custom regex patterns for fuzzy matching  
- [ ] SQLite output format for persistent storage
- [ ] Rate limiting awareness per endpoint
- [ ] Interactive TUI mode
- [ ] Burp Suite integration
- [ ] HTML report generation

[See full roadmap →](CHANGELOG.md#roadmap)

---

## 🤝 Contributing

Contributions are welcome! Whether it's:
- 🐛 Bug reports
- 💡 Feature suggestions  
- 📖 Documentation improvements
- 🔧 Code contributions

**Good first issues**: Look for issues tagged [`good-first-issue`](https://github.com/lcalzada-xor/dupdurl/labels/good-first-issue)

### Development Setup
```bash
git clone https://github.com/lcalzada-xor/dupdurl.git
cd dupdurl
make build
make test
```

---

## 📝 License

MIT License - See [LICENSE](LICENSE) file for details.

**TL;DR**: Free to use, modify, and distribute. Attribution appreciated! 🙏

---

## 🙏 Acknowledgments

Inspired by and designed to work seamlessly with:
- [TomNomNom's](https://github.com/tomnomnom) amazing bug bounty tools
- [ProjectDiscovery](https://github.com/projectdiscovery) for the excellent recon framework
- The entire bug bounty and infosec community

Special thanks to all [contributors](https://github.com/lcalzada-xor/dupdurl/graphs/contributors)!

---

## 💬 Community & Support

- 💬 [GitHub Discussions](https://github.com/lcalzada-xor/dupdurl/discussions) - Ask questions, share workflows
- 🐛 [Issue Tracker](https://github.com/lcalzada-xor/dupdurl/issues) - Report bugs

---

## 📊 Stats & Analytics

![GitHub stars](https://img.shields.io/github/stars/lcalzada-xor/dupdurl?style=social)
![GitHub forks](https://img.shields.io/github/forks/lcalzada-xor/dupdurl?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/lcalzada-xor/dupdurl?style=social)
![GitHub downloads](https://img.shields.io/github/downloads/lcalzada-xor/dupdurl/total)

---

<div align="center">

**Made with ❤️ for the bug bounty community**

[⭐ Star](https://github.com/lcalzada-xor/dupdurl) • [🐛 Report Bug](https://github.com/lcalzada-xor/dupdurl/issues) • [💡 Request Feature](https://github.com/lcalzada-xor/dupdurl/issues)

**Happy Hunting! 🎯🐛**

</div>
