package display

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/t-ishitsuka/claude-usage-go/internal/calculator"
	"github.com/t-ishitsuka/claude-usage-go/internal/models"
)

var (
	headerColor = color.New(color.FgCyan, color.Bold)
	totalColor  = color.New(color.FgYellow, color.Bold)
	modelColor  = color.New(color.FgGreen)
	costColor   = color.New(color.FgRed)
)

func ShowDaily(dailyUsage []models.DailyUsage, ascending bool) error {
	if ascending {
		sort.Slice(dailyUsage, func(i, j int) bool {
			return dailyUsage[i].Date.After(dailyUsage[j].Date)
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Models", "Input", "Output", "Cache Create", "Cache Read", "Total", "Cost (USD)"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetCenterSeparator("│")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
	)

	var totalUsage models.TokenUsage
	var totalCost float64

	for _, daily := range dailyUsage {
		modelNames := getShortModelNames(daily.Models)
		table.Append([]string{
			daily.Date.Format("2006-01-02"),
			strings.Join(modelNames, ", "),
			formatNumber(daily.TokenUsage.InputTokens),
			formatNumber(daily.TokenUsage.OutputTokens),
			formatNumber(daily.TokenUsage.CacheCreateTokens),
			formatNumber(daily.TokenUsage.CacheReadTokens),
			formatNumber(daily.TokenUsage.Total()),
			fmt.Sprintf("$%.4f", daily.CostUSD),
		})

		totalUsage.InputTokens += daily.TokenUsage.InputTokens
		totalUsage.OutputTokens += daily.TokenUsage.OutputTokens
		totalUsage.CacheCreateTokens += daily.TokenUsage.CacheCreateTokens
		totalUsage.CacheReadTokens += daily.TokenUsage.CacheReadTokens
		totalCost += daily.CostUSD
	}

	table.SetFooter([]string{
		"TOTAL", "",
		formatNumber(totalUsage.InputTokens),
		formatNumber(totalUsage.OutputTokens),
		formatNumber(totalUsage.CacheCreateTokens),
		formatNumber(totalUsage.CacheReadTokens),
		formatNumber(totalUsage.Total()),
		fmt.Sprintf("$%.4f", totalCost),
	})
	table.SetFooterColor(
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
	)

	table.Render()
	return nil
}

func ShowDailyWithBreakdown(dailyUsage []models.DailyUsage, messages []models.Message, ascending bool) error {
	if ascending {
		sort.Slice(dailyUsage, func(i, j int) bool {
			return dailyUsage[i].Date.After(dailyUsage[j].Date)
		})
	}

	for _, daily := range dailyUsage {
		fmt.Printf("\n%s %s\n", headerColor.Sprint("Date:"), daily.Date.Format("2006-01-02"))
		
		var dayMessages []models.Message
		for _, msg := range messages {
			if msg.Timestamp.Format("2006-01-02") == daily.Date.Format("2006-01-02") {
				dayMessages = append(dayMessages, msg)
			}
		}

		showModelBreakdown(dayMessages)
	}

	fmt.Println("\n" + strings.Repeat("═", 80))
	ShowDaily(dailyUsage, ascending)
	return nil
}

func ShowMonthly(monthlyUsage []models.MonthlyUsage, ascending bool) error {
	if ascending {
		sort.Slice(monthlyUsage, func(i, j int) bool {
			if monthlyUsage[i].Year != monthlyUsage[j].Year {
				return monthlyUsage[i].Year > monthlyUsage[j].Year
			}
			return monthlyUsage[i].Month > monthlyUsage[j].Month
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Month", "Models", "Input", "Output", "Cache Create", "Cache Read", "Total", "Cost (USD)"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetCenterSeparator("│")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
	)

	var totalUsage models.TokenUsage
	var totalCost float64

	for _, monthly := range monthlyUsage {
		modelNames := getShortModelNames(monthly.Models)
		table.Append([]string{
			fmt.Sprintf("%d-%02d", monthly.Year, monthly.Month),
			strings.Join(modelNames, ", "),
			formatNumber(monthly.TokenUsage.InputTokens),
			formatNumber(monthly.TokenUsage.OutputTokens),
			formatNumber(monthly.TokenUsage.CacheCreateTokens),
			formatNumber(monthly.TokenUsage.CacheReadTokens),
			formatNumber(monthly.TokenUsage.Total()),
			fmt.Sprintf("$%.4f", monthly.CostUSD),
		})

		totalUsage.InputTokens += monthly.TokenUsage.InputTokens
		totalUsage.OutputTokens += monthly.TokenUsage.OutputTokens
		totalUsage.CacheCreateTokens += monthly.TokenUsage.CacheCreateTokens
		totalUsage.CacheReadTokens += monthly.TokenUsage.CacheReadTokens
		totalCost += monthly.CostUSD
	}

	table.SetFooter([]string{
		"TOTAL", "",
		formatNumber(totalUsage.InputTokens),
		formatNumber(totalUsage.OutputTokens),
		formatNumber(totalUsage.CacheCreateTokens),
		formatNumber(totalUsage.CacheReadTokens),
		formatNumber(totalUsage.Total()),
		fmt.Sprintf("$%.4f", totalCost),
	})
	table.SetFooterColor(
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
	)

	table.Render()
	return nil
}

func ShowMonthlyWithBreakdown(monthlyUsage []models.MonthlyUsage, messages []models.Message, ascending bool) error {
	if ascending {
		sort.Slice(monthlyUsage, func(i, j int) bool {
			if monthlyUsage[i].Year != monthlyUsage[j].Year {
				return monthlyUsage[i].Year > monthlyUsage[j].Year
			}
			return monthlyUsage[i].Month > monthlyUsage[j].Month
		})
	}

	for _, monthly := range monthlyUsage {
		fmt.Printf("\n%s %d-%02d\n", headerColor.Sprint("Month:"), monthly.Year, monthly.Month)
		
		var monthMessages []models.Message
		for _, msg := range messages {
			if msg.Timestamp.Year() == monthly.Year && msg.Timestamp.Month() == monthly.Month {
				monthMessages = append(monthMessages, msg)
			}
		}

		showModelBreakdown(monthMessages)
	}

	fmt.Println("\n" + strings.Repeat("═", 80))
	ShowMonthly(monthlyUsage, ascending)
	return nil
}

func ShowSession(sessionUsage []models.SessionUsage, ascending bool) error {
	if !ascending {
		sort.Slice(sessionUsage, func(i, j int) bool {
			return sessionUsage[i].StartTime.After(sessionUsage[j].StartTime)
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Session ID", "Start Time", "Models", "Input", "Output", "Cache Create", "Cache Read", "Total", "Cost (USD)"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetCenterSeparator("│")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Bold},
	)

	var totalUsage models.TokenUsage
	var totalCost float64

	for _, session := range sessionUsage {
		modelNames := getShortModelNames(session.Models)
		table.Append([]string{
			session.SessionID[:8] + "...",
			session.StartTime.Format("2006-01-02 15:04"),
			strings.Join(modelNames, ", "),
			formatNumber(session.TokenUsage.InputTokens),
			formatNumber(session.TokenUsage.OutputTokens),
			formatNumber(session.TokenUsage.CacheCreateTokens),
			formatNumber(session.TokenUsage.CacheReadTokens),
			formatNumber(session.TokenUsage.Total()),
			fmt.Sprintf("$%.4f", session.CostUSD),
		})

		totalUsage.InputTokens += session.TokenUsage.InputTokens
		totalUsage.OutputTokens += session.TokenUsage.OutputTokens
		totalUsage.CacheCreateTokens += session.TokenUsage.CacheCreateTokens
		totalUsage.CacheReadTokens += session.TokenUsage.CacheReadTokens
		totalCost += session.CostUSD
	}

	table.SetFooter([]string{
		"TOTAL", "", "",
		formatNumber(totalUsage.InputTokens),
		formatNumber(totalUsage.OutputTokens),
		formatNumber(totalUsage.CacheCreateTokens),
		formatNumber(totalUsage.CacheReadTokens),
		formatNumber(totalUsage.Total()),
		fmt.Sprintf("$%.4f", totalCost),
	})
	table.SetFooterColor(
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
	)

	table.Render()
	return nil
}

func ShowSessionWithBreakdown(sessionUsage []models.SessionUsage, messages []models.Message, ascending bool) error {
	if !ascending {
		sort.Slice(sessionUsage, func(i, j int) bool {
			return sessionUsage[i].StartTime.After(sessionUsage[j].StartTime)
		})
	}

	for _, session := range sessionUsage {
		fmt.Printf("\n%s %s\n", headerColor.Sprint("Session:"), session.SessionID[:16]+"...")
		fmt.Printf("%s %s - %s\n", headerColor.Sprint("Time Range:"), 
			session.StartTime.Format("2006-01-02 15:04"),
			session.EndTime.Format("15:04"))
		
		var sessionMessages []models.Message
		for _, msg := range messages {
			if msg.SessionID == session.SessionID {
				sessionMessages = append(sessionMessages, msg)
			}
		}

		showModelBreakdown(sessionMessages)
	}

	fmt.Println("\n" + strings.Repeat("═", 80))
	ShowSession(sessionUsage, ascending)
	return nil
}

func showModelBreakdown(messages []models.Message) {
	breakdown := calculator.AggregateByModel(messages)
	
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Model", "Input", "Output", "Cache Create", "Cache Read", "Total", "Cost (USD)"})
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColumnSeparator(" ")
	table.SetAlignment(tablewriter.ALIGN_RIGHT)

	for _, b := range breakdown {
		table.Append([]string{
			modelColor.Sprint(models.GetModelShortName(b.Model)),
			formatNumber(b.TokenUsage.InputTokens),
			formatNumber(b.TokenUsage.OutputTokens),
			formatNumber(b.TokenUsage.CacheCreateTokens),
			formatNumber(b.TokenUsage.CacheReadTokens),
			formatNumber(b.TokenUsage.Total()),
			costColor.Sprintf("$%.4f", b.CostUSD),
		})
	}

	table.Render()
}

func getShortModelNames(modelList []string) []string {
	var shortNames []string
	seen := make(map[string]bool)
	
	for _, model := range modelList {
		short := models.GetModelShortName(model)
		if !seen[short] {
			shortNames = append(shortNames, short)
			seen[short] = true
		}
	}
	
	return shortNames
}

func formatNumber(n int) string {
	if n == 0 {
		return "-"
	}
	return fmt.Sprintf("%d", n)
}

