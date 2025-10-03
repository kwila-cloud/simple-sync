# Feature: Fix API Key Header

## Summary
Change API authentication from `Authorization: Bearer <key>` to `X-API-Key: <key>` header for all endpoints. Update Go/Gin middleware implementation and documentation to reflect the new header requirement. Completely replace Bearer authentication with X-API-Key, rejecting any requests using the old header.

## Technical Decisions
- **Header Standard**: Use X-API-Key as the industry standard for API key authentication instead of Authorization: Bearer (which is typically for OAuth tokens)
- **Complete Replacement**: Reject Authorization: Bearer headers entirely, no backward compatibility
- **Error Handling**: Provide clear error messages for authentication failures
- **Scope**: All protected API endpoints (events, user management) require X-API-Key; health endpoint remains unprotected

## Implementation Details
- **Middleware Update**: Modify `src/middleware/auth.go` to check `X-API-Key` header instead of `Authorization: Bearer`
- **Rejection Logic**: Reject requests containing `Authorization: Bearer` header with clear error message
- **Validation**: Ensure X-API-Key header is present and non-empty for protected routes
- **No Data Model Changes**: API Key entity remains unchanged

## API Changes
- **Authentication Header**: `X-API-Key: <api-key>` required for all protected endpoints
- **Rejected Header**: `Authorization: Bearer <key>` now returns 401 Unauthorized
- **Endpoints Affected**:
  - POST /api/v1/events
  - GET /api/v1/events
  - POST /api/v1/user/generateToken
  - POST /api/v1/user/resetKey
- **Unprotected Endpoints**: GET /health, POST /api/v1/setup/exchangeToken

## Testing Approach
- **Contract Tests**: Updated all existing tests to use X-API-Key header
- **New Tests**: Added tests for X-API-Key acceptance and Bearer rejection
- **Integration Tests**: End-to-end testing of API key generation and authentication flow
- **TDD Process**: Tests written first, then implementation updated to pass them

## Documentation Updates
- **API Docs**: Update `docs/src/content/docs/api/v1.md` with X-API-Key requirement
- **Examples**: Update curl commands in AGENTS.md and README.md
- **Quickstart**: Updated user guide with new header usage

## Acceptance Criteria
1. API key generation works with X-API-Key header
2. X-API-Key header authentication succeeds for protected endpoints
3. Authorization: Bearer header is rejected with clear error
4. All existing functionality preserved with new header
5. Documentation reflects new authentication method

## Implementation Status
- [x] Auth middleware updated
- [x] Contract tests updated
- [x] Integration tests updated
- [x] Documentation updated
- [x] All tests passing