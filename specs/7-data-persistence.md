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
- [ ] Add tests for database-compatible model validation
- [ ] Review model validation for database compatibility

### SQLite Storage Foundation
- [ ] Add tests for SQLite storage initialization and connection management
- [ ] Create SQLite storage implementation with database connection management
- [ ] Add database initialization methods (Initialize, Close)

### Database Schema and Migrations
- [ ] Add tests for database schema creation and migrations
- [ ] Design and implement database schema (events, users, api_keys, setup_tokens, acl_rules tables)
- [ ] Implement table creation with migrations, indexes, and constraints

### Event Storage Implementation
- [ ] Add tests for event storage operations
- [ ] Implement event storage with transaction support and pagination

### User Storage Implementation
- [ ] Add tests for user storage operations
- [ ] Implement user storage with uniqueness constraints

### API Key Storage Implementation
- [ ] Add tests for API key storage operations
- [ ] Implement API key storage with security considerations

### Setup Token Storage Implementation
- [ ] Add tests for setup token storage operations
- [ ] Implement setup token storage with expiration handling

### ACL Rule Storage Implementation
- [ ] Add tests for ACL rule storage operations
- [ ] Implement ACL rule storage with filtering and indexing

### Service Integration
- [ ] Add tests for updated ACL and authentication services
- [ ] Update ACL and authentication services to use SQLite storage

### Main Application Integration
- [ ] Add tests for main application SQLite integration
- [ ] Update main.go to use SQLite storage with database configuration
- [ ] Add database connection management with pooling and health checks

### Docker Configuration
- [ ] Add tests for Docker data persistence
- [ ] Update Docker configuration for data persistence

### Performance and Concurrency Testing
- [ ] Add performance and concurrency tests
- [ ] Validate concurrent access and large dataset handling
- [ ] Test database migrations and schema updates

### Documentation Updates
- [ ] Update AGENTS.md with SQLite storage information
- [ ] Update README.md with data persistence features and setup instructions

### User Documentation
- [ ] Update user-facing documentation in docs/ with SQLite configuration
- [ ] Document backup/restore procedures and security considerations
