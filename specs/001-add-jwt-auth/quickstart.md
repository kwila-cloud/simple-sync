# Quickstart: Add JWT Authentication

**Date**: Sat Sep 20 2025
**Feature**: 001-add-jwt-auth
**Estimated Time**: 30 minutes

## Overview
This quickstart validates the JWT authentication implementation by testing the core authentication flow and protected endpoint access.

## Prerequisites
- Go 1.25 installed
- Simple-sync server running on localhost:8080
- JWT authentication feature implemented
- Test user credentials configured

## Test Scenarios

### Scenario 1: Successful Authentication
**Given** a user has valid credentials
**When** they request a token via POST /auth/token
**Then** they receive a valid JWT token

```bash
# Request authentication token
curl -X POST http://localhost:8080/auth/token \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "testpass123"}'

# Expected response:
# {
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
# }
```

**Success Criteria**:
- HTTP 200 status code
- Response contains "token" field
- Token is a valid JWT string

### Scenario 2: Access Protected Endpoint with Valid Token
**Given** a user has a valid JWT token
**When** they access GET /events with the token
**Then** the request succeeds and returns event data

```bash
# First get a token (from Scenario 1)
TOKEN="your-jwt-token-here"

# Access protected endpoint
curl -X GET http://localhost:8080/events \
  -H "Authorization: Bearer $TOKEN"

# Expected response:
# [
#   {
#     "uuid": "...",
#     "timestamp": 1678886400,
#     "userUuid": "user123",
#     "itemUuid": "item456",
#     "action": "create",
#     "payload": "{}"
#   }
# ]
```

**Success Criteria**:
- HTTP 200 status code
- Response contains event array
- No authentication errors

### Scenario 3: Access Protected Endpoint without Token
**Given** a user attempts to access protected endpoints
**When** they make a request without an Authorization header
**Then** they receive a 401 Unauthorized response

```bash
# Try to access without token
curl -X GET http://localhost:8080/events

# Expected response:
# {
#   "error": "Authorization header required"
# }
```

**Success Criteria**:
- HTTP 401 status code
- Response contains error message
- Request is rejected

### Scenario 4: Access with Invalid Token
**Given** a user has an invalid JWT token
**When** they access protected endpoints
**Then** they receive a 401 Unauthorized response

```bash
# Use invalid token
curl -X GET http://localhost:8080/events \
  -H "Authorization: Bearer invalid-token"

# Expected response:
# {
#   "error": "Invalid token"
# }
```

**Success Criteria**:
- HTTP 401 status code
- Response contains appropriate error message
- Invalid token is rejected

### Scenario 5: Create Events with Authentication
**Given** a user has a valid JWT token
**When** they POST new events
**Then** the events are created successfully

```bash
# Create new events
curl -X POST http://localhost:8080/events \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "uuid": "test-event-123",
      "timestamp": 1678886400,
      "userUuid": "test-user-456",
      "itemUuid": "test-item-789",
      "action": "create",
      "payload": "{\"name\": \"Test Item\"}"
    }
  ]'

# Expected response: Array of all events including the new one
```

**Success Criteria**:
- HTTP 200 status code
- Response contains event array
- New event is included in response

## Error Scenarios to Test

### Invalid Credentials
```bash
curl -X POST http://localhost:8080/auth/token \
  -H "Content-Type: application/json" \
  -d '{"username": "wronguser", "password": "wrongpass"}'

# Should return 401 with "Invalid username or password"
```

### Malformed Request
```bash
curl -X POST http://localhost:8080/auth/token \
  -H "Content-Type: application/json" \
  -d '{"invalid": "request"}'

# Should return 400 with "Invalid request format"
```

### Expired Token
```bash
# Use a token that has expired
curl -X GET http://localhost:8080/events \
  -H "Authorization: Bearer expired-token"

# Should return 401 with "Token has expired"
```

## Validation Checklist

- [ ] POST /auth/token returns valid JWT for correct credentials
- [ ] GET /events works with valid Bearer token
- [ ] POST /events works with valid Bearer token
- [ ] Requests without Authorization header are rejected (401)
- [ ] Requests with invalid tokens are rejected (401)
- [ ] Requests with expired tokens are rejected (401)
- [ ] Error messages are clear and appropriate
- [ ] Token contains expected user information
- [ ] Authentication integrates properly with ACL system

## Troubleshooting

### Common Issues
1. **"Authorization header required"**: Make sure to include `Authorization: Bearer <token>` header
2. **"Invalid token"**: Check that the token is properly formatted and signed
3. **"Token has expired"**: Tokens expire after 24 hours by default
4. **Connection refused**: Make sure the server is running on localhost:8080

### Debug Commands
```bash
# Check server logs (when running locally)
# Replace with your server startup command's log output

# Verify JWT token structure (without exposing secret)
echo "your-token-here" | cut -d'.' -f2 | base64 -d

# Test token validation manually
curl -v -X GET http://localhost:8080/events \
  -H "Authorization: Bearer your-token-here"
```

## Next Steps
Once all scenarios pass:
1. Run the full contract test suite
2. Test integration with ACL system
3. Verify performance meets requirements (<100ms)
4. Update API documentation
5. Consider implementing refresh tokens for better UX