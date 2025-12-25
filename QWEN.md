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
├── Makefile
├── QWEN.md
├── .gitignore
├── build/...
├── cmd/
│   └── devopin/
│       ├── main.go
│       ├── cmd/
│       │   ├── root.go
│       │   ├── serve.go
│       │   └── worker.go
│       └── web/
├── configs/
│   └── config.yaml.example
├── frontend/...
├── internal/
│   ├── config/
│   ├── database/
│   ├── handler/
│   ├── logparser/
│   ├── model/
│   ├── monitoring/
│   ├── repository/
│   ├── services/
│   └── worker/
└── pkg/
    ├── pagination.go
    ├── query-builder.go
    ├── util.go
    └── validator.go
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

#### Building with Makefile

The project includes a Makefile for building the complete application including the frontend:

```bash
# Build the entire application (frontend + backend)
make build_all

# Clean build artifacts and frontend assets
make clean

# Build only the frontend
make build_frontend

# Build only the core backend
make build_core_on_linux
```

### Configuration

The application uses a YAML-based configuration system. The primary configuration file is located at `configs/config.yaml` (based on the example `configs/config.yaml.example`) and contains database connection settings and JWT secret:

```yaml
database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: ""
  dbname: "devopin"
app:
  jwt_secret: "aaaaa"
settings:
  monitoring_interval_seconds: 10
  alarms:
    check_interval_seconds: 60
    repeat_interval_minutes: 15
    thresholds:
      system_cpu_critical_percent: 90.0
      system_disk_critical_percent: 85.0
      system_mem_critical_percent: 85.0
      worker_heartbeat_timeout_seconds: 300
```

### API Endpoints

The application provides REST API endpoints for:
- Authentication (`/auth/*`)
- User management (`/users/*`)
- System metrics (`/system-metrics/*`)
- Worker services (`/worker-services/*`)
- Alarms (`/alarms/*`)
- Docker services (`/docker/*`)

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