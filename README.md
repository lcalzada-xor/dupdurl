<!--
Title: dupdurl - URL Deduplication Tool for Bug Bounty & Pentesting
Description: Fast CLI tool for deduplicating URLs from bug bounty recon tools like waybackurls, katana, and gau. Features fuzzy matching, parameter filtering, and multi-format output (JSON, CSV, text).
Keywords: bug bounty, url deduplication, pentesting, security tools, recon, url normalization, fuzzy matching, waybackurls, katana, gau, CLI tool, golang
-->

# dupdurl - URL Deduplication Tool

Fast, powerful URL deduplication for bug bounty and penetration testing. Fuzzy matching, parameter filtering, and multi-format output.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**Features**: Fuzzy ID matching • Parameter filtering • Extension filtering • Multi-format output (JSON/CSV) • Streaming mode • Diff mode • Scope checking • Parallel processing

[Installation](#installation) • [Quick Start](#quick-start) • [Modes](#deduplication-modes) • [Flags](#common-flags) • [Advanced](#advanced-features) • [Examples](#examples)

---

## Installation

```bash
# Quick install
go install github.com/lcalzada-xor/dupdurl@latest

# Or build from source
git clone https://github.com/lcalzada-xor/dupdurl.git
cd dupdurl
make build
sudo make install
```

**Requirements**: Go 1.21+

---

## Quick Start

```bash
# Basic deduplication
waybackurls target.com | dupdurl

# With fuzzy matching (recommended)
waybackurls target.com | dupdurl -f

# Complete bug bounty workflow
waybackurls target.com | dupdurl -f -ie jpg,png,css -s > unique_urls.txt
```

---

## Deduplication Modes

The `-m` flag controls what part of the URL is used for deduplication:

| Mode | Compares | Use Case | Example |
|------|----------|----------|---------|
| **url** | Full URL (default) | General deduplication | `dupdurl` |
| **path** | Domain + path only | Find unique endpoints | `dupdurl -m path -f` |
| **host** | Domain/subdomain only | Enumerate subdomains | `dupdurl -m host` |
| **params** | Parameter names only | Discover param combos | `dupdurl -m params` |
| **raw** | Exact string match | No normalization | `dupdurl -m raw` |

### Examples

**url mode** (default) - Full URL deduplication:
```bash
# Input:
https://example.com/api/users?id=123&sort=asc
http://www.example.com/api/users?sort=asc&id=123  # Duplicate (normalized)
https://example.com/api/users?id=456

# Output:
https://example.com/api/users?id=123&sort=asc
https://example.com/api/users?id=456
```

**path mode** - Find unique endpoints (RECOMMENDED for APIs):
```bash
# Input:
https://example.com/api/users?id=1
https://example.com/api/users?id=2  # Duplicate (same path)
https://example.com/api/products?id=5

# Output:
example.com/api/users
example.com/api/products

# Best with fuzzy mode:
waybackurls target.com | dupdurl -m path -f
```

**host mode** - Extract unique domains:
```bash
# Input:
https://api.example.com/v1/users
https://api.example.com/v2/products  # Duplicate (same host)
https://web.example.com/home

# Output:
api.example.com
web.example.com
```

**params mode** - Find parameter combinations:
```bash
# Input:
/search?q=test&page=1&sort=asc
/search?q=hello&page=2&sort=desc  # Same params (different values)
/search?q=world&page=1            # Different combo (missing 'sort')

# Output:
page,q,sort
page,q
```

---

## Common Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--mode <mode>` | `-m` | Mode: url, path, host, params, raw (default: url) |
| `--fuzzy` | `-f` | Replace IDs with {id} placeholder |
| `--fuzzy-patterns <list>` | `-fp` | Patterns: numeric, uuid, hash, token (default: numeric) |
| `--ignore-params <list>` | `-ip` | Remove specific params (e.g., utm_source,fbclid) |
| `--sort-params` | `-sp` | Sort parameters alphabetically |
| `--ignore-extensions <ext>` | `-ie` | Skip these extensions (e.g., jpg,png,css) |
| `--filter-extensions <ext>` | `-fe` | Only process these extensions (e.g., js,html,php) |
| `--allow-domains <list>` | `-ad` | Only these domains (whitelist) |
| `--block-domains <list>` | `-bd` | Skip these domains (blacklist) |
| `--output <format>` | `-o` | Format: text, json, csv (default: text) |
| `--counts` | `-c` | Show occurrence counts |
| `--stats` | `-s` | Show statistics |
| `--verbose` | `-v` | Show errors and warnings |
| `--workers <n>` | `-w` | Parallel workers (default: 1, 0=auto) |

**Tip**: Use short flags for faster workflows: `dupdurl -f -s -o json` instead of `dupdurl --fuzzy --stats --output=json`

---

## Advanced Features

### Streaming Mode
Process infinite datasets without memory limits:
```bash
cat huge_urls.txt | dupdurl --stream --stream-interval=5s
tail -f access.log | dupdurl --stream -f
```

### Diff Mode
Compare scans and track changes:
```bash
# Save baseline
waybackurls target.com | dupdurl --save-baseline=day1.json

# Later, compare
waybackurls target.com | dupdurl --diff=day1.json
```

### Scope Checking
Filter by domain patterns with wildcards:
```bash
# Create scope file
cat > scope.txt << EOF
*.example.com
!dev.example.com
EOF

# Filter in-scope only
dupdurl --scope=scope.txt < urls.txt

# Show stats
dupdurl --scope=scope.txt --scope-stats < urls.txt
```

### Config Files
Save your preferred settings:
```bash
cat > ~/.config/dupdurl/config.yml << EOF
mode: url
fuzzy: true
ignore-params: [utm_source, fbclid]
workers: 4
EOF

dupdurl --config ~/.config/dupdurl/config.yml < urls.txt
```

---

## Examples

### Discover API Endpoints
```bash
# Find unique API patterns
waybackurls api.example.com | dupdurl -m path -f

# Input:
#   /api/users/123/profile
#   /api/users/456/profile
# Output:
#   api.example.com/api/users/{id}/profile
```

### Filter by Extensions
```bash
# Only JavaScript files
waybackurls target.com | dupdurl -fe js

# Multiple extensions
waybackurls target.com | dupdurl -fe js,json,php,html

# Exclude images/styles
waybackurls target.com | dupdurl -ie jpg,png,css,woff
```

### JSON Output
```bash
# Export with counts
katana -u target.com | dupdurl -o json -c > results.json

# Analyze with jq
cat results.json | jq '.[] | select(.count > 5) | .url'
```

### Integration with Tools
```bash
# Full recon pipeline
waybackurls target.com | dupdurl -f -ie jpg,png,css | httpx -silent | nuclei

# With gau
gau target.com | dupdurl -f -s > urls.txt

# With katana
katana -u target.com | dupdurl -m path -f > paths.txt
```

### Bug Bounty Workflow
```bash
#!/bin/bash
TARGET="example.com"

# Collect URLs
waybackurls $TARGET > wayback.txt
gau $TARGET > gau.txt

# Deduplicate and filter
cat wayback.txt gau.txt | \
  dupdurl -f -ie jpg,png,css -ip utm_source,fbclid -s > unique.txt
```

---

## FAQ

**What makes dupdurl different?**
Fuzzy ID matching, multi-format output, and advanced filtering designed specifically for bug bounty workflows.

**Can I use it with waybackurls, katana, and gau?**
Yes! That's exactly what it's designed for:
```bash
waybackurls target.com | dupdurl -f | httpx
```

**Is it suitable for large URL lists?**
Yes! Tested with 1M+ URLs. Processes 100K URLs in ~2 seconds.

**How does fuzzy matching work?**
Replaces numeric IDs with `{id}` placeholder. Example: `/users/123/profile` and `/users/456/profile` → `/users/{id}/profile`

**What output formats are supported?**
Text (default), JSON with counts, and CSV.

---

## Related Tools

**URL Collection**: [waybackurls](https://github.com/tomnomnom/waybackurls) • [katana](https://github.com/projectdiscovery/katana) • [gau](https://github.com/lc/gau) • [hakrawler](https://github.com/hakluke/hakrawler)

**Similar Tools**: [urldedupe](https://github.com/ameenmaali/urldedupe) • [uro](https://github.com/s0md3v/uro)

**Workflow Tools**: [httpx](https://github.com/projectdiscovery/httpx) • [nuclei](https://github.com/projectdiscovery/nuclei) • [ffuf](https://github.com/ffuf/ffuf) • [subfinder](https://github.com/projectdiscovery/subfinder)

---

## Contributing

Contributions welcome! Bug reports, feature suggestions, documentation improvements, and code contributions all appreciated.

```bash
git clone https://github.com/lcalzada-xor/dupdurl.git
cd dupdurl
make build && make test
```

---

## License

MIT License - Free to use, modify, and distribute.

---

## Support

If dupdurl saves you time, please star the repository and share with other security researchers!

[⭐ Star](https://github.com/lcalzada-xor/dupdurl) • [🐛 Report Bug](https://github.com/lcalzada-xor/dupdurl/issues) • [💡 Request Feature](https://github.com/lcalzada-xor/dupdurl/issues)
