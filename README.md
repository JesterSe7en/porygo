# Porygo

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/JesterSe7en/porygo)](https://goreportcard.com/report/github.com/JesterSe7en/porygo)

A high-performance, concurrent web scraper built in Go with intelligent caching, retry mechanisms, and flexible data extraction capabilities.

> **Disclaimer:** This project is currently under active development. Features and command-line flags are subject to change.

## What's in a name?

The name "Porygo" is a tribute to a love for Pokémon and programming. It's a portmanteau of **Porygon**, the virtual Pokémon that exists and travels through cyberspace, and **Go**, the language this project is written in.

## Features

- **Concurrent Processing**: Employs a worker pool to manage and execute multiple scraping jobs simultaneously.
- **Intelligent Caching**: Utilizes a BBolt database to cache responses, minimizing redundant network requests.
- **Smart Retry Logic**: Implements exponential backoff with optional jitter to gracefully handle transient network errors.
- **Flexible Data Extraction**: Supports data extraction using CSS selectors (via `goquery`) and regex patterns.
- **Multiple Output Formats**: Presents scraped data in either JSON or plain text formats.
- **Layered Configuration**: Settings can be specified via a `config.toml` file and overridden with command-line flags.
- **Structured Logging**: Provides detailed operational insights using the `zap` logging library.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Prerequisites

- Go 1.25 or higher.

### Build from Source

Clone the repository and build the binary:

```bash
git clone https://github.com/JesterSe7en/porygo.git
cd porygo
go build -o porygo .
```

### Install from Source

Install directly using `go install`:

```bash
go install github.com/JesterSe7en/porygo@latest
```

## Usage

### Basic Scraping

Scrape URLs provided as command-line arguments. The tool can also accept URLs piped from `stdin`.

```bash
# Scrape a single URL
./porygo https://example.com

# Scrape multiple URLs
./porygo https://example.com https://golang.org

# Scrape URLs from a file
cat list.txt | ./porygo
```

### Data Extraction

Use CSS selectors (`-s`) or regex patterns (`-p`) to extract specific content.

```bash
# Extract all h1 and h2 tags
./porygo -s "h1" -s "h2" https://example.com

# Extract email addresses using regex
./porygo -p "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}" https://example.com

# Output in plain text
./porygo -o plain https://example.com
```

### Command-Line Flags

```
Usage:
  porygo [urls...] [flags]
  porygo [command]

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
  -h, --help                   help for porygo
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

The `cache` command helps manage the local data store.

```bash
# Clear all cached results
./porygo cache clear
```

### Configuration Management

The `config` command assists with the configuration file.

```bash
# Create a default 'config.toml' file in the current directory
./porygo config init
```

## Configuration

`porygo` can be configured using a `config.toml` file.

### Configuration Precedence

1.  **Command-line flags** (highest priority)
2.  **Values in `config.toml`**
3.  **Default values** (lowest priority)

### Example `config.toml`

Run `porygo config init` to generate a file with default values.

```toml
# Default number of concurrent workers
concurrency = 5
# Default timeout for each HTTP request
timeout = "10s"
# Default output format ("json" or "plain")
format = "json"
# Default number of retries for failed requests
retry = 3
# Force scraping and ignore existing cache
force = false
# Suppress logs and only show scraped data
quiet = false
# Include response headers in the output
headers = false

[backoff]
  # Base delay for the first retry
  base_delay = "1s"
  # Enable or disable random jitter in retry delays
  jitter = true

[selectors]
  # Default CSS selectors to apply
  select = []
  # Default regex patterns to apply
  pattern = []

[database]
  # Duration for which cached items remain valid
  expiration = "24h"
```

## Architecture

The project follows a modular architecture to separate concerns, making it easier to maintain and extend.

```
├── cmd/                    # Command-line interface (Cobra)
│   ├── cache/              # Cache management commands
│   ├── config/             # Configuration commands
│   └── root.go             # Root command and CLI setup
├── config/                 # Configuration management (TOML)
├── internal/               # Internal application logic
│   ├── app/                # Core application wiring
│   ├── flags/              # CLI flag definitions
│   ├── logger/             # Structured logging (Zap)
│   ├── presenter/          # Output formatting
│   ├── scraper/            # Web scraping logic (Goquery)
│   ├── storage/            # Caching and persistence (BBolt)
│   └── workerpool/         # Concurrent worker management
└── main.go                 # Application entry point
```

### Dependencies

- **[Cobra](https://github.com/spf13/cobra)**: CLI framework
- **[Zap](https://github.com/uber-go/zap)**: Structured logging
- **[BBolt](https://github.com/etcd-io/bbolt)**: Embedded key-value database for caching
- **[GoQuery](https://github.com/PuerkitoBio/goquery)**: HTML parsing and CSS selection
- **[TOML](https://github.com/BurntSushi/toml)**: Configuration file parsing

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Author

**Alexander Chan** - [JesterSe7en](https://github.com/JesterSe7en)
