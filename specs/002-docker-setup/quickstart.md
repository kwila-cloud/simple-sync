# Docker Quick Start Guide

## Prerequisites

- Docker installed and running
- Docker Compose installed
- Git repository cloned

## Quick Start

1. **Clone and navigate to the repository**
   ```bash
   git clone <repository-url>
   cd simple-sync
   ```

2. **Set required environment variable**
   ```bash
   export JWT_SECRET="your-secure-jwt-secret-here"
   # Generate a secure secret: openssl rand -base64 32
   ```

3. **Start the service with Docker Compose**
   ```bash
   docker-compose up
   ```

4. **Verify the service is running**
   - Open http://localhost:8080/health in your browser
   - Expected response: `{"status":"healthy","timestamp":"2025-09-20T...","version":"v0.1.0","uptime":30}`

5. **Test the API**
   ```bash
   # Get authentication token
   curl -X POST http://localhost:8080/auth/token \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","password":"testpass123"}'

   # Use token for protected endpoints
   curl -X GET http://localhost:8080/events \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

## Configuration Options

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JWT_SECRET` | Yes | - | Secret key for JWT token signing |
| `PORT` | No | 8080 | Port the service listens on |

### Docker Compose Override

Create a `docker-compose.override.yml` for custom configuration:

```yaml
version: '3.8'
services:
  simple-sync:
    environment:
      - JWT_SECRET=your-production-secret
      - PORT=3000
    ports:
      - "3000:3000"  # Match PORT setting
    volumes:
      - ./data:/app/data  # Persist SQLite database
```

## Development Workflow

### With Hot Reload (for development)
```bash
# Use volume mounts for source code changes
docker-compose -f docker-compose.dev.yml up
```

### Building Custom Image
```bash
# Build the image
docker build -t simple-sync-custom .

# Run with custom configuration
docker run -p 8080:8080 \
  -e JWT_SECRET="your-secret" \
  simple-sync-custom
```

## Troubleshooting

### Service Won't Start
```bash
# Check if JWT_SECRET is set
echo $JWT_SECRET

# Check Docker Compose logs
docker-compose logs simple-sync

# Check container status
docker-compose ps
```

### Port Already in Use
```bash
# Change the port in docker-compose.yml
ports:
  - "8081:8080"  # Host port 8081, container port 8080

# Or set PORT environment variable
export PORT=8081
```

### Database Persistence
```bash
# Data is stored in ./data directory by default
ls -la data/

# Check SQLite database
docker-compose exec simple-sync sqlite3 /app/data/simple-sync.db ".tables"
```

### Health Check Failing
```bash
# Check service logs
docker-compose logs simple-sync

# Manual health check
curl http://localhost:8080/health

# Check container resource usage
docker-compose stats
```

## Production Deployment

### Using Docker Compose in Production
```bash
# Create production override
cp docker-compose.yml docker-compose.prod.yml
# Edit docker-compose.prod.yml with production settings

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

### Using Docker Registry
```bash
# Tag and push image
docker tag simple-sync your-registry/simple-sync:v1.0.0
docker push your-registry/simple-sync:v1.0.0

# Use in production docker-compose.yml
image: your-registry/simple-sync:v1.0.0
```

## Security Notes

- **Never commit JWT_SECRET** to version control
- Use strong, randomly generated secrets for JWT_SECRET
- Rotate JWT_SECRET periodically in production
- Consider using Docker secrets for sensitive configuration
- Run containers as non-root user (already configured)

## Performance Tuning

### Resource Limits
```yaml
services:
  simple-sync:
    deploy:
      resources:
        limits:
          memory: 256M
          cpus: '0.5'
        reservations:
          memory: 128M
          cpus: '0.25'
```

### Health Check Configuration
The container includes health checks configured for:
- Interval: 30 seconds
- Timeout: 10 seconds
- Retries: 3
- Start period: 40 seconds

## Next Steps

- Explore the API documentation in `docs/api.md`
- Review ACL configuration in `docs/acl.md`
- Check test coverage with `go test ./tests/...`
- Customize configuration for your deployment environment