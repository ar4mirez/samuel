package core

import (
	"strings"
	"testing"
)

func TestGenerateAutoScript_ContainsSetupSection(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "# --- Setup: ensure AI tool is available ---") {
		t.Error("expected setup section header in generated script")
	}
	if !strings.Contains(script, "setup_ai_tool") {
		t.Error("expected setup_ai_tool function in generated script")
	}
	if !strings.Contains(script, "install_nodejs") {
		t.Error("expected install_nodejs function in generated script")
	}
}

func TestGenerateAutoScript_SetupInstallsClaude(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 10,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "claude.ai/install.sh") {
		t.Error("expected native installer URL in setup section")
	}
	if !strings.Contains(script, "native_install_claude") {
		t.Error("expected native_install_claude function in setup section")
	}
}

func TestGenerateAutoScript_SetupInstallsCodex(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 10,
		AITool:        "codex",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "@openai/codex") {
		t.Error("expected codex npm package in setup section")
	}
}

func TestGenerateAutoScript_SetupInstallsAmp(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 10,
		AITool:        "amp",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "install.amp.dev") {
		t.Error("expected amp install URL in setup section")
	}
}

func TestGenerateAutoScript_SetupVerifiesInstall(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 10,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "Failed to install") {
		t.Error("expected post-install verification in setup section")
	}
}

func TestGenerateAutoScript_SetupHandlesNonRootUser(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 10,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	// npm_install_tool should handle non-root installs with --prefix
	if !strings.Contains(script, "npm_install_tool") {
		t.Error("expected npm_install_tool helper in setup section")
	}
	if !strings.Contains(script, "--prefix") {
		t.Error("expected --prefix flag for non-root npm installs")
	}
	if !strings.Contains(script, `id -u`) {
		t.Error("expected root user detection via id -u")
	}
	// ensure_writable_home handles HOME fallback centrally
	if !strings.Contains(script, "ensure_writable_home") {
		t.Error("expected ensure_writable_home function in setup section")
	}
}

func TestGenerateAutoScript_SetupSkipsWhenToolExists(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 10,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	// Early return when tool is already in PATH
	if !strings.Contains(script, "found in PATH") {
		t.Error("expected early return message when tool already exists")
	}
}

func TestGenerateAutoScript_SkipsSetupForDockerSandbox(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
		Sandbox:       SandboxDockerSandbox,
	}

	script := GenerateAutoScript(config)

	if strings.Contains(script, "# --- Setup: ensure AI tool is available ---") {
		t.Error("docker-sandbox mode should NOT include setup section")
	}
	if strings.Contains(script, "setup_ai_tool") {
		t.Error("docker-sandbox mode should NOT include setup_ai_tool")
	}
	// But should still have the main loop and helpers
	if !strings.Contains(script, "# --- Main Loop ---") {
		t.Error("docker-sandbox mode should still include main loop")
	}
	if !strings.Contains(script, "# --- Helpers ---") {
		t.Error("docker-sandbox mode should still include helpers")
	}
}

func TestGenerateAutoScript_IncludesSetupForDockerMode(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
		Sandbox:       SandboxDocker,
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "# --- Setup: ensure AI tool is available ---") {
		t.Error("docker mode SHOULD include setup section")
	}
}

func TestGenerateAutoScript_IncludesSetupForNoneMode(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
		Sandbox:       SandboxNone,
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "# --- Setup: ensure AI tool is available ---") {
		t.Error("none mode SHOULD include setup section")
	}
}

func TestGenerateAutoScript_IncludesSetupForEmptySandbox(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
		Sandbox:       "",
	}

	script := GenerateAutoScript(config)

	if !strings.Contains(script, "# --- Setup: ensure AI tool is available ---") {
		t.Error("empty sandbox mode SHOULD include setup section (backward compat)")
	}
}

func TestGenerateAutoScript_SectionOrder(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	validationIdx := strings.Index(script, "# --- Validation ---")
	setupIdx := strings.Index(script, "# --- Setup:")
	authIdx := strings.Index(script, "# --- Auth Check:")
	helpersIdx := strings.Index(script, "# --- Helpers ---")
	mainLoopIdx := strings.Index(script, "# --- Main Loop ---")

	if validationIdx < 0 || setupIdx < 0 || authIdx < 0 || helpersIdx < 0 || mainLoopIdx < 0 {
		t.Fatal("missing expected section headers in generated script")
	}

	if validationIdx >= setupIdx {
		t.Error("validation section should come before setup section")
	}
	if setupIdx >= authIdx {
		t.Error("setup section should come before auth check section")
	}
	if authIdx >= helpersIdx {
		t.Error("auth check section should come before helpers section")
	}
	if helpersIdx >= mainLoopIdx {
		t.Error("helpers section should come before main loop section")
	}
}

func TestGenerateAutoScript_SectionOrderDockerSandbox(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 50,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
		Sandbox:       SandboxDockerSandbox,
	}

	script := GenerateAutoScript(config)

	// Should go directly from validation to auth check to helpers (no setup)
	validationIdx := strings.Index(script, "# --- Validation ---")
	authIdx := strings.Index(script, "# --- Auth Check:")
	helpersIdx := strings.Index(script, "# --- Helpers ---")
	mainLoopIdx := strings.Index(script, "# --- Main Loop ---")

	if validationIdx < 0 || authIdx < 0 || helpersIdx < 0 || mainLoopIdx < 0 {
		t.Fatal("missing expected section headers in generated script")
	}

	if validationIdx >= authIdx {
		t.Error("validation section should come before auth check section")
	}
	if authIdx >= helpersIdx {
		t.Error("auth check section should come before helpers section")
	}
	if helpersIdx >= mainLoopIdx {
		t.Error("helpers section should come before main loop section")
	}
}

func TestGenerateAutoScript_SetupClaudeNpmFallback(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 10,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	// npm fallback should still be present for claude
	if !strings.Contains(script, "@anthropic-ai/claude-code") {
		t.Error("expected npm fallback with claude-code package")
	}
	if !strings.Contains(script, "Native install failed, falling back to npm") {
		t.Error("expected npm fallback message when native install fails")
	}
}

func TestGenerateAutoScript_SetupEnsuresWritableHome(t *testing.T) {
	config := AutoConfig{
		MaxIterations: 10,
		AITool:        "claude",
		PromptFile:    ".claude/auto/prompt.md",
	}

	script := GenerateAutoScript(config)

	// ensure_writable_home should guard against unwritable HOME
	if !strings.Contains(script, "ensure_writable_home") {
		t.Error("expected ensure_writable_home function")
	}
	if !strings.Contains(script, `export HOME="/tmp"`) {
		t.Error("expected HOME=/tmp fallback in ensure_writable_home")
	}
	if !strings.Contains(script, `-d "$HOME"`) {
		t.Error("expected check that HOME directory exists")
	}
	if !strings.Contains(script, `-w "$HOME"`) {
		t.Error("expected check that HOME directory is writable")
	}
	if !strings.Contains(script, `"$HOME" != "/"`) {
		t.Error("expected guard against HOME=/ (Docker --user default)")
	}
}
