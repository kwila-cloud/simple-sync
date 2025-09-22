package models

// DockerContainer represents the configuration for a containerized application instance
type DockerContainer struct {
	Image         string            `json:"image,omitempty"`          // Container image name and tag
	Ports         []string          `json:"ports,omitempty"`          // Port mappings (host:container)
	Environment   map[string]string `json:"environment,omitempty"`    // Environment variables
	Volumes       []string          `json:"volumes,omitempty"`        // Volume mounts for data persistence
	RestartPolicy string            `json:"restart_policy,omitempty"` // Container restart behavior
	HealthCheck   *HealthCheck      `json:"health_check,omitempty"`   // Health check configuration
}

// HealthCheck represents container health check configuration
type HealthCheck struct {
	Test        []string `json:"test,omitempty"`         // Health check command
	Interval    string   `json:"interval,omitempty"`     // Check interval (e.g., "30s")
	Timeout     string   `json:"timeout,omitempty"`      // Check timeout (e.g., "10s")
	Retries     int      `json:"retries,omitempty"`      // Number of retries
	StartPeriod string   `json:"start_period,omitempty"` // Start period (e.g., "40s")
}

// NewDockerContainer creates a new Docker container configuration
func NewDockerContainer(image string) *DockerContainer {
	return &DockerContainer{
		Image:         image,
		Ports:         []string{},
		Environment:   make(map[string]string),
		Volumes:       []string{},
		RestartPolicy: "unless-stopped",
	}
}

// AddPort adds a port mapping to the container
func (dc *DockerContainer) AddPort(hostPort, containerPort string) {
	dc.Ports = append(dc.Ports, hostPort+":"+containerPort)
}

// SetEnvironment sets an environment variable
func (dc *DockerContainer) SetEnvironment(key, value string) {
	dc.Environment[key] = value
}

// AddVolume adds a volume mount
func (dc *DockerContainer) AddVolume(hostPath, containerPath string) {
	dc.Volumes = append(dc.Volumes, hostPath+":"+containerPath)
}

// SetHealthCheck configures the health check
func (dc *DockerContainer) SetHealthCheck(test []string, interval, timeout, startPeriod string, retries int) {
	dc.HealthCheck = &HealthCheck{
		Test:        test,
		Interval:    interval,
		Timeout:     timeout,
		Retries:     retries,
		StartPeriod: startPeriod,
	}
}
