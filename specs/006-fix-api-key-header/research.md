# Research: Fix API Key Header

## Decisions Made

### Decision: Use X-API-Key Header for Authentication
**Rationale**: X-API-Key is the standard header for API key authentication, providing better convention and compatibility with API clients compared to Authorization: Bearer (which is typically for OAuth tokens).

**Alternatives Considered**:
- Keep Authorization: Bearer - rejected due to user requirement for change
- Support both headers - rejected to avoid confusion and maintenance burden
- Custom header name - rejected as X-API-Key is the industry standard

### Decision: Completely Replace Bearer Header
**Rationale**: User specification requires "rather than" Bearer, indicating complete replacement. This simplifies implementation and avoids dual support complexity.

**Alternatives Considered**:
- Support both for backward compatibility - rejected as it contradicts the "rather than" requirement
- Deprecation period - rejected as not specified and adds unnecessary complexity

### Decision: Reject Requests with Bearer Header
**Rationale**: To enforce the new standard, requests using the old header should be rejected with clear error messages.

**Alternatives Considered**:
- Accept but log deprecation warning - rejected as it doesn't meet the requirement to use X-API-Key instead
- Redirect or transform - rejected as overly complex for a header change

## Technical Research

### Gin Framework Header Handling
- Gin provides `c.GetHeader("X-API-Key")` for retrieving headers
- Existing auth middleware in `middleware/auth.go` needs modification
- Error responses can use `c.JSON()` with appropriate status codes

### Testing Approach
- Update existing contract tests to use X-API-Key header
- Add tests for rejection of Bearer header
- Verify error messages are clear

### Documentation Updates
- Update API documentation in `docs/src/content/docs/api/`
- Update example curl commands in AGENTS.md
- Ensure README reflects new header requirement