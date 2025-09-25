# Feature Specification: Update Auth System

**Feature Branch**: `004-update-auth-system`  
**Created**: 2025-09-25  
**Status**: Draft  
**Input**: User description: "the new feature spec will be in 004-update-auth-system. The new feature should follow all the requirements given in issue #35."
**Note**: No existing users - no migration or backwards compatibility concerns.

## Execution Flow (main)
```
1. Parse user description from Input
   → If empty: ERROR "No feature description provided"
2. Extract key concepts from description
   → Identify: actors, actions, data, constraints
3. For each unclear aspect:
   → Mark with [NEEDS CLARIFICATION: specific question]
4. Fill User Scenarios & Testing section
   → If no clear user flow: ERROR "Cannot determine user scenarios"
5. Generate Functional Requirements
   → Each requirement must be testable
   → Mark ambiguous requirements
6. Identify Key Entities (if data involved)
7. Run Review Checklist
   → If any [NEEDS CLARIFICATION]: WARN "Spec has uncertainties"
   → If implementation details found: ERROR "Remove tech details"
8. Return: SUCCESS (spec ready for planning)
```

---

## ⚡ Quick Guidelines
- ✅ Focus on WHAT users need and WHY
- ❌ Avoid HOW to implement (no tech stack, APIs, code structure)
- 👥 Written for business stakeholders, not developers

### Section Requirements
- **Mandatory sections**: Must be completed for every feature
- **Optional sections**: Include only when relevant to the feature
- When a section doesn't apply, remove it entirely (don't leave as "N/A")

### For AI Generation
When creating this spec from a user prompt:
1. **Mark all ambiguities**: Use [NEEDS CLARIFICATION: specific question] for any assumption you'd need to make
2. **Don't guess**: If the prompt doesn't specify something (e.g., "login system" without auth method), mark it
3. **Think like a tester**: Every vague requirement should fail the "testable and unambiguous" checklist item
4. **Common underspecified areas**:
   - User types and permissions
   - Data retention/deletion policies  
   - Performance targets and scale
   - Error handling behaviors
   - Integration requirements
   - Security/compliance needs

---

## User Scenarios & Testing *(mandatory)*

### Primary User Story
As a user of the simple-sync system, I want to authenticate using long-lived API keys instead of passwords and expiring tokens, so that I can have simpler authentication management and better offline support.

### Acceptance Scenarios
1. **Given** a new user is created, **When** the user account exists, **Then** it has no API key or setup token until explicitly reset.
2. **Given** a user needs initial access, **When** an authorized user calls generateToken or resetKey, **Then** the system generates a setup token for that user.
3. **Given** a user has a setup token, **When** they exchange it for an API key without additional authorization, **Then** they receive a long-lived API key for ongoing authentication.
4. **Given** a user has an API key, **When** they make API requests, **Then** the system authenticates them without requiring password management or token refresh.
5. **Given** a user has multiple API keys, **When** they use any valid key for authentication, **Then** the system authenticates them successfully.
6. **Given** a user needs to reset their access, **When** an authorized user calls generateToken or resetKey, **Then** the system generates a new setup token for that user.
7. **Given** a user attempts to exchange an invalid setup token, **When** they call the exchange endpoint, **Then** the system returns an auth error.
8. **Given** a user attempts operations on a non-existent user, **When** they call resetKey, **Then** the system returns an auth error (preventing user enumeration).
9. **Given** a user attempts to exchange a token for a non-existent user, **When** they call token exchange, **Then** the system returns an auth error (preventing user enumeration).

### Edge Cases
- What happens when a setup token expires before use? (Return auth error)
- How does the system handle attempts to create multiple setup tokens for the same user? (Invalidate previous, allow only one valid at a time)
- What happens when an API key is compromised and needs immediate revocation?
- What happens when a user without proper ACL permissions tries to reset another user's key? (Return auth error)
- How should the system handle resetKey calls for non-existent users? (Return auth error)
- How should the system handle setup token exchange attempts for non-existent users? (Return auth error)

## Requirements *(mandatory)*

### Functional Requirements
- **FR-001**: System MUST NOT generate API keys or setup tokens automatically for new users
- **FR-002**: System MUST provide resetKey endpoint that generates setup tokens when called
- **FR-003**: System MUST restrict resetKey endpoint to users with `.user.resetKey` ACL permission for the target user ID
- **FR-004**: System MUST restrict generateToken endpoint to users with `.user.generateToken` ACL permission for the target user ID
- **FR-005**: System MUST authenticate exchangeToken requests using the setup token itself
- **FR-006**: System MUST allow `.root` user unrestricted access to all user management endpoints
- **FR-007**: System MUST allow exchange of setup tokens for API keys
- **FR-008**: System MUST support multiple API keys per user for simultaneous client authentication
- **FR-009**: System MUST authenticate users using API keys instead of passwords
- **FR-010**: System MUST maintain user identity resolution from API keys
- **FR-011**: System MUST ensure setup tokens expire after 24 hours
- **FR-012**: System MUST enforce single-use constraint on setup tokens
- **FR-013**: System MUST allow only one valid setup token per user at a time
- **FR-014**: System MUST return auth error for invalid setup token exchanges
- **FR-015**: System MUST return auth error for all operations on non-existent users (preventing user enumeration)
- **FR-016**: System MUST eliminate password management complexity

### Key Entities *(include if feature involves data)*
- **API Key**: Cryptographically random credential with sk_ prefix for user authentication, never expires, stored separately from event data
- **Setup Token**: 8-character alphanumeric token with hyphen separator (XXXX-XXXX format) for initial user setup, short-lived and single-use, expires after 24 hours
- **User**: System user with associated API key for authentication and authorization

---

## Review & Acceptance Checklist
*GATE: Automated checks run during main() execution*

### Content Quality
- [ ] No implementation details (languages, frameworks, APIs)
- [ ] Focused on user value and business needs
- [ ] Written for non-technical stakeholders
- [ ] All mandatory sections completed

### Requirement Completeness
- [ ] No [NEEDS CLARIFICATION] markers remain
- [ ] Requirements are testable and unambiguous  
- [ ] Success criteria are measurable
- [ ] Scope is clearly bounded
- [ ] Dependencies and assumptions identified

---

## Execution Status
*Updated by main() during processing*

- [ ] User description parsed
- [ ] Key concepts extracted
- [ ] Ambiguities marked
- [ ] User scenarios defined
- [ ] Requirements generated
- [ ] Entities identified
- [ ] Review checklist passed

---