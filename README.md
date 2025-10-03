# Scrapego

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/JesterSe7en/scrapego)](https://goreportcard.com/report/github.com/JesterSe7en/scrapego)

A high-performance, concurrent web scraper built in Go with intelligent caching, retry mechanisms, and flexible data extraction capabilities.

## Features

- **Concurrent Processing**: Worker pool architecture with configurable concurrency levels
- **Smart Retry Logic**: Exponential backoff with jitter for handling transient failures
- **Intelligent Caching**: BBolt-based caching system to avoid redundant requests
- **Flexible Data Extraction**: CSS selectors and regex pattern matching
- **Multiple Output Formats**: JSON and plain text output options
- **Configurable**: TOML configuration files with CLI flag overrides
- **Structured Logging**: Comprehensive logging with Zap for debugging and monitoring
- **Extensible Architecture**: Clean, modular design with interface-based components

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [Examples](#examples)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Pre-built Binaries

Download the latest release from the [releases page](https://github.com/JesterSe7en/scrapego/releases).

### Install from Source

Requires Go 1.25 or higher:

```bash
go install github.com/JesterSe7en/scrapego@latest
```

### Build from Source

```bash
git clone https://github.com/JesterSe7en/scrapego.git
cd scrapego
go build -o scrapego .
```

## Quick Start

### Basic Usage

Scrape URLs directly from command line:

```bash
# Scrape single URL
scrapego https://example.com

# Scrape multiple URLs
scrapego https://example.com https://golang.org https://github.com
```

### Using Stdin

```bash
# From a file
cat urls.txt | scrapego

# From command output
echo "https://example.com" | scrapego
```

### Extract Specific Data

```bash
# Extract titles using CSS selectors
scrapego -s "title" -s "h1" https://example.com

# Match patterns with regex
scrapego -p "email.*@.*\.com" https://example.com

# Output in plain text format
scrapego -o plain https://example.com
```

## Usage

### Command Line Options

```
Usage:
  scrapego [urls...] [flags]
  scrapego [command]

Available Commands:
  cache       Manage cached scraping results
  config      View and modify CLI configuration
  help        Help about any command

Flags:
  -c, --concurrency int        number of workers (default 5)
      --config string          specify config file
  -d, --debug                  output debug messages
  -f, --force                  ignore cache and scrape fresh data
  -o, --format string          output format (json|plain) (default "json")
  -H, --headers                include response headers
  -h, --help                   help for scrapego
  -l, --log string             file path to write logs
  -p, --pattern strings        regex patterns to match
  -q, --quiet                  only output extracted data
  -r, --retry int              number of retries per URL on failure (default 3)
      --retry-delay duration   base delay between retries (default 1s)
      --retry-jitter           enable jitter for retry delays (default true)
  -s, --select strings         CSS selectors to extract
  -t, --timeout duration       request timeout per URL (default 10s)
  -v, --verbose                show logs for each step
```

### Cache Management

```bash
# Clear all cached results
scrapego cache clear
```

### Configuration Management

```bash
# Initialize default config file
scrapego config init

# View current configuration
scrapego config show
```

## Configuration

### Configuration File

Create a `config.toml` file in your working directory:

```toml
# Worker configuration
concurrency = 10
timeout = "30s"
retry = 5
force = false

# Output configuration
output = "json"

# Retry configuration
[backoff]
  base_delay = "2s"
  jitter = true

# Cache configuration
[database]
  expiration = "24h"
```

### Configuration Precedence

1. **Command-line flags** (highest priority)
2. **Configuration file**
3. **Default values** (lowest priority)

### Environment Variables

You can also use environment variables (useful for CI/CD):

```bash
export SCRAPEGO_CONCURRENCY=10
export SCRAPEGO_TIMEOUT=30s
```

## Examples

### Basic Web Scraping

```bash
# Scrape with custom settings
scrapego -c 10 -t 30s -r 5 https://example.com

# Extract specific elements
scrapego -s "title" -s ".main-content" -s "meta[name='description']" https://example.com

# Match email addresses and phone numbers
scrapego -p "\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b" \
         -p "\b\d{3}-\d{3}-\d{4}\b" \
         https://example.com
```

### Batch Processing

```bash
# Create a URL list file
echo "https://example.com
https://golang.org
https://github.com" > urls.txt

# Process with custom configuration
scrapego -f urls.txt -c 20 -o plain --quiet
```

### Advanced Usage

```bash
# Scrape with headers included and verbose logging
scrapego -H -v -l scrape.log https://api.example.com

# Force fresh scraping (ignore cache) with custom retry settings
scrapego -f --retry-delay 5s --retry-jitter=false https://example.com

# Extract and format specific data
scrapego -s "h1,h2,h3" -p "https?://[^\s]+" -o json https://news.example.com
```

### Integration Examples

```bash
# Use with jq for JSON processing
scrapego https://api.example.com | jq '.data[] | select(.status == "active")'

# Combine with other tools
curl -s https://example.com/sitemap.xml | \
  grep -oP 'https://[^<]+' | \
  scrapego -c 50 -q
```

## Architecture

Scrapego is built with a clean, modular architecture following Go best practices:

```
├── cmd/                    # Command-line interface
│   ├── cache/             # Cache management commands
│   ├── config/            # Configuration commands
│   └── root.go            # Root command and CLI setup
├── config/                # Configuration management
├── internal/              # Internal packages
│   ├── app/               # Main application logic
│   ├── flags/             # CLI flag definitions
│   ├── logger/            # Structured logging
│   ├── presenter/         # Output formatting
│   ├── scraper/           # Web scraping logic
│   ├── storage/           # Caching and persistence
│   └── workerpool/        # Concurrent worker management
└── main.go                # Entry point
```

### Key Design Patterns

- **Worker Pool Pattern**: Manages concurrent scraping operations efficiently
- **Interface-Based Design**: Enables easy testing and component swapping
- **Configuration Precedence**: Clear hierarchy for configuration sources
- **Separation of Concerns**: Each package has a single, well-defined responsibility

### Dependencies

- **[Cobra](https://github.com/spf13/cobra)**: Powerful CLI framework
- **[Zap](https://github.com/uber-go/zap)**: High-performance structured logging
- **[BBolt](https://github.com/etcd-io/bbolt)**: Embedded key-value database for caching
- **[GoQuery](https://github.com/PuerkitoBio/goquery)**: jQuery-like HTML parsing
- **[TOML](https://github.com/BurntSushi/toml)**: Configuration file parsing

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/scraper
```

## Performance

- **Concurrent Processing**: Handles hundreds of URLs simultaneously
- **Memory Efficient**: Streaming processing with controlled memory usage
- **Network Optimized**: Connection pooling and intelligent retry mechanisms
- **Cache Performance**: Fast BBolt-based caching reduces redundant requests

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
git clone https://github.com/JesterSe7en/scrapego.git
cd scrapego
go mod download
go run . --help
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Alexander Chan** - [JesterSe7en](https://github.com/JesterSe7en)

---

<div align="center">
  <strong>⭐ If you found this project helpful, please consider giving it a star! ⭐</strong>
</div>
