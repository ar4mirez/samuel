package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"golang.org/x/term"
)

// Sandbox mode constants
const (
	SandboxNone          = "none"
	SandboxDocker        = "docker"
	DefaultSandboxImage  = "ubuntu:latest"
	DockerContainerMount = "/workspace"
)

// aiToolEnvVarNames is the allowlist of environment variables passed into the
// Docker container. Only variables that are actually set on the host are
// forwarded, preventing accidental secret leakage.
var aiToolEnvVarNames = []string{
	"ANTHROPIC_API_KEY",
	"OPENAI_API_KEY",
	"AMP_API_KEY",
	"AI_TOOL",
	"PAUSE_SECONDS",
	"HOME",
	"TERM",
	"PATH",
}

// DockerSandboxConfig holds the parameters needed to build a docker run command.
type DockerSandboxConfig struct {
	Image        string
	WorkDir      string // host path to bind-mount
	ScriptPath   string // relative path to auto.sh inside mount
	IterOverride int
}

// GetSupportedSandboxModes returns the list of supported sandbox modes.
func GetSupportedSandboxModes() []string {
	return []string{SandboxNone, SandboxDocker}
}

// IsValidSandboxMode checks if the given mode is supported.
func IsValidSandboxMode(mode string) bool {
	return slices.Contains(GetSupportedSandboxModes(), strings.ToLower(mode))
}

// CheckDockerAvailable verifies Docker is installed and the daemon is running.
func CheckDockerAvailable() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker not found in PATH; install Docker or use --sandbox=none")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "docker", "info").Run(); err != nil {
		return fmt.Errorf("docker daemon is not running; start Docker Desktop or the docker service")
	}
	return nil
}

// BuildDockerArgs constructs the argument slice for docker run.
func BuildDockerArgs(config DockerSandboxConfig) []string {
	image := config.Image
	if image == "" {
		image = DefaultSandboxImage
	}

	args := []string{"run", "--rm", "--init"}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		args = append(args, "-it")
	} else {
		args = append(args, "-i")
	}

	// Run as host user to avoid root-owned files
	args = append(args, "--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))

	// Bind-mount project directory
	args = append(args, "-v", fmt.Sprintf("%s:%s", config.WorkDir, DockerContainerMount))
	args = append(args, "-w", DockerContainerMount)

	// Pass through AI-tool environment variables
	args = append(args, getAIToolEnvVars()...)

	args = append(args, image)

	// Command inside container
	scriptInContainer := fmt.Sprintf("%s/%s", DockerContainerMount, config.ScriptPath)
	args = append(args, "bash", scriptInContainer)

	if config.IterOverride > 0 {
		args = append(args, fmt.Sprintf("%d", config.IterOverride))
	}

	return args
}

// getAIToolEnvVars returns -e KEY=VALUE pairs for set environment variables.
func getAIToolEnvVars() []string {
	var envArgs []string
	for _, name := range aiToolEnvVarNames {
		if val, ok := os.LookupEnv(name); ok {
			envArgs = append(envArgs, "-e", fmt.Sprintf("%s=%s", name, val))
		}
	}
	return envArgs
}
