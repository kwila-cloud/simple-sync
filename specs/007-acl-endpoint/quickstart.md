# Quickstart: 007-acl-endpoint

## Prerequisites
- Simple-sync server running on localhost:8080
- Valid API key for a user with ACL permissions
- Admin API key for generating setup tokens (if needed)

## Test the Dedicated ACL Endpoint

1. **Submit ACL Event via Dedicated Endpoint**
   ```bash
   curl -X POST http://localhost:8080/api/v1/acl \
     -H "X-API-Key: YOUR_API_KEY" \
     -H "Content-Type: application/json" \
     -d '[{"user":"testuser","item":"testitem","action":"read"}]'
   ```
   Expected: 200 OK with success message

2. **Verify ACL Event Was Stored**
   ```bash
   curl -X GET "http://localhost:8080/api/v1/events?itemUuid=.acl" \
     -H "X-API-Key: YOUR_API_KEY"
   ```
   Expected: JSON array containing the ACL event with current timestamp

3. **Attempt to Submit ACL Event via /events (Should Fail)**
   ```bash
   curl -X POST http://localhost:8080/api/v1/events \
     -H "X-API-Key: YOUR_API_KEY" \
     -H "Content-Type: application/json" \
     -d '[{"uuid":"acl-test","timestamp":1640995200,"user":"testuser","item":".acl","action":".acl.allow","payload":"{\"user\":\"testuser\",\"item\":\"testitem\",\"action\":\"read\"}"}]'
   ```
   Expected: 400 Bad Request with error message

4. **Test Invalid ACL Data**
   ```bash
   curl -X POST http://localhost:8080/api/v1/acl \
     -H "X-API-Key: YOUR_API_KEY" \
     -H "Content-Type: application/json" \
     -d '[{"user":"","item":"testitem","action":"read"}]'
   ```
   Expected: 400 Bad Request with validation error

## Validation Checklist
- [ ] Dedicated endpoint accepts valid ACL events
- [ ] Events are stored with current timestamp
- [ ] /events endpoint rejects ACL events
- [ ] Invalid data returns appropriate errors
- [ ] Authentication is enforced</content>
</xai:function_call/>
</xai:function_call name="bash">
<parameter name="command">cd /home/aemig/Documents/repos/kwila/simple-sync && .specify/scripts/bash/update-agent-context.sh opencode