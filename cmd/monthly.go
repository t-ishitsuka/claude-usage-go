package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/t-ishitsuka/claude-usage-go/internal/calculator"
	"github.com/t-ishitsuka/claude-usage-go/internal/display"
	"github.com/t-ishitsuka/claude-usage-go/internal/parser"
)

var monthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Show monthly usage report",
	Long:  `Display token usage and costs aggregated by month.`,
	RunE:  runMonthly,
}

func init() {
	rootCmd.AddCommand(monthlyCmd)
}

func runMonthly(cmd *cobra.Command, args []string) error {
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

	monthlyUsage := calculator.AggregateMonthly(messages)

	if opts.JSONOutput {
		return outputJSON(monthlyUsage)
	}

	if opts.Breakdown {
		return display.ShowMonthlyWithBreakdown(monthlyUsage, messages, opts.Ascending)
	}

	return display.ShowMonthly(monthlyUsage, opts.Ascending)
}
