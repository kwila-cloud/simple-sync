# Data Persistence Implementation

https://github.com/kwila-cloud/simple-sync/issues/7

Implement SQLite-based persistent storage for events, users, API keys, setup tokens, and ACL rules.

Continue to use the current in-memory TestStorage for tests, but the new SQLite-based storage in main.go.

## Design Decisions

SQLite chosen over Go marshaling for:
- Rich querying capabilities for ACL lookups and event filtering
- Better scalability with indexing and pagination
- Built-in transaction support and concurrent access
- Future growth potential for complex queries and analytics
- ACID compliance for data integrity

Encryption at rest will be addressed separately (issue #17) using SQLCipher or file-level encryption.

## Task List

### Storage Interface Updates
- [x] Add tests for new ACL storage interface methods
- [x] Add ACL-specific methods to storage interface

### ACL Service Refactoring  
- [x] Add tests expecting ACL service to use storage methods instead of events
- [x] Decouple ACL service from event-based rule loading

### Storage Factory and Error Types
- [x] Add tests for storage factory and database-specific errors
- [x] Create storage factory for better testability
- [x] Add specific error types

### Model Validation Updates
- [x] Add tests for database-compatible model validation
- [x] Review model validation for database compatibility

### SQLite Storage Foundation
- [x] Add tests for SQLite storage initialization and connection management
- [x] Create SQLite storage implementation with database connection management
- [x] Add database initialization methods (Initialize, Close)

### Database Schema and Migrations
- [x] Add tests for database schema creation and migrations
- [x] Design and implement database schema (events, users, api_keys, setup_tokens, acl_rules tables)
- [x] Implement table creation with migrations, indexes, and constraints

### Event Storage Implementation
- [ ] Add tests for event storage operations
- [ ] Implement SaveEvents
- [ ] Implement LoadEvents

### User Storage Implementation
- [ ] Add tests for user storage operations
- [ ] Implement SaveUser
- [ ] Implmement GetUserById

### ACL Rule Storage Implementation
- [ ] Add tests for ACL rule storage operations
- [ ] Implement CreateAclRule
- [ ] Implement GetAclRules

### Setup Token and API Key Storage Implementation
- [ ] Add tests for API key storage operations
- [ ] Add tests for setup token storage operations
- [ ] Implement CreateApiKey
- [ ] Implement GetApiKeyByHash
- [ ] Implement GetAllApiKeys
- [ ] Implement UpdateApiKey
- [ ] Implement InvalidateUserApiKeys
- [ ] Implement CreateSetupToken
- [ ] Implement GetSetupToken
- [ ] Implement UpdateSetupToken
- [ ] Implement InvalidateUserSetupTokens

### Performance and Concurrency Testing
- [ ] Add performance and concurrency tests
- [ ] Validate concurrent access and large dataset handling

### Documentation and Configuration Updates
- [ ] Update Docker configuration for data persistence
- [ ] Update AGENTS.md with SQLite storage information
- [ ] Update README.md with data persistence features and setup instructions
- [ ] Update user-facing documentation in docs/ with SQLite configuration
- [ ] Document backup/restore procedures and security considerations
