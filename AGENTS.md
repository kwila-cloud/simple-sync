# AGENTS.md - AI Development Guide for Simple-Sync

## Project Overview

Simple-sync is a lightweight REST API built in Go that provides event storage and access control functionality. The system allows users to authenticate via setup tokens exchanged for API keys, store timestamped events for specific items, and manage permissions through Access Control Lists (ACLs).

**Technology Stack:**
- Go 1.25 with Gin web framework
- SQLite database storage
- CORS support for web clients

**Core Features:**
- User authentication with API keys
- Event storage with timestamps and metadata
- ACL-based permission system (read/write permissions)

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

**Event Management Testing:**
```bash
# Create event
curl -X POST http://localhost:8080/api/v1/events \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '[{"uuid":"event-123","timestamp":1640995200,"user":"testuser","item":"item-123","action":"create","payload":"{}"}]'

# Get events
curl -X GET http://localhost:8080/api/v1/events \
  -H "X-API-Key: $API_KEY"
```

**CORS Testing:**
```bash
# Test preflight request
curl -X OPTIONS http://localhost:8080/api/v1/events \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET" \
  -H "Access-Control-Request-Headers: X-API-Key" \
  -v
```

## Example curl Commands for Testing

**Complete Authentication Flow:**
```bash
# 1. Generate setup token (requires admin API key)
RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/user/generateToken?user=testuser" \
  -H "X-API-Key: sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E")

# 2. Extract setup token (requires jq)
SETUP_TOKEN=$(echo $RESPONSE | jq -r '.token')
echo "Setup Token: $SETUP_TOKEN"

# 3. Exchange for API key
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/user/exchangeToken \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$SETUP_TOKEN\"}")

# 4. Extract API key (requires jq)
API_KEY=$(echo $RESPONSE | jq -r '.apiKey')
echo "API Key: $API_KEY"
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

