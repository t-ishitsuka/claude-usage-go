package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/t-ishitsuka/claude-usage-go/internal/models"
)

type JSONLEntry struct {
	SessionID string    `json:"sessionId"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Message   *Message  `json:"message,omitempty"`
}

type Message struct {
	Role  string `json:"role"`
	Model string `json:"model"`
	Usage *Usage `json:"usage,omitempty"`
}

type Usage struct {
	InputTokens       int `json:"input_tokens"`
	OutputTokens      int `json:"output_tokens"`
	CacheCreateTokens int `json:"cache_creation_input_tokens"`
	CacheReadTokens   int `json:"cache_read_input_tokens"`
}

func ParseJSONLFiles(directory string) ([]models.Message, error) {
	var messages []models.Message

	// Walk the directory tree to find all .jsonl files
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && filepath.Ext(path) == ".jsonl" {
			msgs, err := parseJSONLFile(path)
			if err != nil {
				return fmt.Errorf("error parsing %s: %w", path, err)
			}
			messages = append(messages, msgs...)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return messages, nil
}

func parseJSONLFile(filename string) ([]models.Message, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var messages []models.Message
	scanner := bufio.NewScanner(file)
	
	// Increase buffer size to handle large lines (10MB)
	const maxCapacity = 10 * 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		var entry JSONLEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		if entry.Type == "assistant" && entry.Message != nil && entry.Message.Role == "assistant" && entry.Message.Usage != nil {
			msg := models.Message{
				SessionID: entry.SessionID,
				Timestamp: entry.Timestamp,
				Model:     entry.Message.Model,
				TokenUsage: models.TokenUsage{
					InputTokens:       entry.Message.Usage.InputTokens,
					OutputTokens:      entry.Message.Usage.OutputTokens,
					CacheCreateTokens: entry.Message.Usage.CacheCreateTokens,
					CacheReadTokens:   entry.Message.Usage.CacheReadTokens,
				},
			}
			messages = append(messages, msg)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func GetClaudeProjectsDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".claude", "projects")
}

func FilterByDateRange(messages []models.Message, since, until *time.Time) []models.Message {
	var filtered []models.Message
	for _, msg := range messages {
		if since != nil && msg.Timestamp.Before(*since) {
			continue
		}
		if until != nil && msg.Timestamp.After(until.Add(24*time.Hour)) {
			continue
		}
		filtered = append(filtered, msg)
	}
	return filtered
}

func FilterByModels(messages []models.Message, modelList []string) []models.Message {
	if len(modelList) == 0 {
		return messages
	}

	modelSet := make(map[string]bool)
	for _, m := range modelList {
		modelSet[strings.ToLower(m)] = true
	}

	var filtered []models.Message
	for _, msg := range messages {
		if modelSet[strings.ToLower(msg.Model)] {
			filtered = append(filtered, msg)
		}
	}
	return filtered
}