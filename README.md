# claude-usage-go

A Go implementation of a Claude API usage and pricing display tool, similar to [ccusage](https://github.com/ryoppippi/ccusage). This CLI tool analyzes Claude API usage from local JSONL files and calculates token usage and estimated costs.

## Features

- **Multiple Report Types**:
  - Daily reports: View token usage and costs aggregated by date
  - Monthly reports: See usage aggregated by month  
  - Session reports: Analyze usage grouped by conversation sessions

- **Comprehensive Token Tracking**:
  - Input tokens
  - Output tokens
  - Cache creation tokens
  - Cache read tokens
  - Total token counts
  - Estimated costs in USD

- **Model-Specific Analysis**:
  - Track usage across all Claude models including Claude 4 (Opus 4, Sonnet 4)
  - Per-model cost breakdown with the `--breakdown` flag
  - Model filtering options

- **Flexible Display Options**:
  - Colorful table-formatted output (default)
  - JSON output format for programmatic use
  - Date range filtering
  - Ascending/descending sorting

## Installation

```bash
go install github.com/t-ishitsuka/claude-usage-go@latest
```

Or clone and build from source:

```bash
git clone https://github.com/t-ishitsuka/claude-usage-go.git
cd claude-usage-go
go build
```

## Usage

The tool reads JSONL files from `~/.claude/projects/` directory.

### Basic Commands

```bash
# Show daily usage
./claude-usage-go daily

# Show monthly usage
./claude-usage-go monthly

# Show session usage
./claude-usage-go session
```

### Options

- `--since YYYYMMDD`: Start date filter
- `--until YYYYMMDD`: End date filter
- `--breakdown`: Show model-specific breakdown
- `--json`: Output as JSON
- `--asc`: Sort in ascending order (default is descending)
- `--models`: Filter by specific models (comma-separated)

### Examples

```bash
# Show daily usage for a specific date range
./claude-usage-go daily --since 20250615 --until 20250616

# Show monthly usage with model breakdown
./claude-usage-go monthly --breakdown

# Show session usage as JSON
./claude-usage-go session --json

# Filter by specific models
./claude-usage-go daily --models claude-opus-4-20250514,claude-3-5-sonnet-20241022

# Show daily usage in ascending order
./claude-usage-go daily --asc
```

## Output Format

The tool displays usage data in a formatted table with:
- Date/Month/Session ID
- Models used
- Token counts (Input, Output, Cache Create, Cache Read, Total)
- Estimated cost in USD

Example output:
```
│────────────│─────────────────────│───────│────────│──────────────│────────────│─────────│────────────│
│    DATE    │       MODELS        │ INPUT │ OUTPUT │ CACHE CREATE │ CACHE READ │  TOTAL  │ COST (USD) │
│────────────│─────────────────────│───────│────────│──────────────│────────────│─────────│────────────│
│ 2025-06-16 │ <synthetic>, Opus 4 │   434 │  16826 │       461951 │    3454156 │ 3933367 │ $15.1113   │
│────────────│─────────────────────│───────│────────│──────────────│────────────│─────────│────────────│
│   TOTAL    │                     │   434 │  16826 │       461951 │    3454156 │ 3933367 │ $15.1113   │
│────────────│─────────────────────│───────│────────│──────────────│────────────│─────────│────────────│
```

## Supported Models and Pricing

The tool includes pricing for:
- **Claude 4 Models**:
  - Opus 4: $15/$75 per 1M tokens (input/output)
  - Sonnet 4: $3/$15 per 1M tokens
- **Claude 3.5 Models**:
  - Sonnet 3.5: $3/$15 per 1M tokens
  - Haiku 3.5: $0.80/$4 per 1M tokens
- **Claude 3 Models**:
  - Opus 3: $15/$75 per 1M tokens
  - Sonnet 3: $3/$15 per 1M tokens
  - Haiku 3: $0.25/$1.25 per 1M tokens

Pricing is based on the official Claude API pricing structure.

## Requirements

- Go 1.21 or higher
- Access to Claude JSONL files in `~/.claude/projects/`

## Troubleshooting

### "Token too long" error
The tool handles large JSONL files with a 10MB buffer size. If you encounter issues with extremely large files, please open an issue.

### No data displayed
Ensure that:
1. Claude JSONL files exist in `~/.claude/projects/`
2. The JSONL files contain assistant messages with usage data
3. You're using the correct date range filters

## License

MIT