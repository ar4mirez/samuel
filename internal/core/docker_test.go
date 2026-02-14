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
	if len(modes) != 2 {
		t.Fatalf("expected 2 modes, got %d", len(modes))
	}
	if modes[0] != SandboxNone {
		t.Errorf("expected first mode %q, got %q", SandboxNone, modes[0])
	}
	if modes[1] != SandboxDocker {
		t.Errorf("expected second mode %q, got %q", SandboxDocker, modes[1])
	}
}

func TestIsValidSandboxMode(t *testing.T) {
	tests := []struct {
		mode string
		want bool
	}{
		{"none", true},
		{"docker", true},
		{"podman", false},
		{"", false},
		{"DOCKER", true},
		{"Docker", true},
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
			name: "default image when empty",
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

func TestCheckDockerAvailable(t *testing.T) {
	// Only run if docker is available on the test machine
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not available, skipping")
	}

	err := CheckDockerAvailable()
	// We can't assert success since the daemon may not be running,
	// but we can verify the function doesn't panic
	if err != nil {
		t.Logf("CheckDockerAvailable returned: %v (daemon may not be running)", err)
	}
}
