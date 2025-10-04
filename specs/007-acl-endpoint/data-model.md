# Data Model: 007-acl-endpoint

## Entities

### ACL Rule (models.AclRule)
Represents an access control rule with the following attributes:

- **user** (string): The user identifier for whom the permission applies
- **item** (string): The item or resource identifier being controlled
- **action** (string): The action being permitted or denied (e.g., "read", "write")
- **type** (string): Either "allow" or "deny"
- **timestamp** (uint64): Unix timestamp set automatically by the server

**Validation Rules**:
- All attributes except timestamp are required and non-empty
- User, item, and action must be valid strings (no control characters)
- Type must be either "allow" or "deny"
- Timestamp is set server-side and cannot be overridden

**Relationships**:
- ACL Rules are converted to events and stored in the event stream
- No direct relationships to other entities

**State Transitions**:
- ACL rules are immutable once created
- New rules can be added but existing ones cannot be modified</content>
</xai:function_call/>
</xai:function_call name="list">
<parameter name="path">/home/aemig/Documents/repos/kwila/simple-sync/specs/007-acl-endpoint