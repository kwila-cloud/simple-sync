
# Implementation Plan: High Performance REST API for simple-sync (Issue #2 Focus)

**Branch**: `001-we-want-to` | **Date**: 2025-09-20 | **Spec**: /home/aemig/Documents/repos/kwila/simple-sync/specs/001-we-want-to/spec.md
**Input**: Feature specification from `/specs/001-we-want-to/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Implement GitHub issue #2: Basic event storage with GET/POST /events endpoints using Go 1.25 and Gin framework. Focus on in-memory storage, thread safety, and core API functionality as foundation for the simple-sync REST API.

## Technical Context
**Language/Version**: Go 1.25  
**Primary Dependencies**: Gin web framework, golang-jwt for future auth  
**Storage**: In-memory for issue #2, file-based JSON for future persistence  
**Testing**: Go built-in testing framework with contract and integration tests  
**Target Platform**: Linux server with Docker deployment  
**Project Type**: Single backend API project  
**Performance Goals**: High-performance REST API handling concurrent requests efficiently  
**Constraints**: Simple deployment, thread-safe in-memory storage, proper error handling  
**Scale/Scope**: MVP for event storage with GET/POST endpoints, foundation for full API

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- RESTful API Design: All planned endpoints MUST use appropriate HTTP methods and follow resource-oriented patterns.
- Event-Driven Architecture: Data model MUST be based on timestamped events with user/item metadata.
- Authentication and Authorization: JWT auth and ACL permissions MUST be integrated into the design.
- Data Persistence: File-based JSON storage MUST be used for data survival.
- Security and Access Control: ACL rules MUST be evaluated with deny-by-default behavior.

## Project Structure

### Documentation (this feature)
```
specs/001-we-want-to/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - focused on issue #2)
```

### Source Code (repository root)
```
# Option 1: Single project (DEFAULT)
src/
├── models/
├── handlers/
├── storage/
└── main.go

tests/
├── contract/
├── integration/
└── unit/
```

**Structure Decision**: Option 1 - Single project for backend API, focused on issue #2 implementation

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
    - Research Go 1.25 best practices for REST APIs with Gin
    - Research thread-safe in-memory storage patterns
    - Research JSON handling and error responses in Gin

2. **Generate and dispatch research agents**:
    ```
    For each unknown in Technical Context:
      Task: "Research {unknown} for issue #2 context"
    For each technology choice:
      Task: "Find best practices for {tech} in Go REST API development"
    ```

3. **Consolidate findings** in `research.md` using format:
    - Decision: [what was chosen]
    - Rationale: [why chosen]
    - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from issue #2** → `data-model.md`:
    - Event entity with required fields
    - Validation rules for Event struct
    - In-memory storage design

2. **Generate API contracts** from issue #2 requirements:
    - GET /events endpoint
    - POST /events endpoint
    - GET /events?fromTimestamp=X endpoint
    - Output OpenAPI schema to `/contracts/`

3. **Generate contract tests** from contracts:
    - One test file per endpoint
    - Assert request/response schemas
    - Tests must fail (no implementation yet)

4. **Extract test scenarios** from issue #2 acceptance criteria:
    - Basic event storage scenarios
    - Timestamp filtering scenarios
    - Error handling scenarios

5. **Update agent file incrementally** (O(1) operation):
    - Run `.specify/scripts/bash/update-agent-context.sh opencode` for your AI assistant
    - If exists: Add only NEW tech from current plan
    - Preserve manual additions between markers
    - Update recent changes (keep last 3)
    - Keep under 150 lines for token efficiency
    - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks focused on issue #2 requirements
- Each contract → contract test task [P]
- Event entity → model creation task
- In-memory storage → storage implementation task
- Each endpoint → handler implementation task
- Main.go setup → integration task

**Ordering Strategy**:
- TDD order: Tests before implementation
- Dependency order: Setup → Tests → Model → Storage → Handlers → Main → Polish
- Mark [P] for parallel execution (independent files)

**Estimated Output**: 20 focused tasks for issue #2 in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [x] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented

---
*Based on Constitution v1.0.0 - See `/memory/constitution.md`*
