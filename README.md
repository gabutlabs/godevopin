# Devopin

A comprehensive logging and system monitoring tool written in Go that provides both REST API, a modern Vue.js web UI, and CLI interfaces for log file analysis and system performance monitoring.

## Project Architecture & Structure

This project follows Go best practices for organizing code, with a clear separation between application code, shared packages, and internal utilities. It combines a robust backend built with Fiber and GORM, a dynamic Vue/Vite frontend (embedded at runtime), and a flexible CLI using Cobra.

```text
devopin/
├── README.md
├── go.mod
├── go.sum
├── Makefile
├── cmd/
│   ├── main.go                # Entry point for the application
│   └── cli/                   # Cobra CLI commands (root, serve, worker)
├── configs/                   # Configuration files (YAML based)
├── frontend/                  # Vue 3 / Vuetify SPA frontend application
├── internal/
│   ├── web/                   # Embedded frontend assets (Vue/Vite dist folder)
│   ├── config/                # Application configuration loader (Viper)
│   ├── database/              # PostgreSQL DB connection & GORM setup
│   ├── handler/               # Fiber HTTP handlers & Socket.io endpoints
│   ├── logparser/             # Log parsing and analysis engines
│   ├── model/                 # GORM database models (User, SystemMetric, etc.)
│   ├── monitoring/            # System & hardware monitoring (CPU, Mem, Network)
│   ├── repository/            # Data access layer
│   ├── services/              # Business logic layer
│   └── worker/                # Background task processing
└── pkg/                       # Shared utility packages
```

## Features

- **Embedded Web Frontend**: A sleek Vue.js 3 / Vuetify single-page admin panel embedded directly into the Go binary. No extra web servers required.
- **Log Tracing**: Parse and analyze `.log` files for specific patterns and issues with support for standard log formats.
- **System Monitoring**: Real-time cross-platform CPU, memory, disk, and network resource monitoring.
- **REST API & WebSockets**: Robust HTTP interface with JWT auth, and WebSocket endpoints (e.g., for real-time Docker logs).
- **CLI Application**: Command-line interface for local operations, worker tasks, and server execution.
- **Alarms & Alerts**: Configurable monitoring thresholds and alerting systems.
- **Docker Services**: Built-in support to monitor and manage Docker services via websockets.
- **Configurable & Extensible**: Modular architecture with flexible YAML configuration.

## Requirements

- Go 1.24.0+
- Node.js (for frontend development)
- PostgreSQL

## Configuration

The application uses a YAML-based configuration. Copy `configs/config.yaml.example` to `configs/config.yaml` and configure your database and JWT settings:

```yaml
database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "password"
  dbname: "devopin"
app:
  jwt_secret: "your-secure-secret"
```

## Building and Running

Devopin includes a comprehensive `Makefile` to handle building both the frontend and backend.

### Development Build

To build the entire application (compiles the Vue frontend and embeds it into the Go backend):

```bash
make build_all
```

To build just the backend or frontend:
```bash
make build_core_on_linux
make build_frontend
```

### Running the Application

Devopin uses a single binary for both Server and Worker modes.

#### API Server Mode

The `serve` command will automatically start the Fiber web server, run database migrations, and serve the embedded Vue frontend.

```bash
# Start server on default port (8080)
./devopin serve

# Start server on a specific port
./devopin serve --port 9000
```

#### Worker Mode (CLI)

```bash
# Run a worker task
./devopin worker --task "example-task"
```

## License

MIT License