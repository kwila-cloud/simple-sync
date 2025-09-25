# Tasks: Update Auth System

**Input**: Design documents from `/specs/004-update-auth-system/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/
**Total Tasks**: 26 numbered tasks (T001-T026)

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
- Paths assume Go project structure with existing codebase

## Phase 3.1: Setup
- [ ] T001 Initialize auth-specific dependencies and imports
- [ ] T002 [P] Configure crypto libraries for API key encryption

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T003 [P] Contract test POST /api/v1/user/resetKey in tests/contract/auth_token_post_test.go
- [ ] T004 [P] Contract test POST /api/v1/user/generateToken in tests/contract/auth_token_post_test.go
- [ ] T005 [P] Contract test POST /api/v1/setup/exchangeToken in tests/contract/auth_token_post_test.go
- [ ] T006 [P] Integration test user setup flow in tests/integration/auth_setup_test.go
- [ ] T007 [P] Integration test error scenarios in tests/integration/auth_errors_test.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T008 [P] API Key model in src/models/api_key.go
- [ ] T009 [P] Setup Token model in src/models/setup_token.go
- [ ] T010 Update auth middleware for API key validation in src/middleware/auth.go
- [ ] T011 POST /api/v1/user/resetKey endpoint in src/handlers/auth_handlers.go
- [ ] T012 POST /api/v1/user/generateToken endpoint in src/handlers/auth_handlers.go
- [ ] T013 POST /api/v1/setup/exchangeToken endpoint in src/handlers/auth_handlers.go
- [ ] T014 ACL permission validation for auth endpoints
- [ ] T015 Remove JWT authentication middleware and related code
- [ ] T016 Remove username/password authentication endpoints
- [ ] T017 Remove password hashing and user credential models

## Phase 3.4: Integration
- [ ] T018 Database schema updates for API keys and setup tokens
- [ ] T019 API key encryption/decryption service
- [ ] T020 Setup token generation and validation service
- [ ] T021 Update existing auth service for API key support

## Phase 3.5: Polish
- [ ] T022 [P] Unit tests for API key encryption in tests/unit/auth_encryption_test.go
- [ ] T023 [P] Unit tests for token validation in tests/unit/auth_token_test.go
- [ ] T024 Performance tests for auth operations
- [ ] T025 Update API documentation
- [ ] T026 Security audit and cleanup

## Dependencies
- Tests (T003-T007) before implementation (T008-T017)
- T008, T009 block T018-T020
- T010 blocks T011-T014
- T019, T020 block T011-T013
- JWT removal (T015-T017) can be parallel with new implementation
- Implementation before polish (T022-T026)

## Parallel Example
```
# Launch T003-T007 together:
Task: "Contract test POST /api/v1/user/resetKey in tests/contract/auth_token_post_test.go"
Task: "Contract test POST /api/v1/user/generateToken in tests/contract/auth_token_post_test.go"
Task: "Contract test POST /api/v1/setup/exchangeToken in tests/contract/auth_token_post_test.go"
Task: "Integration test user setup flow in tests/integration/auth_setup_test.go"
Task: "Integration test error scenarios in tests/integration/auth_errors_test.go"
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
- [ ] No task modifies same file as another [P] task</content>
</xai:function_call">**Output**: tasks.md with 23 numbered tasks covering setup, testing, implementation, integration, and polish phases