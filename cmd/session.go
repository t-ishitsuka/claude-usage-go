package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/claude-usage-go/internal/calculator"
	"github.com/yourusername/claude-usage-go/internal/display"
	"github.com/yourusername/claude-usage-go/internal/parser"
)

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Show session usage report",
	Long:  `Display token usage and costs aggregated by session.`,
	RunE:  runSession,
}

func init() {
	rootCmd.AddCommand(sessionCmd)
}

func runSession(cmd *cobra.Command, args []string) error {
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

	sessionUsage := calculator.AggregateBySession(messages)

	if opts.JSONOutput {
		return outputJSON(sessionUsage)
	}

	if opts.Breakdown {
		return display.ShowSessionWithBreakdown(sessionUsage, messages, opts.Ascending)
	}

	return display.ShowSession(sessionUsage, opts.Ascending)
}