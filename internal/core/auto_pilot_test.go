package core

import (
	"testing"
)

func TestNewPilotConfig(t *testing.T) {
	cfg := NewPilotConfig()

	if cfg.DiscoverInterval != DefaultDiscoverInterval {
		t.Errorf("expected DiscoverInterval=%d, got=%d",
			DefaultDiscoverInterval, cfg.DiscoverInterval)
	}
	if cfg.MaxDiscoveryTasks != DefaultMaxDiscoveryTasks {
		t.Errorf("expected MaxDiscoveryTasks=%d, got=%d",
			DefaultMaxDiscoveryTasks, cfg.MaxDiscoveryTasks)
	}
	if cfg.Focus != "" {
		t.Errorf("expected empty Focus, got=%s", cfg.Focus)
	}
}

func TestCountPendingTasks(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []AutoTask
		expected int
	}{
		{"no tasks", nil, 0},
		{"all pending", []AutoTask{
			{ID: "1", Status: TaskStatusPending},
			{ID: "2", Status: TaskStatusPending},
		}, 2},
		{"mixed statuses", []AutoTask{
			{ID: "1", Status: TaskStatusPending},
			{ID: "2", Status: TaskStatusCompleted},
			{ID: "3", Status: TaskStatusSkipped},
			{ID: "4", Status: TaskStatusPending},
		}, 2},
		{"none pending", []AutoTask{
			{ID: "1", Status: TaskStatusCompleted},
			{ID: "2", Status: TaskStatusSkipped},
		}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prd := &AutoPRD{Tasks: tt.tasks}
			got := CountPendingTasks(prd)
			if got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestShouldRunDiscovery(t *testing.T) {
	tests := []struct {
		name              string
		tasks             []AutoTask
		currentIter       int
		lastDiscoveryIter int
		discoverInterval  int
		expected          bool
	}{
		{
			name:              "no tasks - should discover",
			tasks:             nil,
			currentIter:       1,
			lastDiscoveryIter: 0,
			discoverInterval:  5,
			expected:          true,
		},
		{
			name: "no pending tasks - should discover",
			tasks: []AutoTask{
				{ID: "1", Status: TaskStatusCompleted},
			},
			currentIter:       3,
			lastDiscoveryIter: 1,
			discoverInterval:  5,
			expected:          true,
		},
		{
			name: "first iteration, never discovered - should discover",
			tasks: []AutoTask{
				{ID: "1", Status: TaskStatusPending},
				{ID: "2", Status: TaskStatusPending},
				{ID: "3", Status: TaskStatusPending},
			},
			currentIter:       1,
			lastDiscoveryIter: 0,
			discoverInterval:  5,
			expected:          true,
		},
		{
			name: "interval reached - should discover",
			tasks: []AutoTask{
				{ID: "1", Status: TaskStatusPending},
				{ID: "2", Status: TaskStatusPending},
				{ID: "3", Status: TaskStatusPending},
			},
			currentIter:       6,
			lastDiscoveryIter: 1,
			discoverInterval:  5,
			expected:          true,
		},
		{
			name: "low pending count - should discover preemptively",
			tasks: []AutoTask{
				{ID: "1", Status: TaskStatusPending},
				{ID: "2", Status: TaskStatusCompleted},
			},
			currentIter:       3,
			lastDiscoveryIter: 2,
			discoverInterval:  5,
			expected:          true,
		},
		{
			name: "plenty of tasks, interval not reached - should not discover",
			tasks: []AutoTask{
				{ID: "1", Status: TaskStatusPending},
				{ID: "2", Status: TaskStatusPending},
				{ID: "3", Status: TaskStatusPending},
			},
			currentIter:       3,
			lastDiscoveryIter: 1,
			discoverInterval:  5,
			expected:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prd := &AutoPRD{Tasks: tt.tasks}
			got := ShouldRunDiscovery(prd, tt.currentIter, tt.lastDiscoveryIter, tt.discoverInterval)
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestInitPilotPRD(t *testing.T) {
	pilot := NewPilotConfig()
	pilot.Focus = "testing"

	config := AutoConfig{
		MaxIterations: 30,
		QualityChecks: []string{"go test ./..."},
		AITool:        "claude",
		Sandbox:       SandboxNone,
	}

	prd := InitPilotPRD("/tmp/myproject", config, pilot)

	if prd.Project.Name != "myproject" {
		t.Errorf("expected Name=myproject, got=%s", prd.Project.Name)
	}
	if !prd.Config.PilotMode {
		t.Error("expected PilotMode=true")
	}
	if prd.Config.PilotConfig == nil {
		t.Fatal("expected PilotConfig to be set")
	}
	if prd.Config.PilotConfig.Focus != "testing" {
		t.Errorf("expected Focus=testing, got=%s", prd.Config.PilotConfig.Focus)
	}
	if prd.Config.MaxIterations != 30 {
		t.Errorf("expected MaxIterations=30, got=%d", prd.Config.MaxIterations)
	}
	if prd.Config.AITool != "claude" {
		t.Errorf("expected AITool=claude, got=%s", prd.Config.AITool)
	}
	if prd.Config.DiscoveryPrompt == "" {
		t.Error("expected DiscoveryPrompt to be set")
	}
	if prd.Progress.Status != LoopStatusNotStarted {
		t.Errorf("expected Status=not_started, got=%s", prd.Progress.Status)
	}
	if len(prd.Tasks) != 0 {
		t.Errorf("expected 0 tasks, got=%d", len(prd.Tasks))
	}
}
