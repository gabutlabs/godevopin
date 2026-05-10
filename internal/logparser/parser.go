package logparser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
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

// ParseFile parses a log file starting from a specific byte offset and returns matching entries along with the new file offset
func (p *Parser) ParseFile(filePath string, pattern string, startOffset int64) ([]LogEntry, int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Check file size to handle rotation/truncation
	stat, err := file.Stat()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to stat file: %w", err)
	}
	
	if startOffset > stat.Size() {
		// File was likely truncated/rotated
		startOffset = 0
	}

	if startOffset > 0 {
		if _, err := file.Seek(startOffset, 0); err != nil {
			return nil, 0, fmt.Errorf("failed to seek file: %w", err)
		}
	}

	var regex *regexp.Regexp
	if pattern != "" {
		// Convert {timestamp} to (?P<timestamp>.*?) etc if it's in the builder format
		if strings.Contains(pattern, "{") && strings.Contains(pattern, "}") {
			// First, escape all regex special characters
			pattern = regexp.QuoteMeta(pattern)
			
			// Then replace the escaped placeholders with named capture groups
			// Note: regexp.QuoteMeta escapes { and } as \{ and \}
			pattern = strings.ReplaceAll(pattern, `\{timestamp\}`, "(?P<timestamp>.*?)")
			pattern = strings.ReplaceAll(pattern, `\{level\}`, "(?P<level>.*?)")
			pattern = strings.ReplaceAll(pattern, `\{message\}`, "(?P<message>.*?)")
			pattern = strings.ReplaceAll(pattern, `\{source\}`, "(?P<source>.*?)")
			
			// Add anchors if needed or keep as is
			if !strings.HasPrefix(pattern, "^") {
				pattern = "^" + pattern
			}
			if !strings.HasSuffix(pattern, "$") {
				pattern = pattern + "$"
			}
		}

		regex, err = regexp.Compile(pattern)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid pattern %s: %w", pattern, err)
		}
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// If a pattern is provided, extract named groups
		if regex != nil {
			match := regex.FindStringSubmatch(line)
			if match == nil {
				continue
			}

			entry := LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Message:   line,
				Source:    fmt.Sprintf("%s:%d", filePath, lineNumber),
			}

			result := make(map[string]string)
			for i, name := range regex.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}

			if val, ok := result["timestamp"]; ok {
				if t, err := parseTimestamp(val); err == nil {
					entry.Timestamp = t
				}
			}
			if val, ok := result["level"]; ok {
				entry.Level = val
			}
			if val, ok := result["message"]; ok {
				entry.Message = val
			}

			p.entries = append(p.entries, entry)
		} else {
			// Basic entry if no regex
			entry := LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Message:   line,
				Source:    fmt.Sprintf("%s:%d", filePath, lineNumber),
			}
			p.entries = append(p.entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, fmt.Errorf("error reading file: %w", err)
	}

	// Assuming we read to the end of the file
	return p.entries, stat.Size(), nil
}

// parseTimestamp tries to parse a timestamp string into time.Time
func parseTimestamp(val string) (time.Time, error) {
	// Add common log timestamp formats
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"02/Jan/2006:15:04:05 -0700", // Apache/Nginx
		"2006/01/02 15:04:05",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"Jan _2 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, val); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("could not parse timestamp: %s", val)
}

// GetEntries returns all parsed log entries
func (p *Parser) GetEntries() []LogEntry {
	return p.entries
}

// ClearEntries clears all stored entries
func (p *Parser) ClearEntries() {
	p.entries = make([]LogEntry, 0)
}