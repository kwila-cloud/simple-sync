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

## Implementation Checklist

### PR 1: Storage Interface Updates
- [ ] Add ACL-specific methods to storage interface

### PR 2: ACL Service Refactoring  
- [ ] Decouple ACL service from event-based rule loading

### PR 3: Storage Factory and Error Types
- [ ] Create storage factory for better testability
- [ ] Add database-specific error types

### PR 4: Model Validation Updates
- [ ] Review model validation for database compatibility

### PR 5: SQLite Storage Foundation
- [ ] Create SQLite storage implementation with database connection management
- [ ] Add database initialization methods (Initialize, Close)

### PR 6: Database Schema and Migrations
- [ ] Design and implement database schema (events, users, api_keys, setup_tokens, acl_rules tables)
- [ ] Implement table creation with migrations, indexes, and constraints

### PR 7: Event Storage Implementation
- [ ] Implement event storage with transaction support and pagination

### PR 8: User Storage Implementation
- [ ] Implement user storage with uniqueness constraints

### PR 9: API Key Storage Implementation
- [ ] Implement API key storage with security considerations

### PR 10: Setup Token Storage Implementation
- [ ] Implement setup token storage with expiration handling

### PR 11: ACL Rule Storage Implementation
- [ ] Implement ACL rule storage with filtering and indexing

### PR 12: Service Integration
- [ ] Update ACL and authentication services to use SQLite storage

### PR 13: Main Application Integration
- [ ] Update main.go to use SQLite storage with database configuration
- [ ] Add database connection management with pooling and health checks

### PR 14: Docker Configuration
- [ ] Update Docker configuration for data persistence

### PR 15: Integration Testing
- [ ] Create integration tests for SQLite storage operations
- [ ] Test data persistence for all entities (events, users, API keys, ACLs)

### PR 16: Performance and Concurrency Testing
- [ ] Validate concurrent access and large dataset handling
- [ ] Test database migrations and schema updates

### PR 17: Documentation Updates
- [ ] Update AGENTS.md with SQLite storage information
- [ ] Update README.md with data persistence features and setup instructions

### PR 18: User Documentation
- [ ] Update user-facing documentation in docs/ with SQLite configuration
- [ ] Document backup/restore procedures and security considerations
