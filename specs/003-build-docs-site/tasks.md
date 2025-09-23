# Tasks: Build a docs site to replace our docs directory

**Input**: Design documents from `/specs/003-build-docs-site/`
**Prerequisites**: plan.md (required), research.md, data-model.md, quickstart.md

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
- **Docs site**: `docs/` directory at repository root
- All paths relative to repository root unless specified

## Phase 3.1: Setup
- [X] T001 Rename existing docs/ directory to old-docs/
- [X] T002 Initialize Astro + Starlight project in docs/ with npm create astro
- [X] T003 [P] Configure TypeScript and linting in docs/tsconfig.json and docs/package.json

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [X] T004 [P] Integration test for local build in docs/test_build.js
- [X] T005 [P] Integration test for link checking in docs/test_links.js

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [X] T006 Configure Starlight in docs/astro.config.mjs
- [X] T007 Move existing docs content to docs/src/content/docs/

## Phase 3.4: Integration
- [X] T009 Setup GitHub Actions workflow in .github/workflows/deploy-docs.yml
- [X] T008 Add starlight-to-pdf CLI execution with --contents-links internal option to GitHub Actions workflow in .github/workflows/deploy-docs.yml

## Phase 3.5: Polish
- [X] T010 [P] Add performance testing script in docs/test_performance.js
- [X] T011 Run manual testing per quickstart.md
- [X] T012 Update README.md with docs site link

## Dependencies
- T001 blocks T002
- Tests (T004-T005) before implementation (T006-T007)
- T002 blocks T003, T006
- T006 blocks T007
- T009 blocks T008
- Implementation before polish (T010-T012)

## Parallel Example
```
# Launch T004-T005 together:
Task: "Integration test for local build in docs/test_build.js"
Task: "Integration test for link checking in docs/test_links.js"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task
- Avoid: vague tasks, same file conflicts

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   - Build contract → build integration test [P]

2. **From Data Model**:
   - Project structure → setup tasks
   - Content organization → core tasks

3. **From User Stories**:
   - Build site story → build test [P]
   - Access web page story → link check test [P]
   - Click sections story → navigation test

4. **Ordering**:
   - Setup → Tests → Core → Integration → Polish
   - Dependencies block parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [ ] All contracts have corresponding tests
- [ ] All entities have model tasks
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task
