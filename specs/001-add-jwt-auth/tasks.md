# Tasks: Add JWT Authentication

**Input**: Design documents from `/home/aemig/Documents/repos/kwila/simple-sync/specs/001-add-jwt-auth/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/auth-api.yaml, quickstart.md

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Extract: Go 1.25, Gin framework, golang-jwt/jwt/v5, single project structure
2. Load optional design documents:
   → data-model.md: Extract entities → User, JWT Token, Authentication Request, Token Claims
   → contracts/auth-api.yaml: Extract endpoints → POST /auth/token, GET /events (protected), POST /events (protected)
   → research.md: Extract decisions → JWT library selection, middleware patterns, stateless tokens
   → quickstart.md: Extract test scenarios → 5 integration scenarios for auth flow validation
3. Generate tasks by category:
   → Setup: JWT dependency, environment config, linting
   → Tests: Contract tests for auth endpoints, integration tests for auth flow
   → Core: User/Token models, auth service, JWT utilities, auth endpoint, middleware
   → Integration: Storage connection, event context, ACL integration
   → Polish: Unit tests, performance, docs, quickstart validation
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests? Yes (3 contract tests)
   → All entities have models? Yes (User, JWT Token models)
   → All endpoints implemented? Yes (POST /auth/token, protected /events)
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 3.1: Setup
- [ ] T001 Add golang-jwt/jwt/v5 dependency to go.mod in /home/aemig/Documents/repos/kwila/simple-sync/go.mod
- [ ] T002 Configure JWT_SECRET environment variable in /home/aemig/Documents/repos/kwila/simple-sync/main.go
- [ ] T003 [P] Configure linting and formatting tools for Go project

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T004 [P] Contract test POST /auth/token in /home/aemig/Documents/repos/kwila/simple-sync/tests/contract/test_auth_token_post.go
- [ ] T005 [P] Contract test GET /events (protected) in /home/aemig/Documents/repos/kwila/simple-sync/tests/contract/test_events_get_protected.go
- [ ] T006 [P] Contract test POST /events (protected) in /home/aemig/Documents/repos/kwila/simple-sync/tests/contract/test_events_post_protected.go
- [ ] T007 [P] Integration test successful authentication flow in /home/aemig/Documents/repos/kwila/simple-sync/tests/integration/test_auth_success.go
- [ ] T008 [P] Integration test protected endpoint access in /home/aemig/Documents/repos/kwila/simple-sync/tests/integration/test_protected_access.go
- [ ] T009 [P] Integration test invalid credentials handling in /home/aemig/Documents/repos/kwila/simple-sync/tests/integration/test_auth_invalid.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T010 [P] User model with validation in /home/aemig/Documents/repos/kwila/simple-sync/src/models/user.go
- [ ] T011 [P] JWT Token model with claims in /home/aemig/Documents/repos/kwila/simple-sync/src/models/token.go
- [ ] T012 [P] Authentication service with token generation in /home/aemig/Documents/repos/kwila/simple-sync/src/services/auth_service.go
- [ ] T013 [P] JWT utilities for token validation in /home/aemig/Documents/repos/kwila/simple-sync/src/utils/jwt.go
- [ ] T014 POST /auth/token endpoint implementation in /home/aemig/Documents/repos/kwila/simple-sync/src/handlers/auth.go
- [ ] T015 Authentication middleware for JWT validation in /home/aemig/Documents/repos/kwila/simple-sync/src/middleware/auth.go
- [ ] T016 Integrate auth middleware with existing /events endpoints in /home/aemig/Documents/repos/kwila/simple-sync/src/handlers/events.go

## Phase 3.4: Integration
- [ ] T017 Connect auth service to in-memory user storage in /home/aemig/Documents/repos/kwila/simple-sync/src/storage/memory.go
- [ ] T018 Add authenticated user context to event creation in /home/aemig/Documents/repos/kwila/simple-sync/src/handlers/events.go
- [ ] T019 Update ACL system to use authenticated user for permission checks in /home/aemig/Documents/repos/kwila/simple-sync/src/handlers/acl.go

## Phase 3.5: Polish
- [ ] T020 [P] Unit tests for JWT utilities in /home/aemig/Documents/repos/kwila/simple-sync/tests/unit/test_jwt_utils.go
- [ ] T021 [P] Unit tests for auth service in /home/aemig/Documents/repos/kwila/simple-sync/tests/unit/test_auth_service.go
- [ ] T022 Performance tests for auth endpoints (<100ms) in /home/aemig/Documents/repos/kwila/simple-sync/tests/performance/test_auth_performance.go
- [ ] T023 [P] Update docs/api.md with authentication endpoints and JWT usage
- [ ] T024 Run quickstart.md validation scenarios to verify complete implementation

## Dependencies
- Tests (T004-T009) before implementation (T010-T016)
- T010-T011 blocks T012, T017
- T012 blocks T014-T015
- T014 blocks T016, T018
- T015 blocks T016, T019
- Implementation before polish (T020-T024)

## Parallel Example
```
# Launch T004-T009 together (all test tasks are independent):
Task: "Contract test POST /auth/token in /home/aemig/Documents/repos/kwila/simple-sync/tests/contract/test_auth_token_post.go"
Task: "Contract test GET /events (protected) in /home/aemig/Documents/repos/kwila/simple-sync/tests/contract/test_events_get_protected.go"
Task: "Contract test POST /events (protected) in /home/aemig/Documents/repos/kwila/simple-sync/tests/contract/test_events_post_protected.go"
Task: "Integration test successful authentication flow in /home/aemig/Documents/repos/kwila/simple-sync/tests/integration/test_auth_success.go"
Task: "Integration test protected endpoint access in /home/aemig/Documents/repos/kwila/simple-sync/tests/integration/test_protected_access.go"
Task: "Integration test invalid credentials handling in /home/aemig/Documents/repos/kwila/simple-sync/tests/integration/test_auth_invalid.go"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task
- Avoid: vague tasks, same file conflicts

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   - Each contract file → contract test task [P]
   - Each endpoint → implementation task
   
2. **From Data Model**:
   - Each entity → model creation task [P]
   - Relationships → service layer tasks
   
3. **From User Stories**:
   - Each story → integration test [P]
   - Quickstart scenarios → validation tasks

4. **Ordering**:
   - Setup → Tests → Models → Services → Endpoints → Polish
   - Dependencies block parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [ ] All contracts have corresponding tests
- [ ] All entities have model tasks
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task