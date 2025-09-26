
# Implementation Plan: Implement ACL System

**Branch**: `005-implement-acl-system` | **Date**: 2025-09-26 | **Spec**: /home/aemig/Documents/repos/kwila/simple-sync/specs/005-implement-acl-system/spec.md
**Input**: Feature specification from /home/aemig/Documents/repos/kwila/simple-sync/specs/005-implement-acl-system/spec.md

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
Implement a comprehensive ACL system for access control on event operations. The system will enforce permissions based on user, item, and action rules stored as events, with specificity scoring and timestamp resolution. ACL logic will be centralized in an AclService, integrated into all relevant API handlers replacing TODO(#5) comments.

## Technical Context
**Language/Version**: Go 1.25  
**Primary Dependencies**: Gin web framework, API key authentication  
**Storage**: SQLite database  
**Testing**: Go built-in testing framework  
**Target Platform**: Linux server  
**Project Type**: single (backend API)  
**Performance Goals**: ACL evaluation p95 latency <10ms per request under 100 concurrent evaluations  
**Constraints**: ACL logic centralized in acl_service.go with AclService; integrate into all handlers with TODO(#5) comments  
**Scale/Scope**: Support up to 10,000 ACL rules and 1,000 concurrent users, with expected 100 ACL evaluations per second and linear growth in rules over time

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution file is template with placeholders; no specific principles defined. Proceeding with standard Go development practices.

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
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
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure]
```

**Structure Decision**: Option 1 (single project) - backend API with existing structure

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
    - ACL evaluation algorithm and specificity scoring
    - Best practices for ACL in Go/Gin applications
    - Integration patterns for centralized service in existing handlers

2. **Generate and dispatch research agents**:
    ```
    Task: "Research ACL implementation patterns in Go web services"
    Task: "Find best practices for permission evaluation algorithms"
    Task: "Research SQLite performance for ACL rule queries"
    ```

3. **Consolidate findings** in `research.md` using format:
    - Decision: [what was chosen]
    - Rationale: [why chosen]
    - Alternatives considered: [what else evaluated]

**Output**: research.md with all unknowns resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
    - ACL Rule: user, item, action, allow/deny, timestamp
    - Event: existing structure with ACL rules stored as events
    - Relationships: ACL rules linked to events via .acl item

2. **Generate API contracts** from functional requirements:
    - Existing POST /events: Add ACL check and support .acl events
    - Existing GET /events: Filter by ACL permissions
    - Use REST patterns with JSON
    - Output OpenAPI schema to `/contracts/`

3. **Generate contract tests** from contracts:
    - Test ACL checks on events endpoints
    - Tests must fail initially

4. **Extract test scenarios** from user stories:
    - ACL enforcement on event submission
    - Permission denied scenarios
    - Root user bypass

5. **Update agent file incrementally** (O(1) operation):
    - Run `.specify/scripts/bash/update-agent-context.sh opencode`
      **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
    - Add ACL service and evaluation logic
    - Preserve existing context

**Output**: data-model.md, /contracts/*, quickstart.md, updated AGENTS.md

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs
- ACL model creation [P]
- AclService implementation [P]
- ACL endpoint handlers
- Integration into existing event handlers
- Contract and integration tests

**Ordering Strategy**:
- TDD order: Tests first
- Models → Service → Handlers → Integration
- Parallel where independent

**Estimated Output**: 15-20 tasks in tasks.md

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
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [ ] Complexity deviations documented

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
