# Devopin Project Context

## Project Overview

Devopin is a comprehensive logging and system monitoring tool written in Go that provides both REST API and CLI interfaces for log file analysis and system performance monitoring. The project is built using Go 1.24.0 and follows Go best practices for organizing code with a clear separation between application code, shared packages, and internal utilities.

### Architecture & Components

The project follows a modular architecture with the following key components:

- **CLI & API**: Single binary with multiple modes (server and worker)
- **Log Parser**: Handles parsing and analysis of log files with support for various formats (Apache, Nginx, JSON, Syslog, etc.)
- **System Monitoring**: Collects system resource metrics (CPU, memory, disk, network) with cross-platform support
- **Web Framework**: Uses Fiber for the HTTP API
- **Database**: Uses GORM with PostgreSQL as the primary database
- **Configuration**: Uses Viper for configuration management
- **Command Processing**: Uses Cobra for CLI command handling

### Directory Structure

```
devopin/
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── devopin/
│       ├── main.go
│       └── cmd/
│           ├── root.go
│           ├── serve.go
│           └── worker.go
├── configs/
│   └── config.yaml
├── internal/
│   ├── config/
│   ├── database/
│   ├── handler/
│   ├── logparser/
│   ├── model/
│   ├── monitoring/
│   ├── repository/
│   └── services/
└── pkg/
    └── util.go
```

### Key Technologies & Dependencies

- **Go 1.24.0**: Primary language and runtime
- **Cobra**: Command-line interface framework
- **Fiber**: Web framework for HTTP API
- **Viper**: Configuration management
- **GORM**: Database ORM with PostgreSQL driver
- **JWT**: Authentication and authorization
- **bcrypt**: Password hashing
- **gopsutil**: Cross-platform system and hardware monitoring
- **shirou/gopsutil/v4**: System statistics library

### Building and Running

#### Development Setup

1. Ensure Go 1.24.0+ is installed
2. Clone the repository
3. Install dependencies with `go mod tidy`
4. Build the application with `go build ./cmd/devopin`

#### Running the Application

The application supports two main modes:

**Server Mode (API)**:
```bash
# Build the application
go build -o devopin ./cmd/devopin

# Start API server on default port (8080)
./devopin serve

# Start API server on custom port
./devopin serve --port 9000
```

**Worker Mode (CLI)**:
```bash
# Run a worker task
./devopin worker --task "example-task"
```

#### Testing

Run all tests in the project:
```bash
go test ./...
```

### Configuration

The application uses a YAML-based configuration system. The primary configuration file is located at `configs/config.yaml` and contains database connection settings and JWT secret:

```yaml
database:
  host: "localhost"
  port: "5432"
  user: "madina"
  password: ""
  dbname: "devopin"
app:
  jwt_secret: "1i2u12hnwjkbwuygfbe9u1n21jkn918h31ni"
settings:
  monitoring_interval_seconds: 10
```

### API Endpoints

The application provides REST API endpoints for:
- Authentication (`/auth/*`)
- User management (`/users/*`)
- System metrics (`/system-metrics/*`)

### Development Conventions

1. **Code Structure**: Follows Go standard project layout with internal packages
2. **CLI Commands**: Uses Cobra for a structured command-line interface
3. **HTTP API**: Uses Fiber for web framework with proper HTTP routing
4. **Database**: GORM ORM with PostgreSQL for persistence
5. **Authentication**: JWT-based authentication system
6. **Logging**: Structured logging throughout the application
7. **Testing**: Unit tests for all important functions and services

### Key Features

- **Log Tracing**: Parse and analyze .log files for specific patterns and issues with support for standard log formats
- **System Monitoring**: Real-time CPU, memory, disk, and network resource monitoring across platforms
- **REST API**: HTTP interface for remote monitoring and log analysis
- **CLI Application**: Command-line interface for local operations
- **Configurable**: Flexible configuration system with multiple input methods
- **Extensible**: Modular architecture allowing easy feature additions

### Important Files and Functions

- `cmd/devopin/main.go`: Application entry point
- `cmd/devopin/cmd/serve.go`: Server command implementation
- `cmd/devopin/cmd/worker.go`: Worker command implementation
- `internal/monitoring/system.go`: System metrics collection
- `internal/monitoring/README.md`: Documentation for monitoring package
- `internal/logparser/parser.go`: Log file parsing functionality
- `internal/logparser/README.md`: Documentation for log parser package
- `internal/handler/initiatior.go`: API route setup
- `pkg/util.go`: Shared utility functions
- `internal/model/system-metric.go`: System metrics data structure
- `internal/model/user.go`: User data structure