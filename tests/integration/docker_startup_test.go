package integration

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDockerContainerStartup(t *testing.T) {
	// Skip if Docker is not available
	if !isDockerAvailable() {
		t.Skip("Docker not available, skipping container startup test")
	}

	// Build Docker image
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	buildCmd := exec.CommandContext(ctx, "docker", "build", "-t", "simple-sync-test", "../..")
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build Docker image: %v\nOutput: %s", err, string(output))
	}

	// Run container with test environment
	runCmd := exec.CommandContext(ctx, "docker", "run", "-d",
		"--name", "simple-sync-test-container",
		"-p", "8081:8080",
		"-e", "JWT_SECRET=test-docker-secret-key-for-testing-only-at-least-32-chars",
		"-e", "PORT=8080",
		"simple-sync-test")

	containerID, err := runCmd.Output()
	if err != nil {
		t.Fatalf("Failed to start Docker container: %v", err)
	}
	containerIDStr := strings.TrimSpace(string(containerID))

	// Clean up container when test finishes
	defer func() {
		exec.Command("docker", "rm", "-f", containerIDStr).Run()
	}()

	// Wait for container to be healthy
	time.Sleep(10 * time.Second) // Give container time to start

	// Check if container is running
	psCmd := exec.Command("docker", "ps", "--filter", "id="+containerIDStr, "--format", "{{.Status}}")
	status, err := psCmd.Output()
	if err != nil {
		t.Fatalf("Failed to check container status: %v", err)
	}

	// Container should be running (status should contain "Up")
	assert.Contains(t, string(status), "Up", "Container should be running")

	// Check health endpoint
	healthCmd := exec.Command("docker", "exec", containerIDStr,
		"curl", "-f", "-s", "http://localhost:8080/health")
	err = healthCmd.Run()
	if err != nil {
		// Get container logs for debugging
		logsCmd := exec.Command("docker", "logs", containerIDStr)
		logs, _ := logsCmd.Output()
		t.Logf("Container logs: %s", string(logs))
	}
	assert.NoError(t, err, "Health check should pass")
}

func isDockerAvailable() bool {
	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	return err == nil
}
