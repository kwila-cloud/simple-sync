# Data Model: Add JWT Authentication

**Date**: Sat Sep 20 2025
**Feature**: 001-add-jwt-auth

## Overview
This feature adds JWT-based authentication entities to support secure access to the simple-sync API. The data model focuses on authentication-related entities while integrating with the existing event-driven architecture. For the MVP, authentication data will be managed in-memory with future plans for persistent storage.

## Entities

### User
Represents an authenticated user in the system.

**Fields**:
- `uuid` (string): Unique identifier for the user
- `username` (string): Login username (unique, required)
- `password_hash` (string): Hashed password for authentication
- `created_at` (timestamp): When the user was created
- `is_admin` (boolean): Whether user has administrative privileges

**Validation Rules**:
- Username: 3-50 characters, alphanumeric + underscore/hyphen
- Password: Minimum 8 characters (enforced at creation)
- UUID: Valid UUID format

**Relationships**:
- One-to-many with JWT Token (user can have multiple active tokens)
- Referenced by ACL rules for permission evaluation

### JWT Token
Represents an authentication token issued to a user.

**Fields**:
- `token_string` (string): The actual JWT token
- `user_uuid` (string): Reference to the user who owns this token
- `issued_at` (timestamp): When the token was created
- `expires_at` (timestamp): When the token expires
- `is_revoked` (boolean): Whether this token has been revoked (future feature)

**Validation Rules**:
- Token string: Valid JWT format
- Expiration: Must be in the future when created
- User UUID: Must reference existing user

**Relationships**:
- Many-to-one with User
- Used by middleware for authentication validation

### Authentication Request
Represents a login attempt (not persisted, used for validation).

**Fields**:
- `username` (string): Provided username
- `password` (string): Provided password (plaintext, validated then discarded)

**Validation Rules**:
- Username: Required, non-empty
- Password: Required, non-empty

### Token Claims
Represents the payload embedded in JWT tokens.

**Fields**:
- `user_uuid` (string): The authenticated user's UUID
- `username` (string): The authenticated user's username
- `issued_at` (timestamp): Token issuance time
- `expires_at` (timestamp): Token expiration time
- `is_admin` (boolean): Whether user has admin privileges

**Validation Rules**:
- All timestamps must be valid
- User UUID must reference existing user
- Expiration must be after issuance

## State Transitions

### User States
- `active`: Normal authenticated user
- `inactive`: User account disabled (future feature)
- `admin`: User with administrative privileges

### Token States
- `active`: Valid and usable token
- `expired`: Token past expiration date
- `revoked`: Manually invalidated token

## Data Integrity Rules

1. **Token Expiration**: System must reject expired tokens and clean them up periodically
2. **User Existence**: All tokens must reference valid users (enforced at runtime)
3. **Single Active Session**: MVP allows multiple tokens per user
4. **Password Security**: Passwords must be properly hashed using secure algorithms
5. **Token Uniqueness**: Each token string must be unique (enforced by JWT library)
6. **Memory Management**: Expired tokens must be cleaned up to prevent memory leaks

## Integration with Existing Model

### Event Entity Integration
- Authentication adds `user_uuid` context to events
- ACL evaluation uses authenticated user for permission checks
- Event creation requires valid authentication

### ACL Integration
- User entity provides context for ACL rule evaluation
- Admin status affects ACL permissions
- Authentication middleware provides user context to ACL system

## Storage Considerations

### Current Implementation (In-Memory)
For the MVP, authentication data will be stored in memory using Go data structures:

- **Users**: Stored in a map with username as key for fast lookup
- **Active Tokens**: Stored in a map with token hash as key for validation
- **Token Cleanup**: Periodic cleanup of expired tokens to prevent memory leaks

### Future Persistent Storage
When SQLite integration is implemented (future issue), the following schema will be used:

```sql
-- Users table (future)
CREATE TABLE users (
    uuid TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_admin BOOLEAN DEFAULT FALSE
);

-- Active tokens table (future, for revocation)
CREATE TABLE active_tokens (
    token_hash TEXT PRIMARY KEY,
    user_uuid TEXT NOT NULL,
    issued_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid)
);

-- Indexes for performance (future)
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_active_tokens_user ON active_tokens(user_uuid);
CREATE INDEX idx_active_tokens_expires ON active_tokens(expires_at);
```

### Implementation Strategy
1. Implement in-memory storage for MVP with hardcoded users
2. Add token management with automatic expiration cleanup
3. Design interfaces to support future SQLite integration
4. Update existing event handlers to use authentication context
5. Add authentication middleware to protected routes

## Security Considerations

1. **Password Hashing**: Use bcrypt with appropriate cost factor for secure password storage
2. **Token Security**: Use strong secret keys from environment variables, validate all claims
3. **Session Management**: Implement proper token expiration with periodic cleanup
4. **Memory Security**: Ensure sensitive data is properly managed in memory
5. **Audit Logging**: Log authentication events for security monitoring (future enhancement)
6. **Future Database Security**: When SQLite is implemented, ensure proper encryption and access controls

## Future Extensions

1. **Refresh Tokens**: Add refresh token support for better UX
2. **Multi-factor Authentication**: Extend user model for MFA
3. **OAuth Integration**: Support external authentication providers
4. **Session Management**: Add session tracking and management
5. **Password Reset**: Implement secure password reset flow