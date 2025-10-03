# Tasks: Fix API Key Header

**Input**: Design documents from `/specs/006-fix-api-key-header/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- **Web app**: `backend/src/`, `frontend/src/`
- **Mobile**: `api/src/`, `ios/src/` or `android/src/`
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 3.1: Setup
- [ ] T001 Configure Go linting and formatting tools (gofmt, go vet, golint)

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T002 [P] Create contract test for X-API-Key header acceptance in tests/contract/test_x_api_key_header.go
- [ ] T003 [P] Create contract test for Bearer header rejection in tests/contract/test_bearer_rejection.go
- [ ] T004 [P] Update existing auth contract test in tests/contract/auth_token_post_test.go to use X-API-Key header
- [ ] T005 [P] Update existing events contract tests in tests/contract/events_get_protected_test.go and tests/contract/events_post_protected_test.go to use X-API-Key header
- [ ] T006 [P] Update existing events contract test in tests/contract/get_events_test.go to use X-API-Key header
- [ ] T007 [P] Update existing events contract test in tests/contract/post_events_test.go to use X-API-Key header
- [ ] T008 [P] Integration test for API key generation flow in tests/integration/test_api_key_generation.go
- [ ] T009 [P] Integration test for X-API-Key authentication on protected endpoints in tests/integration/test_x_api_key_protected.go
- [ ] T010 [P] Integration test for Bearer header rejection on protected endpoints in tests/integration/test_bearer_rejection_protected.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T011 Modify auth middleware in src/middleware/auth.go to check X-API-Key header and reject Authorization: Bearer

## Phase 3.4: Integration
- [ ] T012 Update existing integration tests in tests/integration/ to use X-API-Key header instead of Bearer

## Phase 3.5: Polish
- [ ] T013 [P] Update API documentation in docs/src/content/docs/api/v1.md to specify X-API-Key header requirement
- [ ] T014 [P] Update example curl commands in AGENTS.md to use X-API-Key header
- [ ] T015 [P] Update README.md examples to use X-API-Key header instead of Authorization: Bearer
- [ ] T016 Run all tests to ensure implementation works and all tests pass

## Dependencies
- Tests (T002-T010) before implementation (T011)
- T011 blocks T012
- Implementation before polish (T013-T015)

## Parallel Example
```
# Launch T002-T007 together (contract test updates):
Task: "Create contract test for X-API-Key header acceptance in tests/contract/test_x_api_key_header.go"
Task: "Create contract test for Bearer header rejection in tests/contract/test_bearer_rejection.go"
Task: "Update existing auth contract test in tests/contract/auth_token_post_test.go to use X-API-Key header"
Task: "Update existing events contract tests in tests/contract/events_get_protected_test.go and tests/contract/events_post_protected_test.go to use X-API-Key header"
Task: "Update existing events contract test in tests/contract/get_events_test.go to use X-API-Key header"
Task: "Update existing events contract test in tests/contract/post_events_test.go to use X-API-Key header"

# Launch T008-T010 together (integration tests):
Task: "Integration test for API key generation flow in tests/integration/test_api_key_generation.go"
Task: "Integration test for X-API-Key authentication on protected endpoints in tests/integration/test_x_api_key_protected.go"
Task: "Integration test for Bearer header rejection on protected endpoints in tests/integration/test_bearer_rejection_protected.go"
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