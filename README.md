# dedup - Advanced URL Deduplication Tool

A powerful and flexible URL deduplication tool designed specifically for bug bounty pipelines. Perfect for processing output from tools like `katana`, `waybackurls`, `gau`, and other web crawlers.

## 🚀 Features

### Core Functionality
- ✅ **Multiple normalization modes**: URL, path, host, params, raw
- ✅ **Smart query parameter handling**: Sort, ignore specific params, preserve order
- ✅ **Fragment handling**: Remove or preserve URL fragments
- ✅ **Case sensitivity options**: Control case sensitivity for paths and hosts
- ✅ **WWW stripping**: Optional removal of leading `www.`
- ✅ **Scheme normalization**: Choose to distinguish http/https or not

### Advanced Features
- 🔥 **Fuzzy mode**: Replace numeric IDs in paths with `{id}` placeholders
- 🔥 **Extension filtering**: Skip specific file extensions (jpg, png, css, js, etc.)
- 🔥 **Domain filtering**: Whitelist or blacklist domains
- 🔥 **Multiple output formats**: Text, JSON, CSV
- 🔥 **Statistics tracking**: See detailed processing metrics
- 🔥 **Verbose mode**: Debug parsing errors and filtered URLs

### Bug Fixes from Original
- ✅ Fixed `-sort-params` defaulting to `true` (now defaults to `false`)
- ✅ Query params always normalized (not just when sorting)
- ✅ Better error handling and reporting
- ✅ Consistent URL encoding
- ✅ Improved documentation

## 📦 Installation

### Build from source
```bash
go build -o dedup dedup.go
sudo mv dedup /usr/local/bin/
```

### Quick run (without building)
```bash
go run dedup.go < urls.txt
```

## 🔧 Usage

### Basic Examples

```bash
# Basic deduplication (URL mode - default)
cat urls.txt | dedup

# Deduplicate by path only (ignore scheme, params, fragment)
cat urls.txt | dedup -mode=path

# Deduplicate by host only
cat urls.txt | dedup -mode=host

# Extract unique parameter combinations
cat urls.txt | dedup -mode=params

# Show counts for each unique URL
cat urls.txt | dedup -counts
```

### Query Parameter Handling

```bash
# Sort query parameters alphabetically
cat urls.txt | dedup -sort-params

# Ignore specific tracking parameters
cat urls.txt | dedup -ignore-params=utm_source,utm_medium,fbclid

# Combine: ignore params AND sort remaining ones
cat urls.txt | dedup -ignore-params=utm_source,utm_medium -sort-params
```

### Advanced Normalization

```bash
# Fuzzy mode: Replace numeric IDs with {id}
# /user/12345/profile -> /user/{id}/profile
cat urls.txt | dedup -fuzzy

# Ignore image and static file extensions
cat urls.txt | dedup -ignore-extensions=jpg,png,css,js,woff,woff2

# Keep scheme distinction (http vs https)
cat urls.txt | dedup -keep-scheme

# Case-sensitive mode
cat urls.txt | dedup -case-sensitive

# Path mode with query string included
cat urls.txt | dedup -mode=path -path-include-query
```

### Domain Filtering

```bash
# Only process specific domains (whitelist)
cat urls.txt | dedup -allow-domains=example.com,api.example.com

# Exclude specific domains (blacklist)
cat urls.txt | dedup -block-domains=cdn.example.com,static.example.com

# Combine with other options
cat urls.txt | dedup -allow-domains=example.com -ignore-extensions=jpg,png
```

### Output Formats

```bash
# JSON output
cat urls.txt | dedup -output=json

# CSV output with counts
cat urls.txt | dedup -output=csv

# Text with counts
cat urls.txt | dedup -counts
```

### Statistics and Debugging

```bash
# Show processing statistics
cat urls.txt | dedup -stats

# Verbose mode (show errors and warnings)
cat urls.txt | dedup -verbose

# Combine statistics with verbose
cat urls.txt | dedup -verbose -stats
```

## 🎯 Real-World Bug Bounty Examples

### Example 1: Clean waybackurls output
```bash
# Remove tracking params and get unique paths
waybackurls target.com | dedup -mode=path -ignore-params=utm_source,utm_medium,utm_campaign
```

### Example 2: Find unique parameter names
```bash
# Discover all unique parameter combinations
katana -u target.com | dedup -mode=params | sort -u
```

### Example 3: Fuzzy deduplication for endpoints with IDs
```bash
# Identify unique endpoint patterns
gau target.com | dedup -fuzzy -mode=path
# Output: /api/users/{id}/profile instead of multiple /api/users/123/profile, /api/users/456/profile
```

### Example 4: Focus on main domain, ignore CDNs
```bash
cat urls.txt | dedup -allow-domains=target.com,api.target.com -block-domains=cdn.target.com,static.target.com
```

### Example 5: Complete bug bounty pipeline
```bash
# Comprehensive URL collection and deduplication
echo target.com | \
  waybackurls | \
  dedup \
    -fuzzy \
    -ignore-params=utm_source,utm_medium,utm_campaign,fbclid \
    -ignore-extensions=jpg,png,gif,css,js,woff,woff2,svg \
    -stats \
    -verbose 2>errors.log
```

### Example 6: Compare subdomain structures
```bash
# Get unique paths across all subdomains
cat subdomains.txt | while read sub; do
  waybackurls "$sub"
done | dedup -mode=path -counts | sort -rn | head -100
```

### Example 7: JSON output for further processing
```bash
# Export to JSON for analysis with jq
katana -u target.com | dedup -output=json -counts > urls.json
cat urls.json | jq '.[] | select(.count > 5) | .url'
```

## 📊 Output Examples

### Text Output (default)
```
https://example.com/api/users
https://example.com/login
https://example.com/dashboard
```

### Text Output with Counts
```bash
dedup -counts
```
```
5 https://example.com/api/users
2 https://example.com/login
1 https://example.com/dashboard
```

### JSON Output
```bash
dedup -output=json
```
```json
[
  {
    "url": "https://example.com/api/users",
    "count": 5
  },
  {
    "url": "https://example.com/login",
    "count": 2
  }
]
```

### CSV Output
```bash
dedup -output=csv
```
```csv
url,count
https://example.com/api/users,5
https://example.com/login,2
https://example.com/dashboard,1
```

### Statistics Output
```bash
dedup -stats
```
```
=== Statistics ===
Total URLs processed: 1523
Unique URLs:          347
Duplicates removed:   1156
Parse errors:         15
Filtered out:         5
==================
```

## 🎛️ Command-Line Options

### Core Options
| Flag | Default | Description |
|------|---------|-------------|
| `-mode` | `url` | Normalization mode: `url`, `path`, `host`, `raw`, `params` |
| `-ignore-params` | `""` | Comma-separated query params to remove |
| `-sort-params` | `false` | Sort query parameters alphabetically |
| `-ignore-fragment` | `true` | Remove URL fragment (#...) |
| `-case-sensitive` | `false` | Consider case when comparing |
| `-keep-www` | `false` | Don't strip leading www. from host |
| `-keep-scheme` | `false` | Distinguish between http:// and https:// |
| `-trim` | `true` | Trim surrounding spaces |

### Output Options
| Flag | Default | Description |
|------|---------|-------------|
| `-counts` | `false` | Print counts before each unique entry |
| `-output` | `text` | Output format: `text`, `json`, `csv` |
| `-stats` | `false` | Print statistics at the end (to stderr) |
| `-verbose` | `false` | Show warnings and parse errors (to stderr) |

### Advanced Options
| Flag | Default | Description |
|------|---------|-------------|
| `-fuzzy` | `false` | Replace numeric IDs in paths with {id} |
| `-ignore-extensions` | `""` | Comma-separated extensions to skip |
| `-path-include-query` | `false` | In path mode, include normalized query |
| `-allow-domains` | `""` | Comma-separated whitelist of domains |
| `-block-domains` | `""` | Comma-separated blacklist of domains |

## 🔍 Modes Explained

### `url` mode (default)
Full URL normalization with all options applied.
```
Input:  https://www.example.com/path?b=2&a=1#section
Output: https://example.com/path?b=2&a=1
```

### `path` mode
Host + path only, ignoring scheme, query, and fragment.
```
Input:  https://example.com/api/users?id=123#top
Output: example.com/api/users
```

### `host` mode
Extract and normalize hostname only.
```
Input:  https://www.example.com:443/path?query
Output: example.com:443
```

### `params` mode
Extract unique parameter name combinations.
```
Input:  https://example.com?user=john&id=123&sort=asc
Output: id,sort,user
```

### `raw` mode
Minimal processing, optional case normalization.
```
Input:  https://EXAMPLE.COM/Path
Output: https://example.com/Path  (if not case-sensitive)
```

## 🧪 Testing

Run the test suite:
```bash
go test -v
```

Run specific tests:
```bash
go test -v -run TestNormalizePath
go test -v -run TestFuzzy
```

Run with coverage:
```bash
go test -cover
```

## 🔄 Integration with Popular Tools

### With waybackurls
```bash
waybackurls target.com | dedup -fuzzy -ignore-params=utm_source > unique_urls.txt
```

### With katana
```bash
katana -u target.com -d 3 | dedup -mode=path -stats
```

### With gau
```bash
gau target.com | dedup -ignore-extensions=jpg,png,gif,css,js -counts
```

### With subfinder + httpx
```bash
subfinder -d target.com | httpx -silent | dedup -mode=host
```

### Chain multiple tools
```bash
cat targets.txt | \
  waybackurls | \
  dedup -fuzzy -ignore-params=utm_source,fbclid | \
  httpx -silent -mc 200 | \
  nuclei -t ~/nuclei-templates/
```

## 💡 Tips & Best Practices

1. **Use `-fuzzy` for API endpoints**: Helps identify unique endpoint patterns by replacing IDs
   ```bash
   dedup -fuzzy -mode=path
   ```

2. **Always filter tracking params**: Remove noise from analytics
   ```bash
   dedup -ignore-params=utm_source,utm_medium,utm_campaign,fbclid,gclid
   ```

3. **Use `-stats` during development**: Monitor how many URLs are being filtered
   ```bash
   dedup -stats -verbose 2>debug.log
   ```

4. **Combine with sorting for analysis**:
   ```bash
   dedup -counts | sort -rn | head -50  # Top 50 most common URLs
   ```

5. **Filter static resources early**: Save processing time
   ```bash
   dedup -ignore-extensions=jpg,jpeg,png,gif,svg,css,js,woff,woff2,ttf,eot,ico
   ```

6. **Use JSON output for complex analysis**:
   ```bash
   dedup -output=json | jq '.[] | select(.count > 10)'
   ```

## 🚧 Changelog from Original Version

### Bug Fixes
- ✅ Changed `-sort-params` default from `true` to `false`
- ✅ Query parameters now always normalized, not just when sorting
- ✅ Fixed query encoding inconsistencies
- ✅ Improved error handling with proper error messages

### New Features
- 🆕 `-fuzzy` mode for ID normalization
- 🆕 `-ignore-extensions` for file filtering
- 🆕 `-allow-domains` and `-block-domains` for domain filtering
- 🆕 `-keep-scheme` to distinguish http/https
- 🆕 `-path-include-query` for path mode with query strings
- 🆕 `-output=json` and `-output=csv` formats
- 🆕 `-stats` for processing statistics
- 🆕 `-verbose` for debugging
- 🆕 `params` mode for parameter discovery
- 🆕 Comprehensive error messages
- 🆕 Line number tracking in verbose mode

### Improvements
- 📈 Better documentation and examples
- 📈 Comprehensive test suite
- 📈 More consistent URL normalization
- 📈 Better handling of edge cases
- 📈 Improved performance with string builders

## 📝 License

MIT License - Feel free to use in your bug bounty workflows!

## 🤝 Contributing

Contributions welcome! Some ideas for future enhancements:
- [ ] Parallel processing for large inputs
- [ ] Custom regex patterns for fuzzy matching
- [ ] SQLite output for persistent storage
- [ ] Rate limiting for different endpoints
- [ ] Integration with burp suite
- [ ] HTML report generation

## 📧 Support

For bug reports and feature requests, please create an issue with:
- Example input URLs
- Command used
- Expected vs actual output
- Error messages (if any)

Happy hunting! 🎯🐛
