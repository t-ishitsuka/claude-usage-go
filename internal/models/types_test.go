package models

import (
	"testing"
	"time"
)

func TestTokenUsage_Total(t *testing.T) {
	tests := []struct {
		name     string
		usage    TokenUsage
		expected int
	}{
		{
			name: "All token types",
			usage: TokenUsage{
				InputTokens:       100,
				OutputTokens:      200,
				CacheCreateTokens: 50,
				CacheReadTokens:   25,
			},
			expected: 375,
		},
		{
			name: "Only input and output",
			usage: TokenUsage{
				InputTokens:  100,
				OutputTokens: 200,
			},
			expected: 300,
		},
		{
			name:     "Empty usage",
			usage:    TokenUsage{},
			expected: 0,
		},
		{
			name: "Large numbers",
			usage: TokenUsage{
				InputTokens:       1000000,
				OutputTokens:      2000000,
				CacheCreateTokens: 500000,
				CacheReadTokens:   250000,
			},
			expected: 3750000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.usage.Total()
			if result != tt.expected {
				t.Errorf("Total() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestReportOptions(t *testing.T) {
	// Test that ReportOptions can be properly initialized
	now := time.Now()
	opts := ReportOptions{
		Since:      &now,
		Until:      &now,
		Breakdown:  true,
		JSONOutput: true,
		Ascending:  true,
		Models:     []string{"model1", "model2"},
	}

	if opts.Since == nil || opts.Until == nil {
		t.Error("Time pointers should not be nil")
	}
	if !opts.Breakdown || !opts.JSONOutput || !opts.Ascending {
		t.Error("Boolean flags should be true")
	}
	if len(opts.Models) != 2 {
		t.Errorf("Expected 2 models, got %d", len(opts.Models))
	}
}