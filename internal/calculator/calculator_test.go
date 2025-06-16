package calculator

import (
	"testing"
	"time"

	"github.com/t-ishitsuka/claude-usage-go/internal/models"
)

func TestCalculateCost(t *testing.T) {
	tests := []struct {
		name     string
		usage    models.TokenUsage
		model    string
		expected float64
	}{
		{
			name: "Opus 4 with all token types",
			usage: models.TokenUsage{
				InputTokens:       1000000,
				OutputTokens:      1000000,
				CacheCreateTokens: 1000000,
				CacheReadTokens:   1000000,
			},
			model:    "claude-opus-4-20250514",
			expected: 15.0 + 75.0 + 18.75 + 1.5, // 110.25
		},
		{
			name: "Sonnet 4 with partial tokens",
			usage: models.TokenUsage{
				InputTokens:  500000,
				OutputTokens: 500000,
			},
			model:    "claude-sonnet-4-20250514",
			expected: 1.5 + 7.5, // 9.0
		},
		{
			name: "Unknown model returns 0",
			usage: models.TokenUsage{
				InputTokens: 1000000,
			},
			model:    "unknown-model",
			expected: 0,
		},
		{
			name: "Zero tokens",
			usage: models.TokenUsage{
				InputTokens:       0,
				OutputTokens:      0,
				CacheCreateTokens: 0,
				CacheReadTokens:   0,
			},
			model:    "claude-opus-4-20250514",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateCost(tt.usage, tt.model)
			if result != tt.expected {
				t.Errorf("CalculateCost() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAggregateDaily(t *testing.T) {
	baseTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	messages := []models.Message{
		{
			SessionID: "session1",
			Timestamp: baseTime,
			Model:     "claude-opus-4-20250514",
			TokenUsage: models.TokenUsage{
				InputTokens:  1000,
				OutputTokens: 2000,
			},
		},
		{
			SessionID: "session2",
			Timestamp: baseTime.Add(1 * time.Hour),
			Model:     "claude-opus-4-20250514",
			TokenUsage: models.TokenUsage{
				InputTokens:  500,
				OutputTokens: 1500,
			},
		},
		{
			SessionID: "session3",
			Timestamp: baseTime.Add(24 * time.Hour),
			Model:     "claude-sonnet-4-20250514",
			TokenUsage: models.TokenUsage{
				InputTokens:  2000,
				OutputTokens: 3000,
			},
		},
	}

	result := AggregateDaily(messages)

	if len(result) != 2 {
		t.Fatalf("Expected 2 daily aggregations, got %d", len(result))
	}

	// Check first day
	day1 := result[0]
	if day1.TokenUsage.InputTokens != 1500 {
		t.Errorf("Day 1 input tokens = %d, want 1500", day1.TokenUsage.InputTokens)
	}
	if day1.TokenUsage.OutputTokens != 3500 {
		t.Errorf("Day 1 output tokens = %d, want 3500", day1.TokenUsage.OutputTokens)
	}
	if len(day1.Models) != 1 || day1.Models[0] != "claude-opus-4-20250514" {
		t.Errorf("Day 1 models = %v, want [claude-opus-4-20250514]", day1.Models)
	}

	// Check second day
	day2 := result[1]
	if day2.TokenUsage.InputTokens != 2000 {
		t.Errorf("Day 2 input tokens = %d, want 2000", day2.TokenUsage.InputTokens)
	}
	if day2.TokenUsage.OutputTokens != 3000 {
		t.Errorf("Day 2 output tokens = %d, want 3000", day2.TokenUsage.OutputTokens)
	}

	// Check ordering (should be ascending by date)
	if !result[0].Date.Before(result[1].Date) {
		t.Error("Results are not sorted by date in ascending order")
	}
}

func TestAggregateMonthly(t *testing.T) {
	messages := []models.Message{
		{
			SessionID: "session1",
			Timestamp: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			Model:     "claude-opus-4-20250514",
			TokenUsage: models.TokenUsage{
				InputTokens:  1000,
				OutputTokens: 2000,
			},
		},
		{
			SessionID: "session2",
			Timestamp: time.Date(2025, 1, 20, 10, 0, 0, 0, time.UTC),
			Model:     "claude-sonnet-4-20250514",
			TokenUsage: models.TokenUsage{
				InputTokens:  500,
				OutputTokens: 1500,
			},
		},
		{
			SessionID: "session3",
			Timestamp: time.Date(2025, 2, 10, 10, 0, 0, 0, time.UTC),
			Model:     "claude-opus-4-20250514",
			TokenUsage: models.TokenUsage{
				InputTokens:  2000,
				OutputTokens: 3000,
			},
		},
	}

	result := AggregateMonthly(messages)

	if len(result) != 2 {
		t.Fatalf("Expected 2 monthly aggregations, got %d", len(result))
	}

	// Check January
	jan := result[0]
	if jan.Year != 2025 || jan.Month != time.January {
		t.Errorf("First month = %d-%d, want 2025-1", jan.Year, jan.Month)
	}
	if jan.TokenUsage.InputTokens != 1500 {
		t.Errorf("January input tokens = %d, want 1500", jan.TokenUsage.InputTokens)
	}
	if len(jan.Models) != 2 {
		t.Errorf("January models count = %d, want 2", len(jan.Models))
	}

	// Check February
	feb := result[1]
	if feb.Year != 2025 || feb.Month != time.February {
		t.Errorf("Second month = %d-%d, want 2025-2", feb.Year, feb.Month)
	}
}

func TestAggregateBySession(t *testing.T) {
	baseTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	messages := []models.Message{
		{
			SessionID:  "session1",
			Timestamp:  baseTime,
			Model:      "claude-opus-4-20250514",
			TokenUsage: models.TokenUsage{InputTokens: 1000, OutputTokens: 2000},
		},
		{
			SessionID:  "session1",
			Timestamp:  baseTime.Add(1 * time.Hour),
			Model:      "claude-sonnet-4-20250514",
			TokenUsage: models.TokenUsage{InputTokens: 500, OutputTokens: 1500},
		},
		{
			SessionID:  "session2",
			Timestamp:  baseTime.Add(2 * time.Hour),
			Model:      "claude-opus-4-20250514",
			TokenUsage: models.TokenUsage{InputTokens: 2000, OutputTokens: 3000},
		},
	}

	result := AggregateBySession(messages)

	if len(result) != 2 {
		t.Fatalf("Expected 2 session aggregations, got %d", len(result))
	}

	// Check session1
	session1 := result[0]
	if session1.SessionID != "session1" {
		t.Errorf("First session ID = %s, want session1", session1.SessionID)
	}
	if session1.TokenUsage.InputTokens != 1500 {
		t.Errorf("Session1 input tokens = %d, want 1500", session1.TokenUsage.InputTokens)
	}
	if len(session1.Models) != 2 {
		t.Errorf("Session1 models count = %d, want 2", len(session1.Models))
	}
	if session1.StartTime != baseTime {
		t.Errorf("Session1 start time = %v, want %v", session1.StartTime, baseTime)
	}
	if session1.EndTime != baseTime.Add(1*time.Hour) {
		t.Errorf("Session1 end time = %v, want %v", session1.EndTime, baseTime.Add(1*time.Hour))
	}
}

func TestAggregateByModel(t *testing.T) {
	messages := []models.Message{
		{
			SessionID:  "session1",
			Timestamp:  time.Now(),
			Model:      "claude-opus-4-20250514",
			TokenUsage: models.TokenUsage{InputTokens: 1000, OutputTokens: 2000},
		},
		{
			SessionID:  "session2",
			Timestamp:  time.Now(),
			Model:      "claude-opus-4-20250514",
			TokenUsage: models.TokenUsage{InputTokens: 500, OutputTokens: 1500},
		},
		{
			SessionID:  "session3",
			Timestamp:  time.Now(),
			Model:      "claude-sonnet-4-20250514",
			TokenUsage: models.TokenUsage{InputTokens: 2000, OutputTokens: 3000},
		},
	}

	result := AggregateByModel(messages)

	if len(result) != 2 {
		t.Fatalf("Expected 2 model aggregations, got %d", len(result))
	}

	// Should be sorted by cost (descending)
	if result[0].CostUSD < result[1].CostUSD {
		t.Error("Results are not sorted by cost in descending order")
	}

	// Check Opus totals
	var opus *models.ModelBreakdown
	for i := range result {
		if result[i].Model == "claude-opus-4-20250514" {
			opus = &result[i]
			break
		}
	}
	if opus == nil {
		t.Fatal("Opus model not found in results")
	}
	if opus.TokenUsage.InputTokens != 1500 {
		t.Errorf("Opus input tokens = %d, want 1500", opus.TokenUsage.InputTokens)
	}
	if opus.TokenUsage.OutputTokens != 3500 {
		t.Errorf("Opus output tokens = %d, want 3500", opus.TokenUsage.OutputTokens)
	}
}
