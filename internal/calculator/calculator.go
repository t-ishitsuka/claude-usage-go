package calculator

import (
	"sort"
	"time"

	"github.com/t-ishitsuka/claude-usage-go/internal/models"
)

func CalculateCost(usage models.TokenUsage, model string) float64 {
	pricing, ok := models.ModelPricing[model]
	if !ok {
		return 0
	}

	cost := 0.0
	cost += float64(usage.InputTokens) / 1_000_000 * pricing.InputPer1M
	cost += float64(usage.OutputTokens) / 1_000_000 * pricing.OutputPer1M
	cost += float64(usage.CacheCreateTokens) / 1_000_000 * pricing.CacheCreatePer1M
	cost += float64(usage.CacheReadTokens) / 1_000_000 * pricing.CacheReadPer1M

	return cost
}

func AggregateDaily(messages []models.Message) []models.DailyUsage {
	dailyMap := make(map[string]*models.DailyUsage)

	for _, msg := range messages {
		dateKey := msg.Timestamp.Format("2006-01-02")

		if _, exists := dailyMap[dateKey]; !exists {
			dailyMap[dateKey] = &models.DailyUsage{
				Date:   msg.Timestamp.Truncate(24 * time.Hour),
				Models: make([]string, 0),
			}
		}

		daily := dailyMap[dateKey]
		daily.TokenUsage.InputTokens += msg.TokenUsage.InputTokens
		daily.TokenUsage.OutputTokens += msg.TokenUsage.OutputTokens
		daily.TokenUsage.CacheCreateTokens += msg.TokenUsage.CacheCreateTokens
		daily.TokenUsage.CacheReadTokens += msg.TokenUsage.CacheReadTokens

		msg.EstimatedCostUSD = CalculateCost(msg.TokenUsage, msg.Model)
		daily.CostUSD += msg.EstimatedCostUSD

		if !contains(daily.Models, msg.Model) {
			daily.Models = append(daily.Models, msg.Model)
		}
	}

	var result []models.DailyUsage
	for _, daily := range dailyMap {
		result = append(result, *daily)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return result
}

func AggregateMonthly(messages []models.Message) []models.MonthlyUsage {
	monthlyMap := make(map[string]*models.MonthlyUsage)

	for _, msg := range messages {
		monthKey := msg.Timestamp.Format("2006-01")

		if _, exists := monthlyMap[monthKey]; !exists {
			monthlyMap[monthKey] = &models.MonthlyUsage{
				Year:   msg.Timestamp.Year(),
				Month:  msg.Timestamp.Month(),
				Models: make([]string, 0),
			}
		}

		monthly := monthlyMap[monthKey]
		monthly.TokenUsage.InputTokens += msg.TokenUsage.InputTokens
		monthly.TokenUsage.OutputTokens += msg.TokenUsage.OutputTokens
		monthly.TokenUsage.CacheCreateTokens += msg.TokenUsage.CacheCreateTokens
		monthly.TokenUsage.CacheReadTokens += msg.TokenUsage.CacheReadTokens

		msg.EstimatedCostUSD = CalculateCost(msg.TokenUsage, msg.Model)
		monthly.CostUSD += msg.EstimatedCostUSD

		if !contains(monthly.Models, msg.Model) {
			monthly.Models = append(monthly.Models, msg.Model)
		}
	}

	var result []models.MonthlyUsage
	for _, monthly := range monthlyMap {
		result = append(result, *monthly)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Year != result[j].Year {
			return result[i].Year < result[j].Year
		}
		return result[i].Month < result[j].Month
	})

	return result
}

func AggregateBySession(messages []models.Message) []models.SessionUsage {
	sessionMap := make(map[string]*models.SessionUsage)

	for _, msg := range messages {
		if _, exists := sessionMap[msg.SessionID]; !exists {
			sessionMap[msg.SessionID] = &models.SessionUsage{
				SessionID: msg.SessionID,
				StartTime: msg.Timestamp,
				EndTime:   msg.Timestamp,
				Models:    make([]string, 0),
			}
		}

		session := sessionMap[msg.SessionID]
		session.TokenUsage.InputTokens += msg.TokenUsage.InputTokens
		session.TokenUsage.OutputTokens += msg.TokenUsage.OutputTokens
		session.TokenUsage.CacheCreateTokens += msg.TokenUsage.CacheCreateTokens
		session.TokenUsage.CacheReadTokens += msg.TokenUsage.CacheReadTokens

		msg.EstimatedCostUSD = CalculateCost(msg.TokenUsage, msg.Model)
		session.CostUSD += msg.EstimatedCostUSD

		if msg.Timestamp.Before(session.StartTime) {
			session.StartTime = msg.Timestamp
		}
		if msg.Timestamp.After(session.EndTime) {
			session.EndTime = msg.Timestamp
		}

		if !contains(session.Models, msg.Model) {
			session.Models = append(session.Models, msg.Model)
		}
	}

	var result []models.SessionUsage
	for _, session := range sessionMap {
		result = append(result, *session)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].StartTime.Before(result[j].StartTime)
	})

	return result
}

func AggregateByModel(messages []models.Message) []models.ModelBreakdown {
	modelMap := make(map[string]*models.ModelBreakdown)

	for _, msg := range messages {
		if _, exists := modelMap[msg.Model]; !exists {
			modelMap[msg.Model] = &models.ModelBreakdown{
				Model: msg.Model,
			}
		}

		breakdown := modelMap[msg.Model]
		breakdown.TokenUsage.InputTokens += msg.TokenUsage.InputTokens
		breakdown.TokenUsage.OutputTokens += msg.TokenUsage.OutputTokens
		breakdown.TokenUsage.CacheCreateTokens += msg.TokenUsage.CacheCreateTokens
		breakdown.TokenUsage.CacheReadTokens += msg.TokenUsage.CacheReadTokens
		breakdown.CostUSD += CalculateCost(msg.TokenUsage, msg.Model)
	}

	var result []models.ModelBreakdown
	for _, breakdown := range modelMap {
		result = append(result, *breakdown)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CostUSD > result[j].CostUSD
	})

	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
