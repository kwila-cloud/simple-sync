# Research: Update Auth System

## Phase 0 Research Findings

### API Key Generation
**Decision**: Use cryptographically secure random bytes (32 bytes) encoded as base64, prefixed with "sk_"
**Rationale**: Provides sufficient entropy for security while being URL-safe and human-readable
**Alternatives considered**: UUIDs (less secure), sequential IDs (predictable), custom algorithms (complexity)

### Setup Token Format
**Decision**: 8-character alphanumeric with hyphen separator (XXXX-XXXX)
**Rationale**: Human-readable format that's easy to type and share, while providing sufficient randomness
**Alternatives considered**: Pure numeric (harder to read), longer tokens (more secure but cumbersome), QR-only (not accessible)

### Token Expiration
**Decision**: 24-hour expiration for setup tokens
**Rationale**: Balances security (short window) with usability (enough time for setup)
**Alternatives considered**: 1 hour (too short), 7 days (too long), no expiration (insecure)

### API Key Storage
**Decision**: Store API keys encrypted in database using AES-256-GCM, allowing retrieval for exchangeToken responses
**Rationale**: Protects against database compromise while enabling API key distribution to users
**Alternatives considered**: Plain text storage (insecure), hash-only storage (prevents key retrieval), external key store (adds complexity)

### Single Token Constraint
**Decision**: Invalidate previous setup tokens when generating new ones for same user
**Rationale**: Prevents token accumulation and confusion, ensures only one active setup path per user
**Alternatives considered**: Allow multiple tokens (confusion), token reuse (security risk)

### ACL Permission Design
**Decision**: Separate permissions for generateToken (.user.generateToken), exchangeToken (.user.exchangeToken), resetKey (.user.resetKey)
**Rationale**: Provides fine-grained access control for different user management operations
**Alternatives considered**: Single permission for all operations (too broad), no permissions (insecure)

### Error Handling
**Decision**: Return consistent auth error responses for all permission/token failures
**Rationale**: Prevents information leakage about system state while maintaining consistent API behavior
**Alternatives considered**: Specific error codes (information leakage), success responses (confusing)

### Migration Strategy
**Decision**: No migration needed - fresh implementation for new system
**Rationale**: No existing users or JWT tokens to migrate from
**Alternatives considered**: Gradual migration (not applicable), dual auth support (unnecessary complexity)</content>
</xai:function_call">**Output**: research.md with all research findings documented