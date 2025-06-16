package models

var ModelPricing = map[string]struct {
	InputPer1M       float64
	OutputPer1M      float64
	CacheCreatePer1M float64
	CacheReadPer1M   float64
}{
	"claude-opus-4-20250514": {
		InputPer1M:       15.00,
		OutputPer1M:      75.00,
		CacheCreatePer1M: 18.75,
		CacheReadPer1M:   1.50,
	},
	"claude-sonnet-4-20250514": {
		InputPer1M:       3.00,
		OutputPer1M:      15.00,
		CacheCreatePer1M: 3.75,
		CacheReadPer1M:   0.30,
	},
	"claude-3-5-sonnet-20241022": {
		InputPer1M:       3.00,
		OutputPer1M:      15.00,
		CacheCreatePer1M: 3.75,
		CacheReadPer1M:   0.30,
	},
	"claude-3-5-sonnet-20240620": {
		InputPer1M:       3.00,
		OutputPer1M:      15.00,
		CacheCreatePer1M: 3.75,
		CacheReadPer1M:   0.30,
	},
	"claude-3-5-haiku-20241022": {
		InputPer1M:       0.80,
		OutputPer1M:      4.00,
		CacheCreatePer1M: 1.00,
		CacheReadPer1M:   0.08,
	},
	"claude-3-opus-20240229": {
		InputPer1M:       15.00,
		OutputPer1M:      75.00,
		CacheCreatePer1M: 18.75,
		CacheReadPer1M:   1.50,
	},
	"claude-3-sonnet-20240229": {
		InputPer1M:       3.00,
		OutputPer1M:      15.00,
		CacheCreatePer1M: 3.75,
		CacheReadPer1M:   0.30,
	},
	"claude-3-haiku-20240307": {
		InputPer1M:       0.25,
		OutputPer1M:      1.25,
		CacheCreatePer1M: 0.30,
		CacheReadPer1M:   0.03,
	},
}

func GetModelShortName(model string) string {
	shortNames := map[string]string{
		"claude-opus-4-20250514":     "Opus 4",
		"claude-sonnet-4-20250514":   "Sonnet 4",
		"claude-3-5-sonnet-20241022": "Sonnet 3.5",
		"claude-3-5-sonnet-20240620": "Sonnet 3.5",
		"claude-3-5-haiku-20241022":  "Haiku 3.5",
		"claude-3-opus-20240229":     "Opus 3",
		"claude-3-sonnet-20240229":   "Sonnet 3",
		"claude-3-haiku-20240307":    "Haiku 3",
	}
	if short, ok := shortNames[model]; ok {
		return short
	}
	return model
}
