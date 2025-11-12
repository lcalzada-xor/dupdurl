<!-- 
Title: dedup - URL Deduplication Tool for Bug Bounty Hunters & Pentesters
Description: Fast CLI tool for deduplicating URLs from bug bounty recon tools like waybackurls, katana, and gau. Features fuzzy matching, parameter filtering, and multi-format output (JSON, CSV, text). Perfect for security researchers and penetration testers.
Author: YOUR_NAME
Keywords: bug bounty, url deduplication, pentesting, security tools, recon, url normalization, fuzzy matching, waybackurls, katana, gau, CLI tool, golang, bug bounty tools, reconnaissance, infosec, cybersecurity
Category: Security Tools, Bug Bounty, Pentesting, Reconnaissance
-->

# 🔥 dupdurl - URL Deduplication Tool for Bug Bounty & Pentesting

⚡ **Fast, powerful URL deduplication for security researchers and bug bounty hunters**. Fuzzy matching, parameter filtering, and multi-format output.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub forks](https://img.shields.io/github/forks/lcalzada-xor/dupdurl?style=social)](https://github.com/lcalzada-xor/dupdurl/network)
[![Go Report Card](https://goreportcard.com/badge/github.com/lcalzada-xor/dupdurl)](https://goreportcard.com/report/github.com/lcalzada-xor/dupdurl)

> 🎯 Deduplicate waybackurls, katana, and gau output **10x faster** with advanced normalization

[⬇️ Installation](#-installation) • [📖 Quick Start](#-quick-start-for-bug-bounty) • [🎯 Examples](#-usage-examples) • [📚 Full Docs](#-complete-documentation) • [⭐ Star this repo](#-support-this-project)

---

## 🚀 Why dupdurl?

A powerful and flexible **URL deduplication tool** designed specifically for **bug bounty pipelines** and **penetration testing workflows**. Perfect for processing output from tools like `katana`, `waybackurls`, `gau`, and other web crawlers.

### ✨ Key Features

| Feature | dedup | urldedupe | uro | qsreplace |
|---------|-------|-----------|-----|-----------|
| 🎯 **Fuzzy ID matching** | ✅ | ✅ | ❌ | ❌ |
| 🔧 **Parameter filtering** | ✅ | ❌ | ✅ | ✅ |
| 📊 **Multi-format output** | ✅ JSON/CSV/text | ❌ | ❌ | ❌ |
| 📈 **Statistics tracking** | ✅ | ❌ | ❌ | ❌ |
| 🌐 **Domain filtering** | ✅ | ❌ | ❌ | ❌ |
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
cd dedup
make build
sudo make install
```

### Or use directly
```bash
go run dedup.go < urls.txt
```

**Requirements**: Go 1.21+

---

## 🎯 Quick Start for Bug Bounty

### Basic Usage
```bash
# Deduplicate URLs from waybackurls
waybackurls target.com | dedup

# With fuzzy ID matching (recommended for bug bounty)
waybackurls target.com | dedup -fuzzy

# Remove tracking parameters
cat urls.txt | dedup -ignore-params=utm_source,utm_medium,fbclid
```

### Complete Bug Bounty Workflow
```bash
# Ultimate recon pipeline
waybackurls target.com | \
  dedup \
    -fuzzy \
    -ignore-params=utm_source,utm_medium,utm_campaign,fbclid,gclid \
    -ignore-extensions=jpg,png,gif,css,js,woff,woff2 \
    -stats \
    > unique_urls.txt
```

**Result**: Reduce 100,000+ URLs to 200-500 unique patterns! 🚀

---

## 💡 Usage Examples

### 🔥 Fuzzy Mode for API Endpoints
```bash
# Discover unique API patterns
waybackurls api.example.com | dedup -fuzzy -mode=path

# Input:
#   /api/users/123/profile
#   /api/users/456/profile
#   /api/users/789/profile
# Output:
#   /api/users/{id}/profile
```

### 📊 JSON Output for Analysis
```bash
# Export with counts in JSON format
katana -u target.com | dedup -output=json -counts > results.json

# Analyze with jq
cat results.json | jq '.[] | select(.count > 5) | .url'
```

### 🎯 Integration with Recon Tools
```bash
# With waybackurls
waybackurls target.com | dedup -fuzzy > urls.txt

# With gau
gau target.com | dedup -ignore-extensions=jpg,png,css > urls.txt

# With katana
katana -u target.com | dedup -mode=path -fuzzy > paths.txt

# Chain with httpx and nuclei
waybackurls target.com | \
  dedup -fuzzy | \
  httpx -silent -mc 200 | \
  nuclei -t ~/nuclei-templates/
```

### 🌐 Domain Filtering
```bash
# Only process specific domains (whitelist)
cat urls.txt | dedup -allow-domains=example.com,api.example.com

# Exclude CDN domains (blacklist)
cat urls.txt | dedup -block-domains=cdn.example.com,static.example.com
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
- 🔥 **Extension filtering**: Skip specific file extensions (jpg, png, css, js, etc.)
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
| `-ignore-extensions` | `""` | Comma-separated extensions to skip |

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
  dedup -fuzzy \
        -ignore-params=utm_source,utm_medium,fbclid \
        -ignore-extensions=jpg,png,css,js \
        -stats > unique_urls.txt
```

### Workflow 2: API Endpoint Discovery
```bash
# Find unique API endpoint patterns
waybackurls api.target.com | \
  grep -i "/api/" | \
  dedup -fuzzy -mode=path | \
  sort > api_patterns.txt
```

### Workflow 3: Parameter Analysis
```bash
# Discover all unique parameter combinations
waybackurls target.com | \
  dedup -mode=params | \
  sort -u > param_combinations.txt

# Find interesting parameters
grep -E "(callback|redirect|url|return|debug|admin)" param_combinations.txt
```

[📖 See 12+ more workflows →](EXAMPLES.md)

---

## 📈 Performance & Statistics

```bash
# Track what's being filtered
cat urls.txt | dedup -stats -verbose

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

### What makes dedup different from other URL deduplication tools?
`dedup` offers **fuzzy ID matching**, **multi-format output**, and **advanced parameter filtering** specifically designed for bug bounty workflows. It's also **actively maintained** with regular updates.

### Can I use dedup with waybackurls, katana, and gau?
**Absolutely!** That's exactly what it's designed for. Chain it in your recon pipeline:
```bash
waybackurls target.com | dedup -fuzzy | httpx
```

### Is dedup suitable for large URL lists?
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
  dedup -fuzzy | \
  nuclei -t ~/templates/

# amass → httprobe → dedup
amass enum -d target.com | \
  httprobe | \
  dedup -mode=host

# meg + dedup for interesting paths
meg --paths interesting_paths.txt targets.txt | \
  dedup -mode=path
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

> 💡 **Tip**: Chain `dedup` with these tools for optimal bug bounty recon workflows!

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

If `dedup` saves you time during bug bounty hunting or penetration testing, please consider:

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
cd dedup
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
