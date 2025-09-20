# Data Model: High Performance REST API for simple-sync

## Event Entity
- **Fields**:
  - `uuid`: string (unique identifier)
  - `timestamp`: uint64 (Unix timestamp)
  - `userUuid`: string (user who performed action)
  - `itemUuid`: string (item affected)
  - `action`: string (action performed: create, update, delete)
  - `payload`: string (JSON data payload)
- **Validation Rules**:
  - uuid must be valid UUID format
  - timestamp must be positive uint64
  - userUuid, itemUuid must be non-empty strings
  - action must be one of: create, update, delete
  - payload must be valid JSON string
- **Relationships**: None (append-only event log)

## User Entity
- **Fields**:
  - `uuid`: string (unique identifier)
  - `username`: string (login username)
  - `password`: string (hashed password)
- **Validation Rules**:
  - uuid must be valid UUID format
  - username must be unique, non-empty, alphanumeric
  - password must be hashed with bcrypt
- **Relationships**: Referenced by Event.userUuid

## ACL Entry Entity
- **Fields**:
  - `userUuid`: string (user UUID)
  - `itemUuid`: string (item UUID or wildcard)
  - `permissions`: []string (array of "read", "write")
- **Validation Rules**:
  - userUuid must be valid UUID or "*"
  - itemUuid must be valid UUID, "*" or prefix pattern (e.g., "task-*")
  - permissions must contain only "read" and/or "write"
- **Relationships**: Links users to items with permissions

## State Transitions
- **Event Creation**: New events are appended to the event log
- **User Registration**: New users are added to user store
- **ACL Updates**: ACL entries can be added/modified in batch
- **Authentication**: JWT tokens issued on login, validated on requests

## Storage Patterns
- Events: Append-only JSON array in data/events.json
- Users: JSON object map in memory (for MVP)
- ACL: JSON array in data/acl.json
- All data persisted to files with atomic writes