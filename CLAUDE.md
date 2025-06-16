# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

claude-usage-go is a CLI tool that analyzes Claude API usage from local JSONL files (`~/.claude/projects/`) and calculates token usage and estimated costs. It's a Go implementation inspired by ccusage.

## Build and Development Commands

```bash
# Build the project
go build -o claude-usage-go

# Download dependencies
go mod tidy

# Run the tool
./claude-usage-go daily
./claude-usage-go monthly
./claude-usage-go session

# Common development commands
go fmt ./...  # Format code
go vet ./...  # Run static analysis
```

## Architecture

The codebase follows a clean architecture with clear separation of concerns:

1. **CLI Layer** (`/cmd/`): Cobra-based commands that handle user interaction
   - Each command (daily, monthly, session) follows the same pattern: parse options → get data → filter → aggregate → display
   - Global flags are defined in `root.go` and accessed via package-level variables

2. **Data Models** (`/internal/models/`): Core types and pricing data
   - `TokenUsage` tracks all token types (input, output, cache create/read)
   - Model pricing is hardcoded in `pricing.go` - update here when Claude pricing changes
   - Model IDs follow pattern: `claude-{version}-{model}-{date}` (e.g., `claude-opus-4-20250514`)

3. **JSONL Parser** (`/internal/parser/`): Handles file reading and parsing
   - JSONL format: Each line contains a message with nested structure
   - Key parsing logic: Looks for `type: "assistant"` with `message.role: "assistant"` and `message.usage` data
   - Buffer size is set to 10MB to handle large files

4. **Cost Calculator** (`/internal/calculator/`): Aggregates usage and calculates costs
   - All aggregation functions return sorted results
   - Cost calculation: `tokens / 1_000_000 * price_per_1M`
   - Supports daily, monthly, and session aggregation

5. **Display** (`/internal/display/`): Formats output as tables or JSON
   - Table display uses ANSI colors and Unicode box characters
   - Breakdown mode shows per-model costs within each time period

## Key Implementation Details

- **JSONL Structure**: Messages are nested inside a `message` field, not at top level
- **Model Naming**: Use `GetModelShortName()` for display (e.g., "Opus 4" instead of full ID)
- **Date Filtering**: Uses Go's time.Parse with "20060102" format (YYYYMMDD)
- **Model Filtering**: Case-insensitive comparison of model names
- **No External APIs**: All data comes from local JSONL files, no network calls

## Adding New Features

When adding new Claude models:
1. Add pricing to `internal/models/pricing.go` in `ModelPricing` map
2. Add short name mapping in `GetModelShortName()` function

When modifying JSONL parsing:
1. Check actual JSONL structure with: `grep '"type":"message"' ~/.claude/projects/*/**.jsonl | head -1 | python3 -m json.tool`
2. Update `JSONLEntry` and related structs in `internal/parser/jsonl.go`

## Current Limitations

- No test files exist - consider adding tests when making changes
- No automated linting or CI/CD setup
- Hardcoded pricing data (not fetched from API)
- Only reads from `~/.claude/projects/` directory