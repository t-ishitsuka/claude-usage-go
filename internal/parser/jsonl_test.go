package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/t-ishitsuka/claude-usage-go/internal/models"
)

func TestParseJSONLFiles(t *testing.T) {
	// Create temporary directory and test files
	tempDir, err := os.MkdirTemp("", "claude-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test JSONL file with valid data
	testJSONL := `{"sessionId":"test-session","timestamp":"2025-01-15T10:00:00.000Z","type":"assistant","message":{"role":"assistant","model":"claude-opus-4-20250514","usage":{"input_tokens":100,"output_tokens":200,"cache_creation_input_tokens":50,"cache_read_input_tokens":25}}}
{"sessionId":"test-session","timestamp":"2025-01-15T11:00:00.000Z","type":"user","message":{"role":"user"}}
{"sessionId":"test-session","timestamp":"2025-01-15T11:01:00.000Z","type":"assistant","message":{"role":"assistant","model":"claude-sonnet-4-20250514","usage":{"input_tokens":150,"output_tokens":300,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}}
{"type":"summary","summary":"Test summary"}
`

	testFile := filepath.Join(tempDir, "test.jsonl")
	if err := os.WriteFile(testFile, []byte(testJSONL), 0644); err != nil {
		t.Fatal(err)
	}

	// Test parsing
	messages, err := ParseJSONLFiles(tempDir)
	if err != nil {
		t.Fatalf("ParseJSONLFiles() error = %v", err)
	}

	if len(messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(messages))
	}

	// Check first message
	msg1 := messages[0]
	if msg1.SessionID != "test-session" {
		t.Errorf("Message 1 SessionID = %s, want test-session", msg1.SessionID)
	}
	if msg1.Model != "claude-opus-4-20250514" {
		t.Errorf("Message 1 Model = %s, want claude-opus-4-20250514", msg1.Model)
	}
	if msg1.TokenUsage.InputTokens != 100 {
		t.Errorf("Message 1 InputTokens = %d, want 100", msg1.TokenUsage.InputTokens)
	}
	if msg1.TokenUsage.OutputTokens != 200 {
		t.Errorf("Message 1 OutputTokens = %d, want 200", msg1.TokenUsage.OutputTokens)
	}
	if msg1.TokenUsage.CacheCreateTokens != 50 {
		t.Errorf("Message 1 CacheCreateTokens = %d, want 50", msg1.TokenUsage.CacheCreateTokens)
	}
	if msg1.TokenUsage.CacheReadTokens != 25 {
		t.Errorf("Message 1 CacheReadTokens = %d, want 25", msg1.TokenUsage.CacheReadTokens)
	}

	// Check second message
	msg2 := messages[1]
	if msg2.Model != "claude-sonnet-4-20250514" {
		t.Errorf("Message 2 Model = %s, want claude-sonnet-4-20250514", msg2.Model)
	}
	if msg2.TokenUsage.InputTokens != 150 {
		t.Errorf("Message 2 InputTokens = %d, want 150", msg2.TokenUsage.InputTokens)
	}
}

func TestParseJSONLFiles_EmptyDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "claude-test-empty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	messages, err := ParseJSONLFiles(tempDir)
	if err != nil {
		t.Fatalf("ParseJSONLFiles() error = %v", err)
	}

	if len(messages) != 0 {
		t.Errorf("Expected 0 messages for empty directory, got %d", len(messages))
	}
}

func TestParseJSONLFiles_InvalidJSON(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "claude-test-invalid")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file with invalid JSON
	testJSONL := `invalid json content
{"sessionId":"test-session","timestamp":"2025-01-15T10:00:00.000Z","type":"assistant","message":{"role":"assistant","model":"claude-opus-4-20250514","usage":{"input_tokens":100,"output_tokens":200,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}}
`

	testFile := filepath.Join(tempDir, "test.jsonl")
	if err := os.WriteFile(testFile, []byte(testJSONL), 0644); err != nil {
		t.Fatal(err)
	}

	messages, err := ParseJSONLFiles(tempDir)
	if err != nil {
		t.Fatalf("ParseJSONLFiles() error = %v", err)
	}

	// Should still parse the valid line
	if len(messages) != 1 {
		t.Errorf("Expected 1 message (skipping invalid), got %d", len(messages))
	}
}

func TestFilterByDateRange(t *testing.T) {
	baseTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	messages := []models.Message{
		{Timestamp: baseTime.Add(-2 * 24 * time.Hour)}, // Jan 13
		{Timestamp: baseTime},                           // Jan 15
		{Timestamp: baseTime.Add(2 * 24 * time.Hour)},  // Jan 17
		{Timestamp: baseTime.Add(5 * 24 * time.Hour)},  // Jan 20
	}

	tests := []struct {
		name     string
		since    *time.Time
		until    *time.Time
		expected int
	}{
		{
			name:     "No filters",
			since:    nil,
			until:    nil,
			expected: 4,
		},
		{
			name:     "Since only",
			since:    &baseTime,
			until:    nil,
			expected: 3, // Jan 15, 17, 20
		},
		{
			name:     "Until only",
			since:    nil,
			until:    &baseTime,
			expected: 2, // Jan 13, 15
		},
		{
			name: "Both filters",
			since: func() *time.Time {
				t := baseTime.Add(-1 * 24 * time.Hour)
				return &t
			}(),
			until: func() *time.Time {
				t := baseTime.Add(1 * 24 * time.Hour)
				return &t
			}(),
			expected: 2, // Jan 15, 17 (until is inclusive of the next day)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterByDateRange(messages, tt.since, tt.until)
			if len(result) != tt.expected {
				t.Errorf("FilterByDateRange() returned %d messages, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestFilterByModels(t *testing.T) {
	messages := []models.Message{
		{Model: "claude-opus-4-20250514"},
		{Model: "claude-sonnet-4-20250514"},
		{Model: "claude-3-5-sonnet-20241022"},
		{Model: "claude-opus-4-20250514"},
	}

	tests := []struct {
		name     string
		models   []string
		expected int
	}{
		{
			name:     "No filter",
			models:   []string{},
			expected: 4,
		},
		{
			name:     "Single model",
			models:   []string{"claude-opus-4-20250514"},
			expected: 2,
		},
		{
			name:     "Multiple models",
			models:   []string{"claude-opus-4-20250514", "claude-sonnet-4-20250514"},
			expected: 3,
		},
		{
			name:     "Case insensitive",
			models:   []string{"CLAUDE-OPUS-4-20250514"},
			expected: 2,
		},
		{
			name:     "Non-existent model",
			models:   []string{"non-existent-model"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterByModels(messages, tt.models)
			if len(result) != tt.expected {
				t.Errorf("FilterByModels() returned %d messages, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestGetClaudeProjectsDir(t *testing.T) {
	result := GetClaudeProjectsDir()
	
	// Should end with .claude/projects
	expected := filepath.Join(".claude", "projects")
	if !filepath.IsAbs(result) {
		t.Errorf("GetClaudeProjectsDir() should return absolute path, got %s", result)
	}
	if !strings.HasSuffix(result, expected) {
		t.Errorf("GetClaudeProjectsDir() should end with %s, got %s", expected, result)
	}
}