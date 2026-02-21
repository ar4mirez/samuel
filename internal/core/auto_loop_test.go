package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestNewLoopConfig_Defaults(t *testing.T) {
	dir := t.TempDir()
	prd := NewAutoPRD("test", "test project")
	prd.Config.AITool = "claude"
	prd.Config.Sandbox = SandboxNone

	prdPath := filepath.Join(dir, AutoDir, AutoPRDFile)
	if err := prd.Save(prdPath); err != nil {
		t.Fatalf("failed to save prd: %v", err)
	}

	cfg := NewLoopConfig(dir, prd)

	if cfg.ProjectDir != dir {
		t.Errorf("expected ProjectDir=%s, got=%s", dir, cfg.ProjectDir)
	}
	if cfg.AITool != "claude" {
		t.Errorf("expected AITool=claude, got=%s", cfg.AITool)
	}
	if cfg.MaxIterations != 50 {
		t.Errorf("expected MaxIterations=50, got=%d", cfg.MaxIterations)
	}
	if cfg.PauseSecs != 2 {
		t.Errorf("expected PauseSecs=2, got=%d", cfg.PauseSecs)
	}
	if cfg.MaxConsecFails != 3 {
		t.Errorf("expected MaxConsecFails=3, got=%d", cfg.MaxConsecFails)
	}
}

func TestNewLoopConfig_EnvOverrides(t *testing.T) {
	os.Setenv("PAUSE_SECONDS", "5")
	os.Setenv("MAX_CONSECUTIVE_FAILURES", "7")
	defer os.Unsetenv("PAUSE_SECONDS")
	defer os.Unsetenv("MAX_CONSECUTIVE_FAILURES")

	dir := t.TempDir()
	prd := NewAutoPRD("test", "test project")

	cfg := NewLoopConfig(dir, prd)

	if cfg.PauseSecs != 5 {
		t.Errorf("expected PauseSecs=5, got=%d", cfg.PauseSecs)
	}
	if cfg.MaxConsecFails != 7 {
		t.Errorf("expected MaxConsecFails=7, got=%d", cfg.MaxConsecFails)
	}
}

func TestRunAutoLoop_AllTasksCompleted(t *testing.T) {
	dir := t.TempDir()
	prd := NewAutoPRD("test", "test project")
	// No tasks means GetNextTask returns nil -> immediate exit
	prdPath := filepath.Join(dir, AutoDir, AutoPRDFile)
	if err := prd.Save(prdPath); err != nil {
		t.Fatalf("failed to save prd: %v", err)
	}

	cfg := LoopConfig{
		ProjectDir:     dir,
		PRDPath:        prdPath,
		MaxIterations:  10,
		MaxConsecFails: 3,
	}

	err := RunAutoLoop(cfg)
	if err != nil {
		t.Errorf("expected no error when all tasks done, got: %v", err)
	}
}

func TestRunAutoLoop_BadPRDPath(t *testing.T) {
	cfg := LoopConfig{
		ProjectDir:     t.TempDir(),
		PRDPath:        "/nonexistent/prd.json",
		MaxIterations:  1,
		MaxConsecFails: 3,
	}

	err := RunAutoLoop(cfg)
	if err == nil {
		t.Error("expected error for non-existent PRD path")
	}
}

func TestNotifyCallbacks(t *testing.T) {
	startCalled := false
	endCalled := false

	notifyIterStart(func(iter int, iterType string) {
		startCalled = true
		if iter != 1 {
			t.Errorf("expected iter=1, got=%d", iter)
		}
		if iterType != IterationTypeImplementation {
			t.Errorf("expected type=implementation, got=%s", iterType)
		}
	}, 1, IterationTypeImplementation)

	notifyIterEnd(func(iter int, err error) {
		endCalled = true
	}, 1, nil)

	if !startCalled {
		t.Error("OnIterStart callback was not called")
	}
	if !endCalled {
		t.Error("OnIterEnd callback was not called")
	}
}

func TestNotifyCallbacks_NilSafe(t *testing.T) {
	// Should not panic with nil callbacks
	notifyIterStart(nil, 1, IterationTypeImplementation)
	notifyIterEnd(nil, 1, nil)
}

func TestInvokeAgent_RejectsInvalidAITool(t *testing.T) {
	tests := []struct {
		name   string
		aiTool string
	}{
		{"arbitrary binary", "/bin/sh"},
		{"path traversal", "../../../malicious"},
		{"empty string", ""},
		{"shell injection", "claude; rm -rf /"},
		{"unknown tool", "unsupported-tool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := LoopConfig{
				ProjectDir: t.TempDir(),
				AITool:     tt.aiTool,
			}
			err := InvokeAgent(cfg)
			if err == nil {
				t.Errorf("expected error for invalid AI tool %q", tt.aiTool)
			}
		})
	}
}

func TestInvokeAgent_AcceptsValidTools(t *testing.T) {
	for _, tool := range GetSupportedAITools() {
		t.Run(tool, func(t *testing.T) {
			cfg := LoopConfig{
				ProjectDir: t.TempDir(),
				AITool:     tool,
				PromptPath: "/nonexistent/prompt.md",
			}
			err := InvokeAgent(cfg)
			// Should fail with prompt/exec error, NOT with invalid tool error
			if err != nil && err.Error() == fmt.Sprintf(
				"refused to invoke invalid AI tool %q (allowed: %v)",
				tool, GetSupportedAITools()) {
				t.Errorf("valid tool %q was rejected", tool)
			}
		})
	}
}
