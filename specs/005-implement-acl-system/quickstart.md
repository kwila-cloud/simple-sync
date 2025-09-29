# Quickstart: ACL System Testing

## Prerequisites
- Server running on localhost:8080
- User registered and authenticated (API key available as bearer token)

## Test Scenarios

### 1. Set ACL Rules
```bash
# Allow user to write/delete on specific item via event
curl -X POST http://localhost:8080/events \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"itemUuid":".acl","data":{"action":".acl.allow","user":"testuser","item":"item123","action":"delete"}}'
```

### 2. Test Permission Denied
```bash
# Try to create event without permission
curl -X POST http://localhost:8080/events \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"itemUuid":"restricted-item","data":{"action":"write"}}'
# Should return 403 Forbidden
```

### 3. Test Permission Granted
```bash
# Add allow rule first
curl -X POST http://localhost:8080/events \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"itemUuid":".acl","data":{"action":".acl.allow","user":"testuser","item":"allowed-item","action":"write"}}'

# Now create event
curl -X POST http://localhost:8080/events \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"itemUuid":"allowed-item","data":{"action":"write"}}'
# Should succeed with 201 Created
```

### 4. Test Root User Bypass
```bash
# Use .root API key (assuming setup)
# All operations should succeed regardless of ACL
```

### 5. View ACL Rules
```bash
curl -X GET "http://localhost:8080/events?itemUuid=.acl" \
  -H "Authorization: Bearer $API_KEY"
# Returns ACL events
```

## Validation Checklist
- [ ] ACL rules can be set via POST /events with .acl item
- [ ] Unauthorized operations return 403
- [ ] Authorized operations succeed
- [ ] Root user bypasses ACL checks
- [ ] Rules are persisted across restarts