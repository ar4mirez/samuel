package core

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGetSupportedSandboxModes(t *testing.T) {
	modes := GetSupportedSandboxModes()
	if len(modes) != 3 {
		t.Fatalf("expected 3 modes, got %d", len(modes))
	}
	if modes[0] != SandboxNone {
		t.Errorf("expected first mode %q, got %q", SandboxNone, modes[0])
	}
	if modes[1] != SandboxDocker {
		t.Errorf("expected second mode %q, got %q", SandboxDocker, modes[1])
	}
	if modes[2] != SandboxDockerSandbox {
		t.Errorf("expected third mode %q, got %q", SandboxDockerSandbox, modes[2])
	}
}

func TestIsValidSandboxMode(t *testing.T) {
	tests := []struct {
		mode string
		want bool
	}{
		{"none", true},
		{"docker", true},
		{"docker-sandbox", true},
		{"podman", false},
		{"", false},
		{"DOCKER", true},
		{"Docker", true},
		{"DOCKER-SANDBOX", true},
		{"Docker-Sandbox", true},
		{"kubernetes", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("mode=%q", tt.mode), func(t *testing.T) {
			got := IsValidSandboxMode(tt.mode)
			if got != tt.want {
				t.Errorf("IsValidSandboxMode(%q) = %v, want %v", tt.mode, got, tt.want)
			}
		})
	}
}

func TestBuildDockerArgs(t *testing.T) {
	// Clear env vars that would add -e flags to keep assertions predictable
	for _, name := range aiToolEnvVarNames {
		t.Setenv(name, "")
		os.Unsetenv(name)
	}

	tests := []struct {
		name         string
		config       DockerSandboxConfig
		wantImage    string
		wantScript   string
		wantIterFlag bool
		wantIterVal  string
	}{
		{
			name: "basic config",
			config: DockerSandboxConfig{
				Image:      "myimage:latest",
				WorkDir:    "/home/user/project",
				ScriptPath: ".claude/auto/auto.sh",
			},
			wantImage:    "myimage:latest",
			wantScript:   "/workspace/.claude/auto/auto.sh",
			wantIterFlag: false,
		},
		{
			name: "default image when empty (node:lts)",
			config: DockerSandboxConfig{
				Image:      "",
				WorkDir:    "/tmp/proj",
				ScriptPath: ".claude/auto/auto.sh",
			},
			wantImage:    DefaultSandboxImage,
			wantScript:   "/workspace/.claude/auto/auto.sh",
			wantIterFlag: false,
		},
		{
			name: "with iteration override",
			config: DockerSandboxConfig{
				Image:        "ubuntu:22.04",
				WorkDir:      "/home/user/project",
				ScriptPath:   ".claude/auto/auto.sh",
				IterOverride: 10,
			},
			wantImage:    "ubuntu:22.04",
			wantScript:   "/workspace/.claude/auto/auto.sh",
			wantIterFlag: true,
			wantIterVal:  "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := BuildDockerArgs(tt.config)

			// Verify it starts with "run --rm --init"
			if args[0] != "run" || args[1] != "--rm" || args[2] != "--init" {
				t.Errorf("expected args to start with [run --rm --init], got %v", args[:3])
			}

			// Find the image in args (appears before "bash")
			bashIdx := -1
			for i, a := range args {
				if a == "bash" {
					bashIdx = i
					break
				}
			}
			if bashIdx < 0 {
				t.Fatal("expected 'bash' in args, not found")
			}

			gotImage := args[bashIdx-1]
			if gotImage != tt.wantImage {
				t.Errorf("image = %q, want %q", gotImage, tt.wantImage)
			}

			gotScript := args[bashIdx+1]
			if gotScript != tt.wantScript {
				t.Errorf("script = %q, want %q", gotScript, tt.wantScript)
			}

			// Check volume mount
			joined := strings.Join(args, " ")
			wantMount := fmt.Sprintf("-v %s:%s", tt.config.WorkDir, DockerContainerMount)
			if !strings.Contains(joined, wantMount) {
				t.Errorf("expected mount %q in args: %s", wantMount, joined)
			}

			// Check working directory
			if !strings.Contains(joined, fmt.Sprintf("-w %s", DockerContainerMount)) {
				t.Errorf("expected -w %s in args: %s", DockerContainerMount, joined)
			}

			// Check --user flag is present
			if !strings.Contains(joined, "--user") {
				t.Error("expected --user flag in args")
			}

			// Check iteration override
			if tt.wantIterFlag {
				lastArg := args[len(args)-1]
				if lastArg != tt.wantIterVal {
					t.Errorf("last arg = %q, want iteration %q", lastArg, tt.wantIterVal)
				}
			}
		})
	}
}

func TestBuildDockerSandboxArgs(t *testing.T) {
	tests := []struct {
		name     string
		config   DockerSandboxRunConfig
		wantArgs []string
	}{
		{
			name: "basic claude sandbox",
			config: DockerSandboxRunConfig{
				Agent:   "claude",
				WorkDir: "/home/user/project",
			},
			wantArgs: []string{
				"sandbox", "run", "claude", "/home/user/project",
			},
		},
		{
			name: "with template",
			config: DockerSandboxRunConfig{
				Agent:    "claude",
				WorkDir:  "/home/user/project",
				Template: "python:3-alpine",
			},
			wantArgs: []string{
				"sandbox", "run", "--template", "python:3-alpine",
				"claude", "/home/user/project",
			},
		},
		{
			name: "with name",
			config: DockerSandboxRunConfig{
				Agent:   "claude",
				WorkDir: "/home/user/project",
				Name:    "my-sandbox",
			},
			wantArgs: []string{
				"sandbox", "run", "--name", "my-sandbox",
				"claude", "/home/user/project",
			},
		},
		{
			name: "with agent args",
			config: DockerSandboxRunConfig{
				Agent:   "claude",
				WorkDir: "/home/user/project",
				AgentArgs: []string{
					"--print", "--dangerously-skip-permissions",
					"/home/user/project/.claude/auto/prompt.md",
				},
			},
			wantArgs: []string{
				"sandbox", "run", "claude", "/home/user/project",
				"--",
				"--print", "--dangerously-skip-permissions",
				"/home/user/project/.claude/auto/prompt.md",
			},
		},
		{
			name: "default agent when empty",
			config: DockerSandboxRunConfig{
				Agent:   "",
				WorkDir: "/tmp/proj",
			},
			wantArgs: []string{
				"sandbox", "run", DefaultDockerSandboxAgent, "/tmp/proj",
			},
		},
		{
			name: "default workdir when empty",
			config: DockerSandboxRunConfig{
				Agent:   "claude",
				WorkDir: "",
			},
			wantArgs: []string{
				"sandbox", "run", "claude", ".",
			},
		},
		{
			name: "with name and template and args",
			config: DockerSandboxRunConfig{
				Agent:     "codex",
				WorkDir:   "/proj",
				Name:      "test-run",
				Template:  "node:20",
				AgentArgs: []string{"--prompt-file", "/proj/prompt.md"},
			},
			wantArgs: []string{
				"sandbox", "run", "--name", "test-run",
				"--template", "node:20",
				"codex", "/proj",
				"--", "--prompt-file", "/proj/prompt.md",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := BuildDockerSandboxArgs(tt.config)
			if len(args) != len(tt.wantArgs) {
				t.Fatalf("got %d args %v, want %d args %v",
					len(args), args, len(tt.wantArgs), tt.wantArgs)
			}
			for i, got := range args {
				if got != tt.wantArgs[i] {
					t.Errorf("arg[%d] = %q, want %q", i, got, tt.wantArgs[i])
				}
			}
		})
	}
}

func TestGetAgentArgs(t *testing.T) {
	tests := []struct {
		aiTool     string
		promptPath string
		wantArgs   []string
	}{
		{
			"claude",
			"/path/prompt.md",
			[]string{"--print", "--dangerously-skip-permissions", "/path/prompt.md"},
		},
		{
			"codex",
			"/path/prompt.md",
			[]string{"--prompt-file", "/path/prompt.md", "--auto"},
		},
		{
			"amp",
			"/path/prompt.md",
			[]string{"--prompt-file", "/path/prompt.md"},
		},
		{
			"unknown-tool",
			"/path/prompt.md",
			[]string{"/path/prompt.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.aiTool, func(t *testing.T) {
			args := GetAgentArgs(tt.aiTool, tt.promptPath)
			if len(args) != len(tt.wantArgs) {
				t.Fatalf("got %d args %v, want %d args %v",
					len(args), args, len(tt.wantArgs), tt.wantArgs)
			}
			for i, got := range args {
				if got != tt.wantArgs[i] {
					t.Errorf("arg[%d] = %q, want %q", i, got, tt.wantArgs[i])
				}
			}
		})
	}
}

func TestGetAIToolEnvVars(t *testing.T) {
	// Clear all tracked env vars first
	for _, name := range aiToolEnvVarNames {
		os.Unsetenv(name)
	}

	t.Setenv("ANTHROPIC_API_KEY", "sk-test-123")
	t.Setenv("AI_TOOL", "claude")

	envArgs := getAIToolEnvVars()

	// Should contain -e pairs for the two vars we set
	joined := strings.Join(envArgs, " ")
	if !strings.Contains(joined, "-e ANTHROPIC_API_KEY=sk-test-123") {
		t.Errorf("expected ANTHROPIC_API_KEY in env args: %v", envArgs)
	}
	if !strings.Contains(joined, "-e AI_TOOL=claude") {
		t.Errorf("expected AI_TOOL in env args: %v", envArgs)
	}

	// Should NOT contain vars we didn't set
	if strings.Contains(joined, "OPENAI_API_KEY") {
		t.Errorf("did not expect OPENAI_API_KEY in env args: %v", envArgs)
	}
}

func TestGetAIToolEnvVars_NoHomeOrPath(t *testing.T) {
	// HOME and PATH must NOT be forwarded — they leak host paths into
	// the container where they don't exist.
	for _, name := range aiToolEnvVarNames {
		os.Unsetenv(name)
	}

	// Set HOME and PATH to verify they are NOT in the forwarded list
	t.Setenv("HOME", "/Users/testuser")
	t.Setenv("PATH", "/usr/local/bin:/usr/bin")
	t.Setenv("ANTHROPIC_API_KEY", "sk-test")

	envArgs := getAIToolEnvVars()
	joined := strings.Join(envArgs, " ")

	if strings.Contains(joined, "HOME=") {
		t.Errorf("HOME should NOT be forwarded to container: %v", envArgs)
	}
	if strings.Contains(joined, "PATH=") {
		t.Errorf("PATH should NOT be forwarded to container: %v", envArgs)
	}
}

func TestCheckDockerAvailable(t *testing.T) {
	// Only run if docker is available on the test machine
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not available, skipping")
	}

	err := CheckDockerAvailable()
	if err != nil {
		t.Logf("CheckDockerAvailable returned: %v (daemon may not be running)", err)
	}
}

func TestCheckDockerSandboxAvailable(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not available, skipping")
	}

	err := CheckDockerSandboxAvailable()
	if err != nil {
		t.Logf("CheckDockerSandboxAvailable returned: %v (sandbox plugin may not be installed)", err)
	}
}
