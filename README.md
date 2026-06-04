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

## Core Features Overview

Devopin is packed with built-in tools for holistic infrastructure management. Here is a high-level overview of the main capabilities:

### 1. 📊 System Monitoring Dashboard
Real-time and historical visualization of system performance. It tracks key metrics such as **CPU, Memory, and Disk Usage**. Data is persistently stored using **PostgreSQL with TimescaleDB** for efficient time-series querying, allowing you to filter charts from the last 1 hour up to the last 30 days with seamless auto-refresh capabilities.

### 2. 🚨 Alarms & Alerts Management
Monitors system metrics against predefined thresholds. Whenever a threshold is breached, Devopin triggers an alert. You can manage **Active** alarms (FIRING) and review the **History** of past or ACKNOWLEDGED alarms to keep your infrastructure healthy.

### 3. 🐳 Docker Management
Get full visibility into your Docker environment directly from the UI. You can monitor and manage **Containers, Images, Networks, and Volumes**. Integrated with the Docker Engine API, it provides real-time status updates via WebSockets.

### 4. ⚙️ Worker Services Control Panel
A built-in control panel to orchestrate background tasks and internal daemon processes. You can monitor the lifecycle of each worker (**Start, Stop, Restart**), track their desired state, PIDs, health status, and last heartbeats.

### 5. 📝 Log Tracing Engine
Designed to assist developers and sysadmins in analyzing `.log` files. Devopin includes a parsing engine that can extract meaningful information (Error Levels, Messages, Timestamps) from standard log formats, making it easier to trace specific issues across large log files.

### 6. 🔐 User & Access Management
A complete CRUD interface for managing user access to the Devopin dashboard. It features secure **JWT-based authentication**, role-based access control, and password hashing (bcrypt) to ensure your system monitoring is safely guarded.

### 7. 🚀 Embedded Web Frontend & CLI
- **Single Binary Deployment:** The entire Vue.js 3 / Vuetify frontend is embedded directly into the Go binary. No separate web server (like Nginx) is required to serve the UI!
- **Cobra CLI:** Powerful command-line interface for local operations, worker task executions, and launching the server.

## Requirements

- Go 1.24.0+
- Node.js (for frontend development)
- PostgreSQL with TimescaleDB extension enabled
  > After installing PostgreSQL, activate the TimescaleDB extension on your database:
  > ```sql
  > CREATE EXTENSION IF NOT EXISTS timescaledb;
  > ```

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

## Installation

You can install Devopin either by downloading the pre-built binaries from GitHub Releases or by building it from source.

### 1. Auto-Installation Script (Recommended for Linux/macOS)

The easiest way to install Devopin globally on your system is by using the automatic installer script. This will download the latest binary, install it to `/usr/local/bin`, and set up the global configuration directory at `/opt/devopin/`.

```bash
curl -sSL https://raw.githubusercontent.com/gabutlabs/godevopin/main/install.sh | bash
```

After installation, you can simply run:
```bash
devopin serve
```

### 2. Download from GitHub Releases (Manual)

Pre-built binaries are available for Linux, macOS, and Windows.

**Option A: Download via Web Browser**
1. Go to the [Releases page](https://github.com/gabutlabs/godevopin/releases) of this repository.
2. Download the appropriate binary for your operating system and architecture.

**Option B: Download via CLI (curl)**
You can also download the binary directly from the terminal. For example, to download the `linux-arm64` binary for version `v0.2.0-beta`:
```bash
curl -L -o devopin https://github.com/gabutlabs/godevopin/releases/download/v0.2.0-beta/godevopin-linux-arm64
```
*(Make sure to adjust the version and architecture to match your needs)*

**After downloading manually (both options):**
1. Make the binary executable (Linux/macOS): `chmod +x devopin`
2. Run the application: `./devopin serve`

### 2. Build from Source

Devopin includes a comprehensive `Makefile` to handle building both the frontend and backend.

To build the entire application (compiles the Vue frontend and embeds it into the Go backend):

```bash
make build_all
```

To build just the backend or frontend:
```bash
make build_core
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