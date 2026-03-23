package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

// LogEntry represents a parsed log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp,omitempty"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Level rank map for filtering
var levelRank = map[string]int{
	"debug": 0,
	"info":  1,
	"warn":  2,
	"error": 3,
	"fatal": 4,
}

func main() {
	filePath := flag.String("file", "", "path to log file (default: stdin)")
	level := flag.String("level", "info", "minimum log level to display")
	format := flag.String("format", "json", "output format: json or logfmt")
	flag.Parse()

	// Validate level
	minLevel, ok := levelRank[*level]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: invalid level '%s'. Valid levels: debug, info, warn, error, fatal\n", *level)
		os.Exit(1)
	}

	// Validate format
	if *format != "json" && *format != "logfmt" {
		fmt.Fprintf(os.Stderr, "Error: invalid format '%s'. Valid formats: json, logfmt\n", *format)
		os.Exit(1)
	}

	// Open file or use stdin
	var scanner *bufio.Scanner
	if *filePath != "" {
		file, err := os.Open(*filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	// Process each line
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		entry := parseLine(line)

		// Filter by level
		entryLevel := strings.ToLower(entry.Level)
		entryRank, exists := levelRank[entryLevel]
		if !exists {
			entryRank = levelRank["info"]
		}

		if entryRank < minLevel {
			continue
		}

		// Output based on format
		if *format == "json" {
			output, err := json.Marshal(entry)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(output))
		} else {
			fmt.Printf("time=\"%s\" level=\"%s\" msg=\"%s\"\n", entry.Timestamp, entry.Level, entry.Message)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

func parseLine(line string) LogEntry {
	entry := LogEntry{
		Fields: make(map[string]interface{}),
	}

	// Detect level from line content
	lineUpper := strings.ToUpper(line)
	if strings.Contains(lineUpper, "ERROR") {
		entry.Level = "error"
	} else if strings.Contains(lineUpper, "WARN") {
		entry.Level = "warn"
	} else if strings.Contains(lineUpper, "DEBUG") {
		entry.Level = "debug"
	} else {
		entry.Level = "info"
	}

	// Try to parse as JSON first
	if strings.HasPrefix(strings.TrimSpace(line), "{") {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonData); err == nil {
			// Extract common fields
			if ts, ok := jsonData["timestamp"].(string); ok {
				entry.Timestamp = ts
			} else if ts, ok := jsonData["time"].(string); ok {
				entry.Timestamp = ts
			} else if ts, ok := jsonData["ts"].(string); ok {
				entry.Timestamp = ts
			}

			if lvl, ok := jsonData["level"].(string); ok {
				entry.Level = strings.ToLower(lvl)
			} else if lvl, ok := jsonData["lvl"].(string); ok {
				entry.Level = strings.ToLower(lvl)
			}

			if msg, ok := jsonData["message"].(string); ok {
				entry.Message = msg
			} else if msg, ok := jsonData["msg"].(string); ok {
				entry.Message = msg
			}

			// Store remaining fields
			for k, v := range jsonData {
				if k != "timestamp" && k != "time" && k != "ts" &&
					k != "level" && k != "lvl" &&
					k != "message" && k != "msg" {
					entry.Fields[k] = v
				}
			}
			return entry
		}
	}

	// Parse as logfmt or plain text
	parts := strings.SplitN(line, " ", 2)
	if len(parts) > 0 {
		// Check if first part looks like a timestamp
		if strings.Contains(parts[0], "T") || strings.Contains(parts[0], "-") {
			entry.Timestamp = parts[0]
			if len(parts) > 1 {
				entry.Message = parts[1]
			}
		} else {
			entry.Message = line
		}
	} else {
		entry.Message = line
	}

	// Try to extract key=value pairs for logfmt
	parseLogfmtFields(line, entry.Fields)

	return entry
}

func parseLogfmtFields(line string, fields map[string]interface{}) {
	// Simple logfmt parser for key="value" or key=value patterns
	parts := strings.Fields(line)
	for _, part := range parts {
		if idx := strings.Index(part, "="); idx > 0 {
			key := part[:idx]
			value := part[idx+1:]
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			fields[key] = value
		}
	}
}
