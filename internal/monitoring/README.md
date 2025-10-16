# System Monitoring

This internal package provides core functionality for monitoring system resources such as CPU usage, memory consumption, disk I/O, and network activity. It is designed specifically for the Devopin application and should not be imported by external projects.

## Purpose

The monitoring package handles:
- Real-time collection of system resource metrics
- Performance data aggregation and analysis
- Threshold monitoring and alerting capabilities
- Cross-platform system information gathering
- Historical data collection for trend analysis

## Structure

```
monitoring/
├── README.md (This file)
├── cpu.go
├── memory.go
├── disk.go
├── network.go
├── system.go
├── metrics.go
└── collector.go
```

## Features

- **CPU Monitoring**: Detailed CPU usage statistics including per-core metrics
- **Memory Monitoring**: Physical and virtual memory usage tracking
- **Disk I/O**: Disk usage and I/O performance monitoring
- **Network Monitoring**: Network interface statistics and bandwidth tracking
- **Cross-Platform**: Compatible with Linux, macOS, and Windows systems
- **Low Overhead**: Efficient monitoring with minimal system impact
- **Configurable Intervals**: Adjustable monitoring intervals based on requirements
- **Threshold Alerts**: Configurable thresholds for triggering alerts

## Key Components

- **CPU**: CPU usage statistics, load averages, and frequency information
- **Memory**: Physical and virtual memory statistics
- **Disk**: Storage usage and I/O statistics
- **Network**: Network interface statistics
- **System**: Overall system information and health metrics
- **Metrics**: Standardized data structures for metrics
- **Collector**: Aggregates metrics from different sources

## Usage

The monitoring package is primarily used by both the CLI and API applications:

```go
// Example usage within internal packages
collector := monitoring.NewCollector()
stats, err := collector.GetSystemStats()
if err != nil {
    // handle error
}

fmt.Printf("CPU Usage: %.2f%%\n", stats.CPU.UsagePercent)
fmt.Printf("Memory Usage: %.2f%%\n", stats.Memory.UsagePercent)
```

## Supported Platforms

- **Linux**: Full support using /proc and /sys filesystems
- **macOS**: Support through native APIs and system calls
- **Windows**: Support via Windows Performance Counters and WMI

## Performance Considerations

- Monitoring operations are designed to be lightweight
- Configurable sampling intervals to balance accuracy with performance
- Concurrent metric collection where possible
- Cached results for frequently accessed static information