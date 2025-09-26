# AGENTS.md - AI Development Guide for Simple-Sync

## Project Overview

Simple-sync is a lightweight REST API built in Go that provides event storage and access control functionality. The system allows users to authenticate, store timestamped events for specific items, and manage permissions through Access Control Lists (ACLs).

**Technology Stack:**
- Go 1.25 with Gin web framework
- SQLite database storage
- CORS support for web clients

**Core Features:**
- User authentication with API keys
- Event storage with timestamps and metadata
- ACL-based permission system (read/write permissions)
- Persistent SQLite database storage for data survival across restarts

## Key Files & Documentation

**Main Application Files:**
- `main.go` - Application entry point and server setup
- `handlers/` - HTTP endpoint handlers
- `models/` - Data structures for events, users, ACL
- `middleware/` - Authentication and CORS middleware
- `storage/` - SQLite persistence layer

**Configuration & Data:**
- `config/` - Environment variables and configuration
- `data/simple-sync.db` - SQLite database for persistent storage
- `data/backups/` - Database backup files for data safety

**Documentation:**
- `README.md` - Setup and deployment instructions
- `docs/src/content/docs/api/` - Complete API specification and examples
- `docs/src/content/docs/acl.md` - ACL system documentation and permission model
- `AGENTS.md` - This file - AI development guide
- `docker-compose.yml` - Local development environment
- Frontend integration examples in `/examples`

## GitHub CLI Instructions

## Development Workflow

**1. Issue Selection and Setup:**
```bash
# Get issue details
gh issue view <issue-number>

# Create feature branch
gh issue main <issue-number> --checkout
# Or manually: git checkout -b feat/issue-<number>-description
```

**2. Development Process:**
- Read the full issue description and acceptance criteria
- Check `docs/src/content/docs` for the project design
- Check dependencies listed in the issue
- Implement features incrementally, testing as you go
- Follow the existing code patterns and structure
- Update relevant documentation in `docs/` if making API changes

**3. Implementation Guidelines:**
- Start with data models in `models/` directory
- Add storage layer functions in `storage/`
- Implement HTTP handlers in `handlers/`
- Add middleware if needed in `middleware/`
- Update main.go to wire everything together
- Follow Go testing conventions: test files must end with "_test.go" (not start with "test_")

**4. Commit and Push:**
```bash
# Stage changes
git add .

# IMPORTANT: Update CHANGELOG.md for user-facing changes
# Add entry to CHANGELOG.md documenting new features, bug fixes, etc.

# Commit with descriptive message
git commit -m "feat: implement ACL endpoints with simple allow/deny logic

- Add ACL model and storage layer
- Implement GET /acl and PUT /acl endpoints
- Add permission checking to event endpoints
- Return 403 for unauthorized access

Closes #<issue-number>"

# Push to remote
git push origin feat/issue-<number>-description
```

**5. View Pull Request:**
```bash
# This will automatically get the PR for the current branch
gh pr view
```

**CHANGELOG Reminder:**
- **ALWAYS update CHANGELOG.md** for any pull request that introduces user-facing changes
- Document new features, enhancements, bug fixes, and breaking changes
- Follow the existing format with PR links and clear descriptions
- Keep entries concise but descriptive for users and maintainers

## Verification & Testing

**Authentication Testing:**
```bash
# Register new user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}'

# Login and get token
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}'

# Save token for subsequent requests
export TOKEN="your-jwt-token-here"
```

**Event Management Testing:**
```bash
# Create event
curl -X POST http://localhost:8080/events \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"itemUuid":"item-123","data":{"key":"value"}}'

# Get events
curl -X GET http://localhost:8080/events \
  -H "Authorization: Bearer $TOKEN"

# Get events for specific item
curl -X GET "http://localhost:8080/events?itemUuid=item-123" \
  -H "Authorization: Bearer $TOKEN"
```

**ACL Testing:**
```bash
# Set permissions
curl -X PUT http://localhost:8080/acl \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '[{"userUuid":"user-123","itemUuid":"item-123","permissions":["read","write"]}]'

# Get ACL entries
curl -X GET http://localhost:8080/acl \
  -H "Authorization: Bearer $TOKEN"
```

**Database Persistence Verification:**
```bash
# Check SQLite database exists and has data
ls -la data/
sqlite3 data/simple-sync.db "SELECT COUNT(*) FROM events;"
sqlite3 data/simple-sync.db "SELECT COUNT(*) FROM acl;"

# Restart server and verify data persists
# Stop server, restart, then test GET endpoints
```

**CORS Testing:**
```bash
# Test preflight request
curl -X OPTIONS http://localhost:8080/events \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET" \
  -H "Access-Control-Request-Headers: Authorization" \
  -v
```

## Example curl Commands for Testing

**Complete Authentication Flow:**
```bash
# 1. Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}'

# 2. Login
RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}')

# 3. Extract token (requires jq)
TOKEN=$(echo $RESPONSE | jq -r '.token')
echo "Token: $TOKEN"
```

**End-to-End Event Flow:**
```bash
# Create event
curl -X POST http://localhost:8080/events \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "itemUuid": "product-456",
    "data": {
      "action": "view",
      "userId": "user-789",
      "timestamp": "2025-01-15T10:30:00Z"
    }
  }'

# Verify event was stored
curl -X GET "http://localhost:8080/events?itemUuid=product-456" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

## Docker Setup Instructions

**Development Environment:**
```bash
# Build and start services
docker-compose up --build

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f simple-sync

# Stop services
docker-compose down
```

**Manual Docker Commands:**
```bash
# Build image
docker build -t simple-sync .

# Run with volume for data persistence
docker run -d \
  --name simple-sync \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  simple-sync

# Check container logs
docker logs simple-sync

# Access container shell
docker exec -it simple-sync /bin/sh
```

**Testing with Docker:**
```bash
# Test health endpoint
curl http://localhost:8080/health

# Full integration test
./scripts/integration-test.sh  # If available
```

## Code Structure & Patterns

**File Organization:**
```
simple-sync/
├── main.go                 # Application entry point
├── handlers/
│   ├── auth.go            # Authentication endpoints
│   ├── events.go          # Event management endpoints
│   └── acl.go             # ACL management endpoints
├── middleware/
│   ├── auth.go            # API key authentication middleware
│   └── cors.go            # CORS middleware
├── models/
│   ├── user.go            # User data structures
│   ├── event.go           # Event data structures
│   └── acl.go             # ACL data structures
├── storage/
│   ├── interface.go       # Storage interface definition
│   ├── sqlite.go          # SQLite storage implementation
│   └── memory.go          # In-memory storage (for tests)
└── config/
    └── config.go          # Configuration management
```

**Coding Patterns to Follow:**

1. **Error Handling:**
```go
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
    return
}
```

2. **Response Format:**
```go
// Success responses
c.JSON(http.StatusOK, gin.H{"data": result})

// Error responses
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
```

3. **Middleware Usage:**
```go
// Protected routes
protected := router.Group("/")
protected.Use(middleware.AuthMiddleware())
protected.GET("/events", handlers.GetEvents)
```

4. **Storage Interface:**
```go
type Storage interface {
    SaveEvents(events []Event) error
    LoadEvents() ([]Event, error)
    SaveACL(acl []ACLEntry) error
    LoadACL() ([]ACLEntry, error)
}
```

## Common Troubleshooting

**Authentication Issues:**
- Check API key format in Authorization header: "Bearer <api-key>"
- Validate user exists in storage

**Database Issues:**
- Check data directory permissions and SQLite database file permissions
- Verify disk space availability
- Look for SQLite errors in logs (corruption, locking, etc.)
- Check for concurrent access issues with WAL mode
- Verify SQLite database integrity with PRAGMA integrity_check

**CORS Problems:**
- Verify Origin header in requests
- Check preflight OPTIONS responses
- Ensure all required headers are allowed
- Test with browser developer tools

**Performance Issues:**
- Monitor file I/O operations
- Check for file locking bottlenecks
- Consider data file sizes
- Verify concurrent access handling

**Common Error Codes:**
- 401: Invalid or missing authentication
- 403: Insufficient permissions (ACL)
- 400: Invalid request format
- 500: Server/storage errors

Remember to always test your changes incrementally and verify that existing functionality still works after implementing new features.

## Debug Commands
```bash
# Check server logs
docker-compose logs simple-sync

# Verify SQLite database exists
ls -la data/
sqlite3 data/simple-sync.db ".tables"
sqlite3 data/simple-sync.db "SELECT name FROM sqlite_master WHERE type='table';"

# Inspect events table
sqlite3 data/simple-sync.db "SELECT * FROM events LIMIT 5;"

# Inspect ACL table
sqlite3 data/simple-sync.db "SELECT * FROM acl LIMIT 5;"

# Restart server and verify data persists
# Stop server, restart, then test GET endpoints
```
