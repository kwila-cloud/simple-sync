package integration

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentVariableConfiguration(t *testing.T) {
	// Skip if Docker is not available
	if !isDockerAvailable() {
		t.Skip("Docker not available, skipping environment config test")
	}

	// Test with custom PORT
	t.Run("CustomPort", func(t *testing.T) {
		testContainerWithEnv(t, "8082", "test-jwt-secret-custom-port-at-least-32-characters-long")
	})

	// Test with default PORT
	t.Run("DefaultPort", func(t *testing.T) {
		testContainerWithEnv(t, "", "test-jwt-secret-default-port-at-least-32-characters-long")
	})
}

func testContainerWithEnv(t *testing.T, port, jwtSecret string) {
	ctx := t.Context()

	// Build image if not exists
	buildCmd := exec.CommandContext(ctx, "docker", "build", "-t", "simple-sync-env-test", "../..")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build image: %v\nOutput: %s", err, string(output))
	}

	// Prepare docker run command
	args := []string{"run", "-d", "--name", "simple-sync-env-test-" + port}

	// Set environment variables
	if jwtSecret != "" {
		args = append(args, "-e", "JWT_SECRET="+jwtSecret)
	}
	if port != "" {
		args = append(args, "-e", "PORT="+port)
		args = append(args, "-p", port+":"+port)
	} else {
		args = append(args, "-p", "8083:8080") // Default port mapping
	}

	args = append(args, "simple-sync-env-test")

	// Run container
	runCmd := exec.CommandContext(ctx, "docker", args...)
	containerID, err := runCmd.Output()
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	containerIDStr := strings.TrimSpace(string(containerID))

	// Debug: show the docker run command
	t.Logf("Docker run command: docker %s", strings.Join(args, " "))

	// Clean up
	defer func() {
		exec.Command("docker", "rm", "-f", containerIDStr).Run()
	}()

	// Wait for startup
	time.Sleep(10 * time.Second)

	// Check if container is running (including stopped containers)
	psCmd := exec.Command("docker", "ps", "-a", "--filter", "id="+containerIDStr, "--format", "{{.Status}}")
	status, err := psCmd.Output()
	if err != nil {
		t.Fatalf("Failed to check status: %v", err)
	}

	// Get container logs for debugging
	logsCmd := exec.Command("docker", "logs", containerIDStr)
	logs, _ := logsCmd.Output()
	t.Logf("Container status: %s", string(status))
	t.Logf("Container logs: %s", string(logs))

	assert.Contains(t, string(status), "Up", "Container should be running")

	// Verify JWT_SECRET is required (container should fail without it)
	if jwtSecret == "" {
		// If no JWT_SECRET provided, container should not be healthy
		logsCmd := exec.Command("docker", "logs", containerIDStr)
		logs, _ := logsCmd.Output()
		assert.Contains(t, string(logs), "JWT_SECRET", "Should mention missing JWT_SECRET")
	} else {
		// If JWT_SECRET provided, health check should work
		// Check the internal port (PORT env var or 8080 default) since we're running curl inside the container
		internalPort := port
		if internalPort == "" {
			internalPort = "8080" // Default port
		}
		healthCmd := exec.Command("docker", "exec", containerIDStr,
			"curl", "-f", "-s", "http://localhost:"+internalPort+"/health")
		err := healthCmd.Run()
		assert.NoError(t, err, "Health check should pass with valid JWT_SECRET")
	}
}
