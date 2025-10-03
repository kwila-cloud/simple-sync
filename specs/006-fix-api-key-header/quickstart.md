# Quickstart: Fix API Key Header

## Overview
This feature changes the API authentication from `Authorization: Bearer <key>` to `X-API-Key: <key>` header.

## Getting Started

### 1. Obtain an API Key
First, generate a setup token for your user:
```bash
curl -X POST http://localhost:8080/api/v1/user/generateToken?user=your-username \
  -H "Authorization: Bearer sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E"
```

Exchange the setup token for an API key:
```bash
curl -X POST http://localhost:8080/api/v1/setup/exchangeToken \
  -H "Content-Type: application/json" \
  -d '{"token":"your-setup-token"}'
```

### 2. Use X-API-Key Header
Use the `X-API-Key` header instead of `Authorization: Bearer` for all authenticated requests:

```bash
# Create events
curl -X POST http://localhost:8080/api/v1/events \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '[{"uuid":"event-123","timestamp":1640995200,"user":"your-username","item":"item-123","action":"create","payload":"{}"}]'

# Get events
curl -X GET http://localhost:8080/api/v1/events \
  -H "X-API-Key: your-api-key"
```

### 3. Important Changes
- `Authorization: Bearer <key>` is no longer supported and will result in authentication errors
- All protected endpoints now require `X-API-Key: <key>`
- The `/health` endpoint remains unprotected

## Validation Steps
1. Verify API key generation works
2. Confirm X-API-Key header authentication succeeds
3. Confirm Authorization: Bearer header is rejected
4. Test all protected endpoints with new header