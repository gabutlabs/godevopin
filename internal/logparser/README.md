# Log Parser

This internal package provides core functionality for parsing and analyzing log files. It is designed specifically for the Devopin application and should not be imported by external projects.

## Purpose

The log parser package handles:
- Reading and parsing various log file formats
- Pattern matching and filtering of log entries
- Extraction of relevant information from log lines
- Performance-optimized searching through large log files
- Common log analysis operations

## Structure

```
logparser/
├── README.md (This file)
├── parser.go
├── filters.go
├── models.go
├── regex.go
└── utils.go
```

## Features

- **Format Agnostic**: Supports common log formats (Apache, Nginx, JSON, Syslog, etc.)
- **Pattern Matching**: Advanced regex-based pattern matching capabilities
- **Performance Optimized**: Efficient algorithms for processing large log files
- **Filtering**: Rich filtering options for analyzing specific log entries
- **Timestamp Handling**: Automatic timestamp parsing and comparison
- **Context Extraction**: Ability to extract surrounding log entries for context

## Key Components

- **Parser**: Core parser that handles reading and initial analysis of log files
- **Filters**: Various filter implementations for different search criteria
- **Models**: Data structures representing log entries and analysis results
- **Regex**: Pre-compiled regular expressions for common log formats

## Usage

The log parser is primarily used by the CLI and API applications but can be used directly within the internal packages:

```go
// Example usage within internal packages
parser := logparser.NewParser()
results, err := parser.ParseFile("/path/to/logfile.log", filters...)
if err != nil {
    // handle error
}
for _, entry := range results.Entries {
    // Process log entry
}
```

## Supported Log Formats

- **Standard text logs**: Plain text log files with various formats
- **JSON logs**: Structured logs in JSON format
- **Syslog**: Standard syslog format
- **Apache/Nginx**: Common web server log formats
- **Custom formats**: Configurable format definitions

## Performance Considerations

- The parser is designed to handle large files efficiently using buffered reading
- For very large files, the parser implements streaming to avoid memory exhaustion
- Pattern matching is optimized using pre-compiled regular expressions
- Concurrency is used where appropriate for I/O operations