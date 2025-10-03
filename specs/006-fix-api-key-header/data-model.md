# Data Model: Fix API Key Header

## Existing Entities

### API Key
- **Purpose**: Represents user authentication credentials for API access
- **Fields**: 
  - `key` (string): The API key value
  - `user_id` (string): Associated user identifier
  - `created_at` (timestamp): When the key was created
  - `expires_at` (timestamp, optional): Expiration time
- **Validation**: Key must be non-empty, unique
- **Relationships**: One-to-one with User entity

## Changes
No new entities required. The API Key entity remains unchanged. The change is in how the API key is transmitted in HTTP requests (from Authorization: Bearer to X-API-Key header).

## State Transitions
No state transitions affected by this change.