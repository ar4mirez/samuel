package core

import (
	"path/filepath"
	"time"
)

// NewPilotConfig creates a PilotConfig with default values.
func NewPilotConfig() *PilotConfig {
	return &PilotConfig{
		DiscoverInterval:  DefaultDiscoverInterval,
		MaxDiscoveryTasks: DefaultMaxDiscoveryTasks,
	}
}

// ShouldRunDiscovery determines if a discovery iteration should run.
// Discovery triggers when:
//   - no pending tasks exist
//   - the discover interval has elapsed
//   - fewer than MinPendingTasksForDiscovery tasks remain (preemptive)
func ShouldRunDiscovery(prd *AutoPRD, currentIter, lastDiscoveryIter, discoverInterval int) bool {
	pending := CountPendingTasks(prd)

	if pending == 0 {
		return true
	}

	if lastDiscoveryIter == 0 {
		return true
	}

	if currentIter-lastDiscoveryIter >= discoverInterval {
		return true
	}

	if pending < MinPendingTasksForDiscovery {
		return true
	}

	return false
}

// CountPendingTasks returns the number of pending tasks in the PRD.
func CountPendingTasks(prd *AutoPRD) int {
	count := 0
	for _, t := range prd.Tasks {
		if t.Status == TaskStatusPending {
			count++
		}
	}
	return count
}

// InitPilotPRD creates a new AutoPRD configured for pilot mode.
func InitPilotPRD(projectDir string, config AutoConfig, pilot *PilotConfig) *AutoPRD {
	dirName := filepath.Base(projectDir)
	now := time.Now().UTC().Format(time.RFC3339)

	prd := &AutoPRD{
		Version: AutoSchemaVer,
		Project: AutoProject{
			Name:        dirName,
			Description: "Autonomous pilot mode - AI-discovered tasks",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Config: AutoConfig{
			MaxIterations:   config.MaxIterations,
			QualityChecks:   config.QualityChecks,
			AITool:          config.AITool,
			PromptFile:      filepath.Join(AutoDir, AutoPromptFile),
			Sandbox:         config.Sandbox,
			SandboxImage:    config.SandboxImage,
			SandboxTemplate: config.SandboxTemplate,
			PilotMode:       true,
			PilotConfig:     pilot,
			DiscoveryPrompt: filepath.Join(AutoDir, AutoDiscoveryPromptFile),
		},
		Tasks: []AutoTask{},
		Progress: AutoProgress{
			Status: LoopStatusNotStarted,
		},
	}

	return prd
}
