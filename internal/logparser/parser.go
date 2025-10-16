package logparser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

// LogEntry represents a parsed log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Source    string
}

// Parser handles log parsing operations
type Parser struct {
	entries []LogEntry
}

// NewParser creates a new log parser
func NewParser() *Parser {
	return &Parser{
		entries: make([]LogEntry, 0),
	}
}

// ParseFile parses a log file and returns matching entries
func (p *Parser) ParseFile(filePath string, pattern string) ([]LogEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var regex *regexp.Regexp
	if pattern != "" {
		regex, err = regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern: %w", err)
		}
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// If a pattern is provided, only add entries that match
		if regex != nil && !regex.MatchString(line) {
			continue
		}

		// Create a basic log entry
		entry := LogEntry{
			Timestamp: time.Now(), // In a real implementation, parse timestamp from the log line
			Level:     "INFO",     // In a real implementation, extract log level
			Message:   line,
			Source:    fmt.Sprintf("%s:%d", filePath, lineNumber),
		}

		p.entries = append(p.entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return p.entries, nil
}

// GetEntries returns all parsed log entries
func (p *Parser) GetEntries() []LogEntry {
	return p.entries
}

// ClearEntries clears all stored entries
func (p *Parser) ClearEntries() {
	p.entries = make([]LogEntry, 0)
}