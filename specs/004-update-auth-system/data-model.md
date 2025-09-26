# Data Model: Update Auth System

## Entities

### API Key
**Purpose**: Long-lived authentication credential for users
**Storage**: Encrypted in database, separate from event history for security
**Fields**:
- `uuid`: string (primary key)
- `user_id`: string (references user ID)
- `encrypted_key`: string (AES-256-GCM encrypted API key)
- `key_hash`: string (bcrypt hash for authentication verification)
- `created_at`: timestamp
- `last_used_at`: timestamp (nullable)
- `description`: string (optional, for user to identify keys)
**Validation**:
- `uuid`: required, valid UUID format
- `user_id`: required, valid user ID format
- `encrypted_key`: required, AES-256-GCM encrypted data
- `key_hash`: required, bcrypt hash format
- `created_at`: required, valid timestamp
**Relationships**:
- Many-to-one with User (each user can have multiple API keys for different clients)
**Constraints**:
- Multiple API keys allowed per user
- Keys never expire
- Each key must be unique

### Setup Token
**Purpose**: Short-lived token for initial user authentication setup
**Storage**: In-memory or database with automatic cleanup
**Fields**:
- `token`: string (XXXX-XXXX format)
- `user_id`: string (references user ID)
- `expires_at`: timestamp
- `used`: boolean (default false)
**Validation**:
- `token`: required, 8-char alphanumeric with hyphen
- `user_id`: required, valid user ID format
- `expires_at`: required, future timestamp
- `used`: required, boolean
**Relationships**:
- One-to-one with User (each user has at most one active setup token)
**Constraints**:
- Single use only
- Expires after 24 hours
- Only one active token per user

### User
**Purpose**: System user identity
**Storage**: Existing user storage (unchanged)
**Fields**: (existing fields maintained)
- `uuid`: string (primary key)
- `username`: string
- `created_at`: timestamp
**Relationships**:
- One-to-one with API Key
- One-to-one with Setup Token
**Constraints**: (existing constraints maintained)

## State Transitions

### API Key Lifecycle
1. **Non-existent** → **Created** (via successful token exchange)
2. **Active** → **Active** (additional keys can be created for same user)
3. **Active** → **Invalid** (individual key deleted or user deleted)

### Setup Token Lifecycle
1. **Non-existent** → **Generated** (via generateToken/resetKey endpoint)
2. **Active** → **Used** (successful exchange)
3. **Active** → **Expired** (time exceeded)
4. **Active** → **Invalidated** (new token generated for same user)

## Data Integrity Rules
- API keys must be unique across all users (no duplicate keys)
- Setup tokens must be unique and unpredictable
- User deletion cascades to all associated API keys and setup tokens
- Expired setup tokens are automatically cleaned up
- Each user can have multiple active API keys simultaneously</content>
</xai:function_call">**Output**: data-model.md with entity definitions and relationships