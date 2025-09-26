# Research: ACL System Implementation

## ACL Evaluation Algorithm

**Decision**: Implement specificity scoring with item > user > action > timestamp hierarchy. For same specificity, use latest timestamp. Exact matches have higher specificity than wildcards.

**Rationale**: Matches the clarified requirements from spec. Provides clear precedence rules for rule evaluation.

**Alternatives considered**:
- Flat permission model (rejected: too simplistic for complex rules)
- Role-based only (rejected: needs item/action granularity)

## ACL Storage in SQLite

**Decision**: Store ACL rules as events on .acl item with .acl.allow/.acl.deny actions. Use JSON data field for rule details.

**Rationale**: Leverages existing event storage system. Maintains audit trail and append-only nature.

**Alternatives considered**:
- Separate ACL table (rejected: duplicates event functionality)
- In-memory cache (rejected: needs persistence across restarts)

## Go ACL Service Pattern

**Decision**: Create acl_service.go with AclService struct implementing permission checks. Centralize all ACL logic here.

**Rationale**: Follows service pattern used in project (like auth_service.go). Allows easy testing and reuse across handlers.

**Alternatives considered**:
- Inline checks in handlers (rejected: violates DRY, hard to test)
- Middleware only (rejected: needs complex state management)

## Performance Considerations

**Decision**: Cache ACL rules in memory with periodic refresh. Target <10ms evaluation time.

**Rationale**: SQLite queries for every request would be too slow. Caching balances performance and consistency.

**Alternatives considered**:
- No caching (rejected: performance impact)
- External cache like Redis (rejected: adds complexity for small scale)

## Integration Points

**Decision**: Add ACL checks to POST /events and GET /events handlers. Bypass for .root user.

**Rationale**: Covers all event operations. Root bypass ensures administrative access.

**Alternatives considered**:
- Check all endpoints (rejected: unnecessary for read-only ops like health)
- Client-side checks (rejected: insecure)