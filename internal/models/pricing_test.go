package models

import (
	"testing"
)

func TestModelPricing(t *testing.T) {
	// Test that all expected models have pricing
	expectedModels := []string{
		"claude-opus-4-20250514",
		"claude-sonnet-4-20250514",
		"claude-3-5-sonnet-20241022",
		"claude-3-5-sonnet-20240620",
		"claude-3-5-haiku-20241022",
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
	}

	for _, model := range expectedModels {
		pricing, exists := ModelPricing[model]
		if !exists {
			t.Errorf("Model %s missing from ModelPricing", model)
			continue
		}

		// Verify all prices are non-negative
		if pricing.InputPer1M < 0 {
			t.Errorf("Model %s has negative input price: %f", model, pricing.InputPer1M)
		}
		if pricing.OutputPer1M < 0 {
			t.Errorf("Model %s has negative output price: %f", model, pricing.OutputPer1M)
		}
		if pricing.CacheCreatePer1M < 0 {
			t.Errorf("Model %s has negative cache create price: %f", model, pricing.CacheCreatePer1M)
		}
		if pricing.CacheReadPer1M < 0 {
			t.Errorf("Model %s has negative cache read price: %f", model, pricing.CacheReadPer1M)
		}

		// Verify output is more expensive than input (common pattern)
		if pricing.OutputPer1M <= pricing.InputPer1M {
			t.Errorf("Model %s output price (%f) should be higher than input price (%f)",
				model, pricing.OutputPer1M, pricing.InputPer1M)
		}
	}
}

func TestModelPricing_SpecificValues(t *testing.T) {
	tests := []struct {
		model               string
		expectedInput       float64
		expectedOutput      float64
		expectedCacheCreate float64
		expectedCacheRead   float64
	}{
		{
			model:               "claude-opus-4-20250514",
			expectedInput:       15.00,
			expectedOutput:      75.00,
			expectedCacheCreate: 18.75,
			expectedCacheRead:   1.50,
		},
		{
			model:               "claude-sonnet-4-20250514",
			expectedInput:       3.00,
			expectedOutput:      15.00,
			expectedCacheCreate: 3.75,
			expectedCacheRead:   0.30,
		},
		{
			model:               "claude-3-haiku-20240307",
			expectedInput:       0.25,
			expectedOutput:      1.25,
			expectedCacheCreate: 0.30,
			expectedCacheRead:   0.03,
		},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			pricing, exists := ModelPricing[tt.model]
			if !exists {
				t.Fatalf("Model %s not found in ModelPricing", tt.model)
			}

			if pricing.InputPer1M != tt.expectedInput {
				t.Errorf("Input price = %f, want %f", pricing.InputPer1M, tt.expectedInput)
			}
			if pricing.OutputPer1M != tt.expectedOutput {
				t.Errorf("Output price = %f, want %f", pricing.OutputPer1M, tt.expectedOutput)
			}
			if pricing.CacheCreatePer1M != tt.expectedCacheCreate {
				t.Errorf("Cache create price = %f, want %f", pricing.CacheCreatePer1M, tt.expectedCacheCreate)
			}
			if pricing.CacheReadPer1M != tt.expectedCacheRead {
				t.Errorf("Cache read price = %f, want %f", pricing.CacheReadPer1M, tt.expectedCacheRead)
			}
		})
	}
}

func TestGetModelShortName(t *testing.T) {
	tests := []struct {
		model    string
		expected string
	}{
		{"claude-opus-4-20250514", "Opus 4"},
		{"claude-sonnet-4-20250514", "Sonnet 4"},
		{"claude-3-5-sonnet-20241022", "Sonnet 3.5"},
		{"claude-3-5-sonnet-20240620", "Sonnet 3.5"},
		{"claude-3-5-haiku-20241022", "Haiku 3.5"},
		{"claude-3-opus-20240229", "Opus 3"},
		{"claude-3-sonnet-20240229", "Sonnet 3"},
		{"claude-3-haiku-20240307", "Haiku 3"},
		{"unknown-model-12345", "unknown-model-12345"}, // Unknown models return as-is
		{"", ""}, // Empty string returns empty
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			result := GetModelShortName(tt.model)
			if result != tt.expected {
				t.Errorf("GetModelShortName(%s) = %s, want %s", tt.model, result, tt.expected)
			}
		})
	}
}
