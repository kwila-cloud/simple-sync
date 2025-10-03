# Research: 007-acl-endpoint

## Decisions

### Dedicated ACL Endpoint Design
**Decision**: Implement POST /api/v1/acl endpoint for submitting ACL events  
**Rationale**: Follows REST principles, resource-oriented (/acl), uses POST for creation. Consistent with existing /events endpoint pattern.  
**Alternatives considered**: PUT /api/v1/acl (rejected: not creating a single resource), POST /api/v1/events with special handling (rejected: defeats purpose of dedicated endpoint).

### ACL Event Rejection in /events
**Decision**: Check incoming events in /events handler - reject if item == ".acl" and action starts with ".acl."  
**Rationale**: Prevents submission of ACL events with arbitrary timestamps. Matches existing ACL event format from internal events documentation.  
**Alternatives considered**: Accept but override timestamp (rejected: allows past timestamps), separate validation middleware (rejected: overkill for single check).

### Timestamp Handling
**Decision**: Automatically set timestamp to current time on ACL event submission  
**Rationale**: Ensures events are always current, prevents historical ACL changes. Aligns with feature requirement to avoid past timestamps.  
**Alternatives considered**: Allow client-provided timestamps (rejected: violates security requirement), use server time with offset (rejected: unnecessary complexity).

### Authentication Enforcement
**Decision**: Apply existing API key authentication middleware to the new endpoint  
**Rationale**: Consistent with other protected endpoints, uses clarified API key requirement.  
**Alternatives considered**: JWT authentication (constitution preference, but rejected: would require broader auth system changes).

### Concurrent Submission Handling
**Decision**: Rely on SQLite transactions for sequential processing  
**Rationale**: SQLite handles concurrency via locking, ensures data integrity. Simple and sufficient for expected load.  
**Alternatives considered**: Explicit queuing mechanism (rejected: over-engineering), reject concurrent requests (rejected: poor UX).

### Error Handling for Invalid Data
**Decision**: Return 400 Bad Request with descriptive error message  
**Rationale**: Standard HTTP practice, provides feedback to client. Matches clarified requirement.  
**Alternatives considered**: Silent failure (rejected: poor debugging), 422 Unprocessable Entity (rejected: 400 is more appropriate for validation).

### Documentation Updates
**Decision**: Update docs/src/content/docs/acl.mdx, docs/src/content/docs/internal-events.mdx, docs/src/content/docs/api/v1.md  
**Rationale**: Covers all mentioned documentation areas, ensures API consumers are informed.  
**Alternatives considered**: Update only API docs (rejected: misses ACL and internal event context).</content>
</xai:function_call/>
</xai:function_call name="edit">
<parameter name="filePath">/home/aemig/Documents/repos/kwila/simple-sync/specs/007-acl-endpoint/plan.md