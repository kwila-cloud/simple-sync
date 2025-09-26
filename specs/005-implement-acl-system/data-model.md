# Data Model: ACL System

## Entities

### ACL Rule
Represents a permission rule for access control.

**Fields**:
- user: string (supports wildcards *, prefix.*)
- item: string (supports wildcards *, prefix.*)
- action: string (supports wildcards *, prefix.*)
- type: string ("allow" or "deny")
- timestamp: int64 (Unix timestamp)

**Validation**:
- user, item, action cannot be empty
- type must be "allow" or "deny"
- timestamp must be valid Unix time

**Relationships**:
- Stored as Event with item=".acl" and action=".acl.allow" or ".acl.deny"
- Data field contains the rule details as JSON

### Event (Extended)
Existing Event entity extended to support ACL rules.

**Additional Logic**:
- When item == ".acl", treat as ACL rule
- Validate ACL rules against current permissions before storage
- Filter events based on ACL during retrieval

## State Transitions

### ACL Rule Lifecycle
1. Created: New rule submitted via PUT /acl
2. Validated: Checked against current ACL (unless .root)
3. Stored: Saved as event if validation passes
4. Applied: Used in permission checks for subsequent operations

### Permission Evaluation
1. Collect all applicable rules (matching user/item/action)
2. Sort by specificity (item > user > action > timestamp)
3. For same specificity, use latest timestamp
4. Apply first matching rule (allow/deny)