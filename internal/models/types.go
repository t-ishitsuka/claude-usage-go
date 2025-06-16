package models

import (
	"time"
)

type TokenUsage struct {
	InputTokens       int
	OutputTokens      int
	CacheCreateTokens int
	CacheReadTokens   int
}

func (t TokenUsage) Total() int {
	return t.InputTokens + t.OutputTokens + t.CacheCreateTokens + t.CacheReadTokens
}

type Message struct {
	SessionID        string    `json:"session_id"`
	Timestamp        time.Time `json:"timestamp"`
	Model            string    `json:"model"`
	TokenUsage       TokenUsage
	EstimatedCostUSD float64
}

type DailyUsage struct {
	Date       time.Time
	Models     []string
	TokenUsage TokenUsage
	CostUSD    float64
}

type MonthlyUsage struct {
	Year       int
	Month      time.Month
	Models     []string
	TokenUsage TokenUsage
	CostUSD    float64
}

type SessionUsage struct {
	SessionID  string
	StartTime  time.Time
	EndTime    time.Time
	Models     []string
	TokenUsage TokenUsage
	CostUSD    float64
}

type ModelBreakdown struct {
	Model      string
	TokenUsage TokenUsage
	CostUSD    float64
}

type ReportOptions struct {
	Since      *time.Time
	Until      *time.Time
	Breakdown  bool
	JSONOutput bool
	Ascending  bool
	Models     []string
}
