# Tasks: Issue #2 - Basic Event Storage with GET/POST /events Endpoints

**Input**: GitHub issue #2 description and acceptance criteria
**Prerequisites**: Go 1.25 installed, basic understanding of Gin framework

## Execution Flow (main)
```
1. Load issue #2 description and requirements
   → Extract: Event model fields, API endpoints, storage requirements
2. Generate focused tasks for issue #2 implementation
   → Setup: Project structure and dependencies
   → Tests: Contract tests for endpoints
   → Core: Event model, in-memory storage, handlers
   → Integration: Main application setup
   → Polish: Basic validation and documentation
3. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
4. Number tasks sequentially (T001, T002...)
5. Generate dependency graph
6. Create parallel execution examples
7. Validate task completeness for issue #2
8. Return: SUCCESS (tasks ready for issue #2 implementation)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- Paths assume Go project structure with src/ and tests/

## Phase 3.1: Setup
- [ ] T001 Initialize Go 1.25 module with go.mod in repository root
- [ ] T002 Add Gin web framework dependency to go.mod
- [ ] T003 Create src/ and tests/ directories at repository root
- [ ] T004 [P] Configure gofmt for code formatting

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T005 [P] Contract test for GET /events endpoint in tests/contract/test_get_events.go
- [ ] T006 [P] Contract test for POST /events endpoint in tests/contract/test_post_events.go
- [ ] T007 [P] Contract test for GET /events?fromTimestamp=X endpoint in tests/contract/test_get_events_timestamp.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T008 Create Event struct with fields: uuid, timestamp, userUuid, itemUuid, action, payload in src/models/event.go
- [ ] T009 Implement thread-safe in-memory storage using slice and RWMutex in src/storage/memory.go
- [ ] T010 Implement GET /events handler that returns all events as JSON array in src/handlers/events.go
- [ ] T011 Implement GET /events?fromTimestamp=X handler with timestamp filtering in src/handlers/events.go
- [ ] T012 Implement POST /events handler that accepts event array and returns updated history in src/handlers/events.go
- [ ] T013 Add JSON serialization/deserialization with error handling for malformed requests in src/handlers/events.go

## Phase 3.4: Integration
- [ ] T014 Set up main.go with Gin router, register /events routes, and configurable port (default 8080)
- [ ] T015 Add proper HTTP status codes (200, 400, 500) to all handlers
- [ ] T016 Implement concurrent request handling with thread-safe storage access

## Phase 3.5: Polish
- [ ] T017 [P] Add basic unit tests for Event model validation in tests/unit/test_event_model.go
- [ ] T018 Run all contract tests to ensure they pass after implementation
- [ ] T019 Test server startup on port 8080 and basic endpoint responses
- [ ] T020 [P] Add inline documentation and comments following Go best practices

## Dependencies
- Setup (T001-T004) before everything
- Tests (T005-T007) before implementation (T008-T016)
- Event model (T008) before storage (T009)
- Storage (T009) before handlers (T010-T013)
- Handlers (T010-T013) before main.go (T014)
- Main.go (T014) before status codes and concurrency (T015-T016)
- Implementation before polish (T017-T020)

## Parallel Example
```
# Launch T005-T007 together (all contract test files are independent):
Task: "Contract test for GET /events endpoint in tests/contract/test_get_events.go"
Task: "Contract test for POST /events endpoint in tests/contract/test_post_events.go"
Task: "Contract test for GET /events?fromTimestamp=X endpoint in tests/contract/test_get_events_timestamp.go"

# Launch T017 and T020 together (unit test and docs are independent):
Task: "Add basic unit tests for Event model validation in tests/unit/test_event_model.go"
Task: "Add inline documentation and comments following Go best practices"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task
- Focus only on issue #2 requirements - no auth, ACL, or persistence yet
- Use Go 1.25 idioms and best practices
- Ensure all Event struct fields match specification exactly
- Handle malformed JSON gracefully with appropriate error responses

## Task Generation Rules
*Applied during main() execution*

1. **From Issue Requirements**:
   - Each API endpoint → contract test task [P]
   - Event model → model creation task
   - In-memory storage → storage implementation task
   - Each handler → implementation task (sequential in same file)

2. **From Acceptance Criteria**:
   - Project structure → setup tasks
   - Thread safety → storage task
   - HTTP status codes → integration task
   - Error handling → handler tasks

3. **Ordering**:
   - Setup → Tests → Model → Storage → Handlers → Main → Polish
   - Dependencies block parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [x] All API endpoints have corresponding contract tests
- [x] Event model matches specification exactly
- [x] Thread-safe storage implemented
- [x] All tests come before implementation
- [x] Parallel tasks truly independent
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Focused only on issue #2 scope