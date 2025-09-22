# Docker Configuration Data Model

## Overview
This feature adds Docker containerization support to the simple-sync service. The data model focuses on configuration entities and runtime requirements rather than application data models.

## Configuration Entities

### DockerContainer
**Purpose**: Represents the containerized application instance
**Attributes**:
- `image`: Container image name and tag
- `ports`: Port mappings (host:container)
- `environment`: Environment variables
- `volumes`: Volume mounts for data persistence
- `restart_policy`: Container restart behavior
- `health_check`: Health check configuration

**Validation Rules**:
- Image must be specified
- Ports must include host:container mapping
- JWT_SECRET environment variable must be provided
- PORT defaults to 8080 if not specified

### EnvironmentConfiguration
**Purpose**: Manages environment-specific settings
**Attributes**:
- `jwt_secret`: Required JWT signing secret (string, 32+ characters recommended)
- `port`: Service port number (integer, default 8080, range 80-65535)
- `environment`: Deployment environment (development/production)

**Validation Rules**:
- jwt_secret is mandatory and cannot be empty
- port must be valid TCP port number
- environment affects logging and debug settings

### HealthCheckResponse
**Purpose**: Standard response format for container health checks
**Attributes**:
- `status`: Health status ("healthy" | "unhealthy")
- `timestamp`: ISO 8601 timestamp of check
- `version`: Application version string
- `uptime`: Service uptime in seconds

**Validation Rules**:
- status must be either "healthy" or "unhealthy"
- timestamp must be valid ISO 8601 format
- version should match application version
- uptime should be non-negative integer

## Relationships

```
DockerContainer
├── requires EnvironmentConfiguration
├── provides HealthCheckResponse
└── persists data via volumes

EnvironmentConfiguration
├── configures DockerContainer
└── validates jwt_secret and port
```

## State Transitions

### Container Lifecycle
1. **Created**: Container image built, configuration validated
2. **Starting**: Container starting, environment variables loaded
3. **Healthy**: Health check passes, service accepting requests
4. **Unhealthy**: Health check fails, container may restart
5. **Stopped**: Container terminated, data preserved in volumes

### Configuration States
1. **Valid**: All required environment variables present and valid
2. **Invalid**: Missing or invalid configuration, container fails to start
3. **Development**: Default settings for local development
4. **Production**: Secure settings for production deployment

## Data Flow

### Startup Sequence
1. Docker loads environment variables
2. Application validates JWT_SECRET requirement
3. SQLite database initializes (if volume mounted)
4. Authentication service loads default user
5. HTTP server starts on configured port
6. Health check endpoint becomes available

### Runtime Operation
1. Client requests hit container on exposed port
2. Environment variables configure JWT validation
3. SQLite database persists data to mounted volume
4. Health checks monitor service status
5. Logs output to container stdout/stderr

## Security Considerations

- JWT_SECRET must be provided externally (never baked into image)
- Container runs as non-root user
- Sensitive data stored in environment variables only
- No secrets committed to version control
- Database files protected by volume permissions

## Performance Characteristics

- **Startup Time**: <5 seconds for container initialization
- **Memory Usage**: <100MB for base application
- **Health Check Response**: <10ms
- **Port Availability**: Configurable to avoid conflicts
- **Volume I/O**: SQLite performance maintained through host filesystem
