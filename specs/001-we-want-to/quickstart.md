# Quick Start: High Performance REST API for simple-sync

## Prerequisites
- Docker and Docker Compose installed
- Git repository cloned

## Setup
1. Clone the repository
2. Navigate to project root
3. Set environment variables:
   ```bash
   export JWT_SECRET="your-secure-secret-here"
   export PORT=8080
   ```

## Deployment
1. Build and start the service:
   ```bash
   docker-compose up --build
   ```

2. Verify service is running:
   ```bash
   curl http://localhost:8080/health
   ```

## Basic Usage
1. Register a user (if admin endpoint available):
   ```bash
   curl -X POST http://localhost:8080/admin/users \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","password":"testpass"}'
   ```

2. Authenticate and get token:
   ```bash
   TOKEN=$(curl -s -X POST http://localhost:8080/auth/token \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","password":"testpass"}' | jq -r '.token')
   ```

3. Post an event:
   ```bash
   curl -X POST http://localhost:8080/events \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '[{"uuid":"123e4567-e89b-12d3-a456-426614174000","timestamp":1640995200,"userUuid":"user123","itemUuid":"item456","action":"create","payload":"{}"}]'
   ```

4. Retrieve events:
   ```bash
   curl -X GET http://localhost:8080/events \
     -H "Authorization: Bearer $TOKEN"
   ```

## ACL Management
1. Set permissions:
   ```bash
   curl -X PUT http://localhost:8080/acl \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '[{"userUuid":"user123","itemUuid":"item456","permissions":["read","write"]}]'
   ```

2. View ACL:
   ```bash
   curl -X GET http://localhost:8080/acl \
     -H "Authorization: Bearer $TOKEN"
   ```

## Testing
- Run contract tests: `go test ./tests/contract/...`
- Run integration tests: `go test ./tests/integration/...`
- All tests should pass after implementation

## Troubleshooting
- Check Docker logs: `docker-compose logs`
- Verify environment variables are set
- Ensure data directory exists and is writable
- Check JWT secret is consistent across restarts