# Implementation Plan: High Performance REST API for simple-sync

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
Build a high-performance REST API backend for simple-sync that supports event storage, ACL-based access control, JWT authentication, and easy deployment via Docker Compose. The API will be flexible for multiple frontend apps, using Go with Gin framework, file-based JSON storage, and following REST principles.

## Technical Context
**Language/Version**: Go 1.25  
**Primary Dependencies**: Gin web framework, golang-jwt for authentication, SQLite for data storage  
**Storage**: File-based JSON storage for events and ACL data  
**Testing**: Go built-in testing framework with contract and integration tests  
**Target Platform**: Linux server with Docker deployment  
**Project Type**: Single backend API project  
**Performance Goals**: High-performance REST API handling concurrent requests efficiently  
**Constraints**: Simple deployment with docker-compose, minimal configuration, data persistence across restarts  
**Scale/Scope**: MVP for event synchronization with ACL, supporting multiple frontend apps

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
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
# Option 1: Single project (DEFAULT)
src/
├── models/
├── handlers/
├── middleware/
├── storage/
└── config/

tests/
├── contract/
├── integration/
└── unit/

data/
├── events.json
└── acl.json
```

**Structure Decision**: Option 1 - Single project as this is a backend API without frontend components.

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - Research Go best practices for REST APIs with Gin
   - Research JWT implementation patterns in Go
   - Research file-based JSON storage patterns for concurrent access
   - Research ACL implementation patterns

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Event: uuid, timestamp, userUuid, itemUuid, action, payload
   - User: uuid, username, password
   - ACL Entry: userUuid, itemUuid, permissions
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - GET/POST /events endpoints
   - GET/PUT /acl endpoints
   - POST /auth/token endpoint
   - Admin endpoints for user management
   - Use standard REST patterns
   - Output OpenAPI schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Docker deployment scenario
   - Authentication and event operations
   - ACL permission enforcement
   - Quickstart test = story validation steps

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
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- Each contract → contract test task [P]
- Each entity → model creation task [P] 
- Each user story → integration test task
- Implementation tasks to make tests pass
- Work on GitHub issues #2, #7, #4, #5, #6, #1 in order

**Ordering Strategy**:
- TDD order: Tests before implementation 
- Dependency order: Models before handlers before middleware
- Mark [P] for parallel execution (independent files)
- Stop for PR after each issue completion

**Estimated Output**: 25-30 numbered, ordered tasks in tasks.md

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
| None | N/A | N/A |

## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented

---
*Based on Constitution v1.0.0 - See `/memory/constitution.md`*