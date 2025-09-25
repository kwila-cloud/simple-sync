# Quickstart: Update Auth System

## Overview
This guide demonstrates the new API key authentication system that replaces JWT tokens with long-lived API keys and short-lived setup tokens.

## Prerequisites
- Simple Sync server running
- Admin user with appropriate ACL permissions (.user.generateToken, .user.resetKey, .user.exchangeToken)
- API testing tool (curl, Postman, etc.)

## Scenario 1: New User Setup Flow

### Step 1: Create a User
First, create a user account (this would typically be done through existing user creation mechanisms):

```bash
# This step assumes you have a way to create users
# The exact endpoint depends on your current user creation process
```

### Step 2: Generate Setup Token
Generate a setup token for the new user:

```bash
curl -X POST "http://localhost:8080/api/v1/user/generateToken?user=<USER_ID>" \
  -H "Authorization: Bearer <ADMIN_API_KEY>" \
  -H "Content-Type: application/json"
```

**Expected Response:**
```json
{
  "token": "ABCD-1234",
  "expiresAt": "2025-09-26T12:00:00Z"
}
```

### Step 3: Exchange Setup Token for API Key
Exchange the setup token for a long-lived API key:

```bash
curl -X POST "http://localhost:8080/api/v1/setup/exchangeToken" \
  -H "Content-Type: application/json" \
  -d '{
    "token": "ABCD-1234",
    "description": "Desktop Client"
  }'
```

**Expected Response:**
```json
{
  "keyUuid": "550e8400-e29b-41d4-a716-446655440000",
  "apiKey": "sk_abcdefghijklmnopqrstuvwxyz1234567890",
  "user": "<USER_ID>",
  "description": "Desktop Client"
}
```

### Step 4: Use API Key for Authentication
Use the API key for subsequent requests:

```bash
curl -X GET "http://localhost:8080/api/v1/events" \
  -H "Authorization: Bearer sk_abcdefghijklmnopqrstuvwxyz1234567890"
```

## Scenario 2: Reset Existing User Access

### Reset User API Key
If a user needs to reset their access:

```bash
curl -X POST "http://localhost:8080/api/v1/user/resetKey?user=<USER_ID>" \
  -H "Authorization: Bearer <ADMIN_API_KEY>" \
  -H "Content-Type: application/json"
```

Then follow Steps 3-4 above to exchange the new setup token.

## Error Scenarios

### Insufficient Permissions
```bash
curl -X POST "http://localhost:8080/api/v1/user/generateToken?user=<USER_ID>" \
  -H "Authorization: Bearer <INSUFFICIENT_PERMISSIONS_KEY>"
```

**Expected Response:**
```json
{
  "error": "Unauthorized"
}
```

### Expired Setup Token
```bash
curl -X POST "http://localhost:8080/api/v1/setup/exchangeToken" \
  -H "Content-Type: application/json" \
  -d '{
    "token": "EXPIRED-TOKEN"
  }'
```

**Expected Response:**
```json
{
  "error": "Unauthorized"
}
```

### Non-existent User
```bash
curl -X POST "http://localhost:8080/api/v1/user/generateToken?user=non-existent-user" \
  -H "Authorization: Bearer <ADMIN_API_KEY>"
```

**Expected Response:**
```json
{
  "error": "Unauthorized"
}
```

## Key Points
- Setup tokens expire after 24 hours
- Each user can only have one active setup token
- Users can have multiple API keys simultaneously (one per client/device)
- API keys never expire and are suitable for offline use
- All operations require appropriate ACL permissions
- Error responses are consistent for security (no user enumeration)
- Non-existent user operations return auth errors to prevent information leakage</content>
</xai:function_call">**Output**: quickstart.md with step-by-step testing scenarios