package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/t-ishitsuka/claude-usage-go/internal/calculator"
	"github.com/t-ishitsuka/claude-usage-go/internal/display"
	"github.com/t-ishitsuka/claude-usage-go/internal/models"
	"github.com/t-ishitsuka/claude-usage-go/internal/parser"
)

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Show daily usage report",
	Long:  `Display token usage and costs aggregated by day.`,
	RunE:  runDaily,
}

func init() {
	rootCmd.AddCommand(dailyCmd)
}

func runDaily(cmd *cobra.Command, args []string) error {
	opts, err := parseOptions()
	if err != nil {
		return err
	}

	projectsDir := parser.GetClaudeProjectsDir()
	messages, err := parser.ParseJSONLFiles(projectsDir)
	if err != nil {
		return fmt.Errorf("error parsing JSONL files: %w", err)
	}

	messages = parser.FilterByDateRange(messages, opts.Since, opts.Until)
	messages = parser.FilterByModels(messages, opts.Models)

	dailyUsage := calculator.AggregateDaily(messages)

	if opts.JSONOutput {
		return outputJSON(dailyUsage)
	}

	if opts.Breakdown {
		return display.ShowDailyWithBreakdown(dailyUsage, messages, opts.Ascending)
	}

	return display.ShowDaily(dailyUsage, opts.Ascending)
}

func parseOptions() (*models.ReportOptions, error) {
	opts := &models.ReportOptions{
		Breakdown:  breakdown,
		JSONOutput: jsonOutput,
		Ascending:  ascending,
		Models:     modelFilter,
	}

	if since != "" {
		t, err := time.Parse("20060102", since)
		if err != nil {
			return nil, fmt.Errorf("invalid since date format: %w", err)
		}
		opts.Since = &t
	}

	if until != "" {
		t, err := time.Parse("20060102", until)
		if err != nil {
			return nil, fmt.Errorf("invalid until date format: %w", err)
		}
		opts.Until = &t
	}

	return opts, nil
}

func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
