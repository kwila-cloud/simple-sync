package models

import (
	"errors"
	"strconv"
)

// EnvironmentConfiguration manages environment-specific settings
type EnvironmentConfiguration struct {
	JWT_SECRET  string `json:"jwt_secret"`  // Required JWT signing secret
	Port        int    `json:"port"`        // Service port number (default 8080)
	Environment string `json:"environment"` // Deployment environment (development/production)
}

// NewEnvironmentConfiguration creates a new environment configuration with defaults
func NewEnvironmentConfiguration() *EnvironmentConfiguration {
	return &EnvironmentConfiguration{
		Port:        8080,
		Environment: "development",
	}
}

// LoadFromEnv loads configuration from environment variables
func (ec *EnvironmentConfiguration) LoadFromEnv(getenv func(string) string) error {
	// JWT_SECRET is required
	jwtSecret := getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("JWT_SECRET environment variable is required")
	}
	ec.JWT_SECRET = jwtSecret

	// PORT is optional, defaults to 8080
	portStr := getenv("PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return errors.New("PORT must be a valid integer")
		}
		if port < 1024 || port > 65535 {
			return errors.New("PORT must be between 1024 and 65535")
		}
		ec.Port = port
	}

	// ENVIRONMENT is optional, defaults to "development"
	env := getenv("ENVIRONMENT")
	if env != "" {
		ec.Environment = env
	}

	return nil
}

// Validate checks if the configuration is valid
func (ec *EnvironmentConfiguration) Validate() error {
	if ec.JWT_SECRET == "" {
		return errors.New("JWT_SECRET is required")
	}

	if len(ec.JWT_SECRET) < 32 {
		return errors.New("JWT_SECRET should be at least 32 characters long")
	}

	if ec.Port < 1024 || ec.Port > 65535 {
		return errors.New("PORT must be between 1024 and 65535")
	}

	return nil
}

// IsProduction returns true if running in production environment
func (ec *EnvironmentConfiguration) IsProduction() bool {
	return ec.Environment == "production"
}

// IsDevelopment returns true if running in development environment
func (ec *EnvironmentConfiguration) IsDevelopment() bool {
	return ec.Environment == "development"
}
