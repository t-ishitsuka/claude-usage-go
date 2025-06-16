package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	since       string
	until       string
	breakdown   bool
	jsonOutput  bool
	ascending   bool
	modelFilter []string
)

var rootCmd = &cobra.Command{
	Use:   "claude-usage-go",
	Short: "Analyze Claude API usage and costs from local JSONL files",
	Long: `claude-usage-go is a CLI tool that analyzes Claude API usage from local JSONL files.
It calculates token usage and estimated costs based on current Claude API pricing.

Similar to ccusage, this tool helps you understand your Claude usage patterns and costs.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&since, "since", "", "Start date (YYYYMMDD format)")
	rootCmd.PersistentFlags().StringVar(&until, "until", "", "End date (YYYYMMDD format)")
	rootCmd.PersistentFlags().BoolVar(&breakdown, "breakdown", false, "Show model breakdown")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&ascending, "asc", false, "Sort in ascending order")
	rootCmd.PersistentFlags().StringSliceVar(&modelFilter, "models", []string{}, "Filter by models")
}