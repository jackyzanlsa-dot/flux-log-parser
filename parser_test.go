package main

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedLevel string
	}{
		{
			name:          "ERROR level detection",
			input:         `{"timestamp":"2024-01-15T10:30:00Z","level":"ERROR","message":"database connection failed"}`,
			expectedLevel: "error",
		},
		{
			name:          "WARN level detection",
			input:         `{"timestamp":"2024-01-15T10:30:00Z","level":"WARN","message":"high memory usage"}`,
			expectedLevel: "warn",
		},
		{
			name:          "DEBUG level detection",
			input:         `{"timestamp":"2024-01-15T10:30:00Z","level":"DEBUG","message":"entering function foo"}`,
			expectedLevel: "debug",
		},
		{
			name:          "INFO level detection (plain line)",
			input:         `2024-01-15T10:30:00Z Application started successfully`,
			expectedLevel: "info",
		},
		{
			name:          "ERROR in message content",
			input:         `2024-01-15T10:30:00Z ERROR: something went wrong`,
			expectedLevel: "error",
		},
		{
			name:          "WARN in message content",
			input:         `2024-01-15T10:30:00Z WARN: deprecated API call`,
			expectedLevel: "warn",
		},
		{
			name:          "DEBUG in message content",
			input:         `2024-01-15T10:30:00Z DEBUG: variable x = 42`,
			expectedLevel: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := parseLine(tt.input)
			if entry.Level != tt.expectedLevel {
				t.Errorf("parseLine(%q) returned Level = %q, want %q", tt.input, entry.Level, tt.expectedLevel)
			}
		})
	}
}
