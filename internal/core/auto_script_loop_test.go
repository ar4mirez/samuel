package core

import (
	"strings"
	"testing"
)

func TestGenerateAutoScript_ConfigContainsCircuitBreaker(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "MAX_CONSECUTIVE_FAILURES") {
		t.Error("expected MAX_CONSECUTIVE_FAILURES in config section")
	}
	if !strings.Contains(script, `MAX_CONSECUTIVE_FAILURES:-3`) {
		t.Error("expected default value of 3 for MAX_CONSECUTIVE_FAILURES")
	}
}

func TestGenerateAutoScript_MainLoopCircuitBreaker(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "consecutive_failures=0") {
		t.Error("expected consecutive_failures counter initialization")
	}
	if !strings.Contains(script, "consecutive_failures=$((consecutive_failures + 1))") {
		t.Error("expected consecutive_failures increment on error")
	}
	if !strings.Contains(script, `"$consecutive_failures" -ge "$MAX_CONSECUTIVE_FAILURES"`) {
		t.Error("expected threshold check against MAX_CONSECUTIVE_FAILURES")
	}
	if !strings.Contains(script, "FATAL:") {
		t.Error("expected FATAL message when circuit breaker triggers")
	}
}

func TestGenerateAutoScript_MainLoopResetsOnSuccess(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	// The success branch should reset the counter
	mainLoop := script[strings.Index(script, "# --- Main Loop ---"):]
	runIdx := strings.Index(mainLoop, "if run_ai_tool; then")
	if runIdx < 0 {
		t.Fatal("expected 'if run_ai_tool; then' in main loop")
	}

	// After "if run_ai_tool; then" the next line should reset counter
	afterRun := mainLoop[runIdx:]
	resetIdx := strings.Index(afterRun, "consecutive_failures=0")
	elseIdx := strings.Index(afterRun, "else")
	if resetIdx < 0 || elseIdx < 0 {
		t.Fatal("expected both reset and else branches")
	}
	if resetIdx > elseIdx {
		t.Error("consecutive_failures=0 reset should be in the success branch (before else)")
	}
}

func TestGenerateAutoScript_AuthCheckSection(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "# --- Auth Check:") {
		t.Error("expected auth check section header")
	}
	if !strings.Contains(script, "check_ai_auth") {
		t.Error("expected check_ai_auth function")
	}
	if !strings.Contains(script, "ANTHROPIC_API_KEY") {
		t.Error("expected ANTHROPIC_API_KEY check for claude")
	}
	if !strings.Contains(script, "OPENAI_API_KEY") {
		t.Error("expected OPENAI_API_KEY check for codex")
	}
	if !strings.Contains(script, "AMP_API_KEY") {
		t.Error("expected AMP_API_KEY check for amp")
	}
}

func TestGenerateAutoScript_AuthCheckInDockerSandbox(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
		Sandbox:       SandboxDockerSandbox,
	}

	script := GenerateAutoScript(config)

	// Auth check should be present even when setup is skipped
	if !strings.Contains(script, "# --- Auth Check:") {
		t.Error("docker-sandbox mode should include auth check section")
	}
	if !strings.Contains(script, "check_ai_auth") {
		t.Error("docker-sandbox mode should include check_ai_auth function")
	}
	// But setup should still be absent
	if strings.Contains(script, "# --- Setup: ensure AI tool is available ---") {
		t.Error("docker-sandbox mode should NOT include setup section")
	}
}
