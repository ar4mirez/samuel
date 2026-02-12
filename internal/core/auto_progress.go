package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Progress entry type constants
const (
	ProgressStarted      = "STARTED"
	ProgressCompleted    = "COMPLETED"
	ProgressError        = "ERROR"
	ProgressLearning     = "LEARNING"
	ProgressQualityCheck = "QUALITY_CHECK"
	ProgressCommit       = "COMMIT"
)

// ProgressEntry represents a single line in progress.txt
type ProgressEntry struct {
	Iteration int
	TaskID    string
	Type      string
	Message   string
}

// FormatProgressEntry formats a progress entry for the append-only log
func FormatProgressEntry(entry ProgressEntry) string {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	parts := []string{fmt.Sprintf("[%s]", timestamp)}

	if entry.Iteration > 0 {
		parts = append(parts, fmt.Sprintf("[iteration:%d]", entry.Iteration))
	}
	if entry.TaskID != "" {
		parts = append(parts, fmt.Sprintf("[task:%s]", entry.TaskID))
	}

	parts = append(parts, fmt.Sprintf("%s: %s", entry.Type, entry.Message))
	return strings.Join(parts, " ")
}

// AppendProgress appends a formatted entry to the progress file
func AppendProgress(path string, entry ProgressEntry) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open progress file: %w", err)
	}
	defer f.Close()

	line := FormatProgressEntry(entry) + "\n"
	if _, err := f.WriteString(line); err != nil {
		return fmt.Errorf("failed to write progress entry: %w", err)
	}
	return nil
}

// ReadProgressTail reads the last N lines from the progress file
func ReadProgressTail(path string, lines int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open progress file: %w", err)
	}
	defer f.Close()

	var allLines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read progress file: %w", err)
	}

	if lines <= 0 || lines >= len(allLines) {
		return allLines, nil
	}
	return allLines[len(allLines)-lines:], nil
}
