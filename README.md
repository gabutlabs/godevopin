# Devopin

A comprehensive logging and system monitoring tool written in Go that provides both REST API and CLI interfaces for log file analysis and system performance monitoring.

## Project Structure

This project follows Go best practices for organizing code, with a clear separation between application code, shared packages, and internal utilities.

```
devopin/
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── cli/
│       └── main.go
├── internal/
│   ├── logparser/
│   │   ├── parser.go
│   │   ├── filters.go
│   │   └── models.go
│   └── monitoring/
│       ├── cpu.go
│       ├── memory.go
│       └── system.go
├── pkg/
│   ├── config/
│   ├── utils/
│   └── types/
├── examples/
│   └── example_usage.go
└── docs/
    └── api.md
```

### Directory Descriptions

- **cmd/**: Contains main applications. Single binary with multiple modes.
  - **devopin/**: Main application with both server and worker modes
- **internal/**: Internal packages that should not be imported by external projects
  - **logparser/**: Log parsing and analysis functionality  
  - **monitoring/**: System monitoring and metrics collection utilities
- **pkg/**: Shared libraries that can be used by external projects
- **examples/**: Example usages and code samples
- **docs/**: Documentation files

## Features

- **Log Tracing**: Parse and analyze .log files for specific patterns and issues
- **System Monitoring**: Real-time CPU, memory, and system resource monitoring
- **REST API**: HTTP interface for remote monitoring and log analysis
- **CLI Application**: Command-line interface for local operations
- **Configurable**: Flexible configuration system with multiple input methods
- **Extensible**: Modular architecture allowing easy feature additions

## Installation

```bash
go install github.com/username/devopin@latest
```

## Usage

Devopin now supports both API and worker modes in a single binary:

### API Server Mode

```bash
# Start API server
devopin server --port 8080

# Query endpoints
curl http://localhost:8080/api/logs/analyze
curl http://localhost:8080/api/monitoring/system
```

### Worker Mode (CLI)

```bash
# Trace log file
devopin worker trace -f /path/to/app.log

# Monitor system resources
devopin worker monitor --cpu --memory

# More examples with worker commands
```

## Development

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Build: `go build ./cmd/cli` or `go build ./cmd/api`
4. Run tests: `go test ./...`

## License

MIT License