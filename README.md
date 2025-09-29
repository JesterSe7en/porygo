# Scrapego – Production-Ready Concurrent Web Scraper

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Scrapego is a high-performance, production-ready CLI tool for concurrent web scraping built in Go. It demonstrates enterprise-grade patterns including worker pools, caching, rate limiting, and extensible plugin architecture.

## Features

### Current Implementation
- **Concurrent Scraping**: Worker pool pattern with configurable concurrency
- **Intelligent Caching**: bbolt-based caching with TTL and cross-platform storage
- **Retry Logic**: Exponential backoff with configurable retry attempts
- **Configuration Management**: TOML files + CLI flags with precedence handling
- **Structured Logging**: Production-ready logging with zap
- **Cross-Platform**: Proper cache directory handling for Windows, macOS, and Linux
- **Input Flexibility**: Support for command-line arguments and stdin piping

### Architecture Highlights
- Clean separation of concerns with internal packages
- Interface-based design for extensibility
- Proper error handling and context propagation
- Resource cleanup and graceful operation

---

## Installation

```bash
# Clone and build
git clone https://github.com/JesterSe7en/scrapego.git
cd scrapego
go build -o scrapego

# Or install directly
go install github.com/JesterSe7en/scrapego@latest
```

---

## Usage

### Basic Usage
```bash
# Scrape multiple URLs concurrently
scrapego https://example.com https://golang.org https://github.com

# Pipe URLs from stdin
cat urls.txt | scrapego

# With custom configuration
scrapego -c 10 -t 30s -r 5 https://example.com
```

### Configuration Options
```bash
# Performance tuning
scrapego -c 20 --timeout 15s --retry 3 --backoff 5s <urls>

# Logging and debugging
scrapego --verbose --debug --log ./scrape.log <urls>

# Cache management
scrapego cache clear
scrapego cache list

# Configuration file
scrapego --config ./custom-config.toml <urls>
```

### Configuration File (config.toml)
```toml
concurrency = 10
timeout = "30s"
retry = 3
backoff = "2s"
output = "JSON"

[database]
expiration = "1h"
```

---

## Architecture

### Current Structure
```
scrapego/
├── cmd/                    # CLI commands (Cobra)
│   ├── cache/             # Cache management subcommands
│   └── config/            # Configuration subcommands
├── config/                # Configuration management
├── internal/              # Private packages
│   ├── database/          # bbolt caching layer
│   ├── flags/             # CLI flag definitions
│   ├── logger/            # Structured logging
│   ├── scraper/           # HTTP scraping logic
│   └── workerpool/        # Concurrent worker implementation
└── main.go               # Application entry point
```

### Design Patterns Implemented
- **Worker Pool Pattern**: Efficient concurrent processing
- **Configuration Precedence**: CLI flags > config file > defaults
- **Cross-Platform Compatibility**: Proper cache directory handling
- **Resource Management**: Proper cleanup and connection handling

---

## Current Capabilities

### Scraping Engine
- HTTP/HTTPS support with custom User-Agent
- Configurable timeouts and retry logic
- Status code validation and error handling
- Context-based cancellation support

### Caching System
- Embedded bbolt key-value storage
- Automatic cache expiration
- Cross-platform cache directory resolution
- Cache management commands

### Configuration Management
- TOML-based configuration files
- CLI flag override capability
- Validation and default value handling
- Environment-specific configurations

---

## Planned Enhancements

### Immediate Priority (Production Readiness)
- [ ] **Database Abstraction Layer**: Interface-based storage for better testability
- [ ] **Structured Output Formats**: JSON, CSV, XML with schema validation
- [ ] **Comprehensive Error Handling**: Detailed error types and recovery strategies
- [ ] **Enhanced Testing Suite**: Unit tests, integration tests, and benchmarks

### Advanced Features (Enterprise Grade)
- [ ] **Plugin Architecture**: Extensible data extraction system
  ```go
  type Extractor interface {
      Name() string
      Extract(*http.Response, string) (interface{}, error)
  }
  ```

- [ ] **Rate Limiting**: Per-domain token bucket implementation
  - Configurable requests per second
  - Burst handling
  - Domain-specific limits

- [ ] **Circuit Breaker Pattern**: Fault tolerance for unreliable endpoints
  - Automatic failure detection
  - Recovery mechanisms
  - Configurable thresholds

- [ ] **Observability Stack**: Production monitoring capabilities
  - Prometheus metrics export
  - Distributed tracing support
  - Health check endpoints
  - Performance profiling hooks

- [ ] **Advanced Scraping Features**:
  - HTML parsing with goquery/colly
  - JavaScript rendering with headless browsers
  - Sitemap discovery and parsing
  - robots.txt compliance
  - Content deduplication

### Operational Excellence
- [ ] **Graceful Shutdown**: Signal handling and resource cleanup
- [ ] **Hot Configuration Reload**: Runtime configuration updates
- [ ] **Metrics Dashboard**: Real-time scraping statistics
- [ ] **Resume Capability**: Pause and resume large scraping jobs
- [ ] **Progress Indicators**: Real-time progress reporting

---

## Development

### Running Tests
```bash
go test ./...
go test -race ./...
go test -bench=. ./...
```

### Building
```bash
# Development build
go build -o scrapego

# Production build with optimizations
go build -ldflags="-s -w" -o scrapego
```

---

## Performance Characteristics

- **Concurrency**: Configurable worker pools (default: 5 workers)
- **Memory**: Efficient channel-based communication
- **Storage**: Embedded database with minimal overhead
- **Network**: Configurable timeouts and retry logic
- **Caching**: TTL-based invalidation with automatic cleanup

---

## Contributing

This project showcases modern Go development practices and is designed to demonstrate:

- **Clean Architecture**: Proper separation of concerns
- **Interface Design**: Extensible and testable code
- **Concurrency Patterns**: Safe and efficient parallel processing
- **Configuration Management**: Flexible and user-friendly setup
- **Error Handling**: Comprehensive error management
- **Testing Strategy**: Unit, integration, and performance testing

---

## Why This Project?

### For Resume/Portfolio
- Demonstrates **concurrent programming** expertise
- Showcases **system design** thinking with clean architecture
- Exhibits **production readiness** with proper logging, caching, and configuration
- Shows **extensibility planning** with plugin architecture roadmap

### Technical Skills Highlighted
- **Go Language Mastery**: Advanced patterns and idioms
- **CLI Development**: User-friendly command-line interfaces
- **Database Integration**: Embedded storage solutions
- **Network Programming**: HTTP clients and error handling
- **Testing**: Comprehensive test coverage and benchmarking
- **DevOps Awareness**: Configuration management and observability

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Example Output

```bash
$ scrapego https://example.com https://golang.org
Scraping 2 URLs with 5 workers...
https://example.com -> text/html; charset=UTF-8
https://golang.org -> text/html; charset=utf-8
Cache entries: 2
Total processing time: 1.2s
```

*Future versions will support structured JSON/CSV output with extracted data, metadata, and comprehensive error reporting.*
