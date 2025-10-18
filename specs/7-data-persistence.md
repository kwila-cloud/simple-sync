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

### Pre-Implementation Refactoring
- [ ] Add ACL-specific methods to storage interface
  - [ ] CreateAclRule(rule *models.AclRule) error
  - [ ] GetAclRules() ([]models.AclRule, error)
  - [ ] GetAclRulesByUser(user string) ([]models.AclRule, error)
  - [ ] UpdateAclRule(rule *models.AclRule) error
  - [ ] DeleteAclRule(user, item, action string) error
- [ ] Decouple ACL service from event-based rule loading
  - [ ] Modify AclService to use storage ACL methods instead of parsing events
  - [ ] Update loadRules method to use GetAclRules()
  - [ ] Maintain backward compatibility during transition
- [ ] Create storage factory for better testability
  - [ ] Add storage factory function for environment-based storage selection
  - [ ] Update tests to use storage factory instead of direct TestStorage instantiation
  - [ ] Add database-specific error types
- [ ] Review model validation for database compatibility
  - [ ] Separate validation concerns for in-memory vs database storage
  - [ ] Update model validation to handle database-generated IDs
  - [ ] Add database-specific model considerations

### Database Implementation
- [ ] Add database initialization method
  - [ ] Initialize() error
  - [ ] Close() error
- [ ] Create SQLite storage implementation
  - [ ] Create sqlite.go file in storage package
  - [ ] Implement storage.Storage interface
  - [ ] Add database connection management
- [ ] Design database schema
  - [ ] Create events table
  - [ ] Create users table  
  - [ ] Create api_keys table
  - [ ] Create setup_tokens table
  - [ ] Create acl_rules table
- [ ] Implement table creation and migration
  - [ ] Create schema migration system
  - [ ] Add indexes for performance
  - [ ] Add foreign key constraints
- [ ] Implement event storage
  - [ ] SaveEvents method with transaction support
  - [ ] LoadEvents with pagination support
  - [ ] Add event indexing by timestamp and user
- [ ] Implement user storage
  - [ ] SaveUser with duplicate prevention
  - [ ] GetUserById with proper error handling
  - [ ] Add user uniqueness constraints
- [ ] Implement API key storage
  - [ ] CreateAPIKey with hash generation
  - [ ] GetAPIKeyByHash with security considerations
  - [ ] GetAllAPIKeys with filtering options
  - [ ] UpdateAPIKey with last_used_at tracking
  - [ ] InvalidateUserAPIKeys with cascade delete
- [ ] Implement setup token storage
  - [ ] CreateSetupToken with expiration
  - [ ] GetSetupToken with validation
  - [ ] UpdateSetupToken with usage tracking
  - [ ] InvalidateUserSetupTokens with batch operations
- [ ] Implement ACL rule storage
  - [ ] CreateAclRule with validation
  - [ ] GetAclRules with filtering
  - [ ] GetAclRulesByUser for user-specific rules
  - [ ] UpdateAclRule with conflict resolution
  - [ ] DeleteAclRule with cascade effects
- [ ] Add ACL-specific indexes
  - [ ] Index on user, item, action combination
  - [ ] Index on user for quick lookups
  - [ ] Index on item for item-based permissions

### Integration and Configuration
- [ ] Update ACL service
  - [ ] Modify loadRules to use ACL storage instead of events
  - [ ] Update AddRule to use dedicated storage
- [ ] Update authentication service
  - [ ] Ensure API key lookup uses SQLite
  - [ ] Add connection pooling for performance
- [ ] Add database connection management
  - [ ] Connection pooling configuration
  - [ ] Proper connection cleanup on shutdown
  - [ ] Health check for database connectivity
- [ ] Update main.go
  - [ ] Replace TestStorage with SQLite storage
  - [ ] Add database path configuration
  - [ ] Add graceful shutdown handling
- [ ] Add environment configuration
  - [ ] Database file path setting
  - [ ] Connection pool settings
  - [ ] Migration auto-run option
- [ ] Update Docker configuration
  - [ ] Add data volume mounting
  - [ ] Database file persistence
  - [ ] Backup strategy documentation

### Testing and Validation
- [ ] Create SQLite storage tests
  - [ ] Unit tests for all storage methods
  - [ ] Transaction rollback tests
  - [ ] Concurrency tests
- [ ] Update integration tests
  - [ ] Test ACL rule persistence
  - [ ] Test user data persistence
  - [ ] Test API key persistence
- [ ] Performance testing
  - [ ] Large dataset handling
  - [ ] Concurrent access testing
  - [ ] Query performance validation

### Documentation and Security
- [ ] Update documentation
  - [ ] Update AGENTS.md storage section
  - [ ] Add SQLite setup instructions
  - [ ] Document backup/restore procedures
- [ ] Security hardening
  - [ ] Database file permissions
  - [ ] Connection security
  - [ ] Input validation for SQL injection prevention
