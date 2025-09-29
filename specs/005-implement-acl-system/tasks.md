# Tasks: Implement ACL System

**Input**: Design documents from /home/aemig/Documents/repos/kwila/simple-sync/specs/005-implement-acl-system/
**Prerequisites**: plan.md (required), research.md, data-model.md, quickstart.md

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Tests: integration tests from quickstart
   → Core: models, services, handler integration
   → Integration: middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All scenarios have tests?
   → All entities have models?
   → All integrations implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: integration tests from quickstart
   → Core: models, services, handler integration
   → Integration: middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All scenarios have tests?
   → All entities have models?
   → All integrations implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [X] T001 [P] Integration test for setting ACL rules via POST /events in tests/integration/test_acl_setup.go
- [X] T002 [P] Integration test for permission denied on unauthorized event in tests/integration/test_acl_denied.go
- [X] T003 [P] Integration test for permission granted on authorized event in tests/integration/test_acl_granted.go
- [X] T004 [P] Integration test for root user bypass in tests/integration/test_acl_root_bypass.go
- [X] T005 [P] Integration test for retrieving ACL events in tests/integration/test_acl_retrieve.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [X] T006 [P] Extend Event model for ACL rules in src/models/event.go
- [X] T007 [P] Create AclService for permission evaluation in src/services/acl_service.go
- [X] T008 Integrate ACL checks into POST /events handler in src/handlers/handlers.go
- [X] T009 Integrate ACL filtering into GET /events handler in src/handlers/handlers.go

## Phase 3.4: Integration
- [X] T010 Ensure ACL service integrates with existing auth middleware
- [X] T011 Add logging for ACL decisions

## Phase 3.5: Polish
- [X] T012 [P] Unit tests for AclService in tests/unit/acl_service_test.go
- [X] T013 Performance validation for ACL evaluation
- [X] T014 [P] Update documentation for ACL system

## Dependencies
- Tests (T001-T005) before implementation (T006-T009)
- T006 blocks T007, T008, T009
- T007 blocks T008, T009, T010
- Implementation before polish (T012-T014)

## Parallel Example
```
# Launch T001-T005 together:
Task: "Integration test for setting ACL rules via POST /events in tests/integration/test_acl_setup.go"
Task: "Integration test for permission denied on unauthorized event in tests/integration/test_acl_denied.go"
Task: "Integration test for permission granted on authorized event in tests/integration/test_acl_granted.go"
Task: "Integration test for root user bypass in tests/integration/test_acl_root_bypass.go"
Task: "Integration test for retrieving ACL events in tests/integration/test_acl_retrieve.go"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task
- Avoid: vague tasks, same file conflicts

## Task Generation Rules
*Applied during main() execution*

1. **From Quickstart**:
   - Each test scenario → integration test task [P]

2. **From Data Model**:
   - Each entity → model extension task [P]
   - ACL service → service creation task [P]

3. **From Plan**:
   - Handler integration → implementation tasks
   - Middleware integration → integration tasks

4. **Ordering**:
   - Setup → Tests → Models → Services → Handpoints → Integration → Polish
   - Dependencies block parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [ ] All quickstart scenarios have corresponding tests
- [ ] All entities have model tasks
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task