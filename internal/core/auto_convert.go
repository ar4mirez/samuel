package core

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// taskLineRegex parses task lines from the generate-tasks skill output format:
//
//	"- [ ] 1.0 Task title [~3,000 tokens - Medium]"
//
// Groups: (1) indentation, (2) checkbox, (3) task ID, (4) title, (5) complexity
var taskLineRegex = regexp.MustCompile(
	`^(\s*)- \[([ xX])\]\s*(\d+\.\d+)\s+(.+?)(?:\s*\[~[\d,]+\s+tokens?\s*-\s*(\w+)\])?\s*$`,
)

// prdTitleRegex extracts the title from a PRD markdown H1 heading
var prdTitleRegex = regexp.MustCompile(`^#\s+(.+)$`)

// ConvertMarkdownToPRD converts a PRD markdown file and optional task list
// into a structured AutoPRD. If tasksPath is empty, only project metadata
// is extracted from the PRD.
func ConvertMarkdownToPRD(prdPath, tasksPath string) (*AutoPRD, error) {
	prdContent, err := os.ReadFile(prdPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PRD: %w", err)
	}

	name, description := extractPRDMetadata(string(prdContent))
	prd := NewAutoPRD(name, description)
	prd.Project.SourcePRD = prdPath
	prd.Project.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	if tasksPath != "" {
		tasksContent, err := os.ReadFile(tasksPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read tasks file: %w", err)
		}
		tasks, err := ParseTaskMarkdown(string(tasksContent))
		if err != nil {
			return nil, fmt.Errorf("failed to parse tasks: %w", err)
		}
		prd.Tasks = tasks
	}

	prd.RecalculateProgress()
	return prd, nil
}

// extractPRDMetadata extracts name and description from PRD markdown content
func extractPRDMetadata(content string) (name, description string) {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if matches := prdTitleRegex.FindStringSubmatch(trimmed); matches != nil {
			name = slugify(matches[1])
			description = matches[1]
			break
		}
	}

	if name == "" {
		name = "unnamed-project"
		description = "Converted from PRD"
	}
	return name, description
}

// slugify converts a title to a URL-friendly slug
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		if r == ' ' || r == '-' || r == '_' {
			return '-'
		}
		return -1
	}, s)
	// Collapse multiple hyphens
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

// ParseTaskMarkdown parses task markdown content into AutoTask structs.
// It handles the format from AICoF's generate-tasks skill:
//
//   - [ ] 1.0 Parent Task Title
//   - [ ] 1.1 Sub-task description [~2,000 tokens - Simple]
func ParseTaskMarkdown(content string) ([]AutoTask, error) {
	lines := strings.Split(content, "\n")
	var tasks []AutoTask
	var currentParentID string

	for _, line := range lines {
		task, err := parseTaskLine(line)
		if err != nil {
			continue // Skip non-task lines leniently
		}
		if task == nil {
			continue
		}

		// Determine parent-child relationship from indentation
		if isChildTask(line) {
			task.ParentID = currentParentID
			if currentParentID != "" {
				task.DependsOn = []string{currentParentID}
			}
		} else {
			currentParentID = task.ID
		}

		tasks = append(tasks, *task)
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("no valid tasks found in markdown")
	}
	return tasks, nil
}

// parseTaskLine parses a single markdown task line into an AutoTask.
// Returns nil, nil for non-task lines (not an error, just not a match).
func parseTaskLine(line string) (*AutoTask, error) {
	matches := taskLineRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, nil
	}

	checkbox := matches[2]
	taskID := matches[3]
	title := strings.TrimSpace(matches[4])
	complexity := strings.ToLower(strings.TrimSpace(matches[5]))

	status := TaskStatusPending
	if checkbox == "x" || checkbox == "X" {
		status = TaskStatusCompleted
	}

	if !isValidComplexity(complexity) {
		complexity = TaskComplexityMedium
	}

	return &AutoTask{
		ID:         taskID,
		Title:      title,
		Status:     status,
		Complexity: complexity,
		Priority:   TaskPriorityMedium,
	}, nil
}

// isChildTask checks if a task line is indented (child task)
func isChildTask(line string) bool {
	return len(line) > 0 && (line[0] == ' ' || line[0] == '\t')
}

func isValidComplexity(c string) bool {
	switch c {
	case TaskComplexitySimple, TaskComplexityMedium, TaskComplexityComplex:
		return true
	default:
		return false
	}
}

// FindTasksFile attempts to locate the corresponding tasks file for a PRD.
// Convention: PRD at ".claude/tasks/0001-prd-feature.md"
// maps to tasks at ".claude/tasks/tasks-0001-prd-feature.md"
func FindTasksFile(prdPath string) string {
	dir := filepath.Dir(prdPath)
	base := filepath.Base(prdPath)
	tasksFile := filepath.Join(dir, "tasks-"+base)

	if _, err := os.Stat(tasksFile); err == nil {
		return tasksFile
	}
	return ""
}
