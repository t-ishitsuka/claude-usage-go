package display

import (
	"testing"
)

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "Zero returns dash",
			input:    0,
			expected: "-",
		},
		{
			name:     "Positive number",
			input:    12345,
			expected: "12345",
		},
		{
			name:     "Large number",
			input:    1000000,
			expected: "1000000",
		},
		{
			name:     "Negative number",
			input:    -100,
			expected: "-100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatNumber(tt.input)
			if result != tt.expected {
				t.Errorf("formatNumber(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetShortModelNames(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Empty list",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "Single model",
			input:    []string{"claude-opus-4-20250514"},
			expected: []string{"Opus 4"},
		},
		{
			name:     "Multiple different models",
			input:    []string{"claude-opus-4-20250514", "claude-sonnet-4-20250514"},
			expected: []string{"Opus 4", "Sonnet 4"},
		},
		{
			name:     "Duplicate models - should dedupe",
			input:    []string{"claude-opus-4-20250514", "claude-opus-4-20250514"},
			expected: []string{"Opus 4"},
		},
		{
			name:     "Multiple Sonnet 3.5 versions - should dedupe to single short name",
			input:    []string{"claude-3-5-sonnet-20241022", "claude-3-5-sonnet-20240620"},
			expected: []string{"Sonnet 3.5"},
		},
		{
			name:     "Unknown model passes through",
			input:    []string{"unknown-model"},
			expected: []string{"unknown-model"},
		},
		{
			name:     "Mixed known and unknown",
			input:    []string{"claude-opus-4-20250514", "unknown-model", "claude-sonnet-4-20250514"},
			expected: []string{"Opus 4", "unknown-model", "Sonnet 4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getShortModelNames(tt.input)

			if len(result) != len(tt.expected) {
				t.Fatalf("getShortModelNames() returned %d items, want %d", len(result), len(tt.expected))
			}

			// Create a map to check existence rather than order
			// since the function doesn't guarantee order for all cases
			resultMap := make(map[string]bool)
			for _, name := range result {
				resultMap[name] = true
			}

			for _, expectedName := range tt.expected {
				if !resultMap[expectedName] {
					t.Errorf("Expected model name %s not found in result %v", expectedName, result)
				}
			}
		})
	}
}

func TestGetShortModelNames_PreservesOrderForUnique(t *testing.T) {
	// When all models are unique, order should be preserved
	input := []string{
		"claude-sonnet-4-20250514",
		"claude-opus-4-20250514",
		"claude-3-haiku-20240307",
	}

	result := getShortModelNames(input)
	expected := []string{"Sonnet 4", "Opus 4", "Haiku 3"}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(result))
	}

	for i, name := range result {
		if name != expected[i] {
			t.Errorf("At index %d: got %s, want %s", i, name, expected[i])
		}
	}
}
