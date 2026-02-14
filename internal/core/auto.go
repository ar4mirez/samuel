package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Auto directory and file constants
const (
	AutoDir          = ".claude/auto"
	AutoPRDFile      = "prd.json"
	AutoProgressFile = "progress.md"
	AutoPromptFile   = "prompt.md"
	AutoScriptFile   = "auto.sh"
	AutoSchemaVer    = "1.0"
)

// Task status constants
const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusSkipped    = "skipped"
	TaskStatusBlocked    = "blocked"
)

// Task priority constants
const (
	TaskPriorityCritical = "critical"
	TaskPriorityHigh     = "high"
	TaskPriorityMedium   = "medium"
	TaskPriorityLow      = "low"
)

// Task complexity constants
const (
	TaskComplexitySimple  = "simple"
	TaskComplexityMedium  = "medium"
	TaskComplexityComplex = "complex"
)

// Loop status constants
const (
	LoopStatusNotStarted = "not_started"
	LoopStatusRunning    = "running"
	LoopStatusPaused     = "paused"
	LoopStatusCompleted  = "completed"
	LoopStatusFailed     = "failed"
)

// AutoPRD represents the machine-readable task state for autonomous loops
type AutoPRD struct {
	Version  string       `json:"version"`
	Project  AutoProject  `json:"project"`
	Config   AutoConfig   `json:"config"`
	Tasks    []AutoTask   `json:"tasks"`
	Progress AutoProgress `json:"progress"`
}

// AutoProject holds project metadata
type AutoProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SourcePRD   string `json:"source_prd,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// AutoConfig holds loop configuration
type AutoConfig struct {
	MaxIterations int      `json:"max_iterations"`
	QualityChecks []string `json:"quality_checks"`
	AITool        string   `json:"ai_tool"`
	PromptFile    string   `json:"ai_prompt_file"`
	Sandbox       string   `json:"sandbox"`
	SandboxImage  string   `json:"sandbox_image,omitempty"`
}

// AutoTask represents a single task in the autonomous loop
type AutoTask struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description,omitempty"`
	Status        string   `json:"status"`
	Priority      string   `json:"priority,omitempty"`
	Complexity    string   `json:"complexity,omitempty"`
	ParentID      string   `json:"parent_id,omitempty"`
	DependsOn     []string `json:"depends_on,omitempty"`
	FilesToCreate []string `json:"files_to_create,omitempty"`
	FilesToModify []string `json:"files_to_modify,omitempty"`
	Guardrails    []string `json:"guardrails,omitempty"`
	CompletedAt   string   `json:"completed_at,omitempty"`
	CommitSHA     string   `json:"commit_sha,omitempty"`
	Iteration     int      `json:"iteration,omitempty"`
}

// AutoProgress holds summary progress data
type AutoProgress struct {
	TotalTasks         int    `json:"total_tasks"`
	CompletedTasks     int    `json:"completed_tasks"`
	CurrentIteration   int    `json:"current_iteration"`
	TotalIterationsRun int    `json:"total_iterations_run"`
	LastIterationAt    string `json:"last_iteration_at,omitempty"`
	Status             string `json:"status"`
}

// NewAutoPRD creates a new AutoPRD with defaults
func NewAutoPRD(name, description string) *AutoPRD {
	now := time.Now().UTC().Format(time.RFC3339)
	return &AutoPRD{
		Version: AutoSchemaVer,
		Project: AutoProject{
			Name:        name,
			Description: description,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Config: AutoConfig{
			MaxIterations: 50,
			QualityChecks: []string{"go test ./...", "go vet ./...", "go build ./..."},
			AITool:        "claude",
			PromptFile:    filepath.Join(AutoDir, AutoPromptFile),
			Sandbox:       SandboxNone,
		},
		Tasks: []AutoTask{},
		Progress: AutoProgress{
			Status: LoopStatusNotStarted,
		},
	}
}

// LoadAutoPRD loads a prd.json file from disk
func LoadAutoPRD(path string) (*AutoPRD, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read prd.json: %w", err)
	}

	var prd AutoPRD
	if err := json.Unmarshal(data, &prd); err != nil {
		return nil, fmt.Errorf("failed to parse prd.json: %w", err)
	}

	return &prd, nil
}

// Save writes the AutoPRD to disk using write-to-temp-then-rename for safety
func (p *AutoPRD) Save(path string) error {
	p.Project.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	p.RecalculateProgress()

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal prd.json: %w", err)
	}
	data = append(data, '\n')

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tmpFile, path); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// priorityRank returns a numeric rank for sorting (lower = higher priority)
func priorityRank(priority string) int {
	switch priority {
	case TaskPriorityCritical:
		return 0
	case TaskPriorityHigh:
		return 1
	case TaskPriorityMedium:
		return 2
	case TaskPriorityLow:
		return 3
	default:
		return 2
	}
}

// GetNextTask returns the highest-priority available pending task
func (p *AutoPRD) GetNextTask() *AutoTask {
	available := p.getAvailableTasks()
	if len(available) == 0 {
		return nil
	}

	sort.Slice(available, func(i, j int) bool {
		pi := priorityRank(available[i].Priority)
		pj := priorityRank(available[j].Priority)
		if pi != pj {
			return pi < pj
		}
		return available[i].ID < available[j].ID
	})

	return available[0]
}

// getAvailableTasks returns pending tasks whose dependencies are all completed
func (p *AutoPRD) getAvailableTasks() []*AutoTask {
	completed := make(map[string]bool)
	for i := range p.Tasks {
		if p.Tasks[i].Status == TaskStatusCompleted || p.Tasks[i].Status == TaskStatusSkipped {
			completed[p.Tasks[i].ID] = true
		}
	}

	var available []*AutoTask
	for i := range p.Tasks {
		if p.Tasks[i].Status != TaskStatusPending {
			continue
		}
		if allDependenciesMet(p.Tasks[i].DependsOn, completed) {
			available = append(available, &p.Tasks[i])
		}
	}
	return available
}

func allDependenciesMet(deps []string, completed map[string]bool) bool {
	for _, dep := range deps {
		if !completed[dep] {
			return false
		}
	}
	return true
}

// findTask returns a pointer to the task with the given ID
func (p *AutoPRD) findTask(id string) *AutoTask {
	for i := range p.Tasks {
		if p.Tasks[i].ID == id {
			return &p.Tasks[i]
		}
	}
	return nil
}

// CompleteTask marks a task as completed with commit info
func (p *AutoPRD) CompleteTask(id, commitSHA string, iteration int) error {
	task := p.findTask(id)
	if task == nil {
		return fmt.Errorf("task not found: %s", id)
	}

	task.Status = TaskStatusCompleted
	task.CompletedAt = time.Now().UTC().Format(time.RFC3339)
	task.CommitSHA = commitSHA
	task.Iteration = iteration
	return nil
}

// SkipTask marks a task as skipped
func (p *AutoPRD) SkipTask(id string) error {
	task := p.findTask(id)
	if task == nil {
		return fmt.Errorf("task not found: %s", id)
	}
	task.Status = TaskStatusSkipped
	return nil
}

// ResetTask resets a task to pending
func (p *AutoPRD) ResetTask(id string) error {
	task := p.findTask(id)
	if task == nil {
		return fmt.Errorf("task not found: %s", id)
	}

	task.Status = TaskStatusPending
	task.CompletedAt = ""
	task.CommitSHA = ""
	task.Iteration = 0
	return nil
}

// AddTask appends a new task to the task list
func (p *AutoPRD) AddTask(task AutoTask) error {
	if task.ID == "" {
		return fmt.Errorf("task ID is required")
	}
	if existing := p.findTask(task.ID); existing != nil {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}
	if task.Status == "" {
		task.Status = TaskStatusPending
	}
	p.Tasks = append(p.Tasks, task)
	return nil
}

// RecalculateProgress updates the progress summary from task data
func (p *AutoPRD) RecalculateProgress() {
	total := 0
	completed := 0
	for _, t := range p.Tasks {
		total++
		if t.Status == TaskStatusCompleted {
			completed++
		}
	}
	p.Progress.TotalTasks = total
	p.Progress.CompletedTasks = completed

	if total > 0 && completed == total {
		p.Progress.Status = LoopStatusCompleted
	}
}

// ValidateAutoPRD checks prd.json for structural issues
func ValidateAutoPRD(prd *AutoPRD) []string {
	var errors []string

	if prd.Version == "" {
		errors = append(errors, "version is required")
	}
	if prd.Project.Name == "" {
		errors = append(errors, "project.name is required")
	}

	errors = append(errors, validateTasks(prd.Tasks)...)
	return errors
}

// validateTasks checks task-level constraints
func validateTasks(tasks []AutoTask) []string {
	var errors []string
	ids := make(map[string]bool)

	for _, t := range tasks {
		if t.ID == "" {
			errors = append(errors, "task missing ID")
			continue
		}
		if ids[t.ID] {
			errors = append(errors, fmt.Sprintf("duplicate task ID: %s", t.ID))
		}
		ids[t.ID] = true

		if t.Title == "" {
			errors = append(errors, fmt.Sprintf("task %s missing title", t.ID))
		}
		if !isValidStatus(t.Status) {
			errors = append(errors, fmt.Sprintf("task %s has invalid status: %s", t.ID, t.Status))
		}
	}

	// Validate dependency references
	for _, t := range tasks {
		for _, dep := range t.DependsOn {
			if !ids[dep] {
				errors = append(errors, fmt.Sprintf("task %s depends on unknown task: %s", t.ID, dep))
			}
		}
	}

	return errors
}

func isValidStatus(status string) bool {
	switch status {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted,
		TaskStatusSkipped, TaskStatusBlocked:
		return true
	default:
		return false
	}
}

// GetAutoPRDPath returns the full path to prd.json in a project directory
func GetAutoPRDPath(projectDir string) string {
	return filepath.Join(projectDir, AutoDir, AutoPRDFile)
}

// GetAutoDir returns the full path to the .claude/auto directory
func GetAutoDir(projectDir string) string {
	return filepath.Join(projectDir, AutoDir)
}
