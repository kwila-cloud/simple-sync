# Feature Specification: Build a docs site to replace our docs directory

**Feature Branch**: `003-build-docs-site`  
**Created**: 2025-09-23  
**Status**: Draft  
**Input**: User description: "Build a docs site to replace our docs directory."

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
As a developer or user of the simple-sync project, I want to access the documentation through a user-friendly web site instead of navigating raw markdown files in the docs directory, so that I can easily browse and read the documentation.

### Acceptance Scenarios
1. **Given** the docs directory contains markdown files, **When** I build the docs site, **Then** I can access a web page with navigation and rendered content.
2. **Given** the docs site is built, **When** I click on a documentation section, **Then** I see the content properly formatted.

### Edge Cases
- What happens when the docs directory is empty?
- How does the system handle broken links or missing files?

## Requirements *(mandatory)*

### Functional Requirements
- **FR-001**: System MUST generate a static web site from markdown files in the docs directory.
- **FR-002**: System MUST provide navigation between documentation pages.
- **FR-003**: System MUST render markdown content properly (headings, lists, code blocks, etc.).
- **FR-004**: System MUST be accessible via a web browser.
- **FR-005**: System MUST preserve the structure of the docs directory in the site navigation.

### Key Entities *(include if feature involves data)*
- **Documentation Page**: Represents a markdown file with title, content, and path.
- **Site Navigation**: Hierarchical structure based on directory organization.

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

## Implementation Plan

**Branch**: `003-build-docs-site` | **Date**: 2025-09-23 | **Spec**: /home/aemig/Documents/repos/kwila/simple-sync/specs/003-build-docs-site/spec.md
**Input**: Feature specification from /home/aemig/Documents/repos/kwila/simple-sync/specs/003-build-docs-site/spec.md

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
Generate a static documentation website from the existing docs directory using Astro and Starlight framework, with GitHub Pages hosting and automated PDF generation.

## Technical Context
**Language/Version**: JavaScript/Node.js (Astro framework)  
**Primary Dependencies**: Astro, Starlight, starlight-to-pdf  
**Storage**: N/A (static site generation)  
**Testing**: Manual testing for build and deployment  
**Target Platform**: Web browsers  
**Project Type**: single (documentation site)  
**Performance Goals**: Fast static site loading  
**Constraints**: Static generation, GitHub Pages hosting  
**Scale/Scope**: Documentation for simple-sync project

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- RESTful API Design: N/A (documentation site, no API endpoints)
- Event-Driven Architecture: N/A (static site, no data model)
- Authentication and Authorization: N/A (public documentation)
- Data Persistence: N/A (static files)
- Security and Access Control: N/A (public docs)

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

**Structure Decision**: Single project with docs/ directory for Astro/Starlight source

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

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
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
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

**Ordering Strategy**:
- TDD order: Tests before implementation 
- Dependency order: Models before services before UI
- Mark [P] for parallel execution (independent files)

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
*Based on Constitution v1.1.0 - See `/memory/constitution.md`*

---

## Research

## Astro Framework Research
**Decision**: Use Astro as the static site generator for the documentation site.

**Rationale**: Astro is designed for content-focused websites like documentation, provides excellent performance with static generation, and has great integration with Starlight for documentation sites. It's lightweight and focused on content delivery.

**Alternatives considered**:
- Next.js: More complex for a docs site, overkill for static content
- Hugo: Good for docs but less flexible for customization
- Docusaurus: Similar to Starlight but Astro provides better performance

## Starlight Theme Research
**Decision**: Use Starlight as the documentation theme on top of Astro.

**Rationale**: Starlight is specifically built for documentation sites, provides excellent navigation, search, and theming out of the box. It integrates seamlessly with Astro and supports markdown content structure.

**Alternatives considered**:
- Docusaurus: Similar features but heavier
- MkDocs: Python-based, not JavaScript
- Custom theme: Would require more development time

## GitHub Pages Hosting Research
**Decision**: Host the site on GitHub Pages with automatic deployment.

**Rationale**: GitHub Pages is free, integrates directly with the repository, and supports custom domains. GitHub Actions can automate the build and deployment process.

**Alternatives considered**:
- Netlify: More features but adds external dependency
- Vercel: Similar to Netlify
- Self-hosted: More complex infrastructure

## PDF Generation Research
**Decision**: Use starlight-to-pdf tool for generating PDF versions of the documentation.

**Rationale**: The tool is specifically designed for Starlight sites, integrates into the build process, and provides a clean PDF output for offline reading.

**Alternatives considered**:
- Puppeteer custom script: More complex to implement
- Other PDF tools: May not integrate as well with Starlight

## Build Process Research
**Decision**: Use GitHub Actions for CI/CD with build and deployment automation.

**Rationale**: Integrates with GitHub Pages, can run on every push to main branch, and can include PDF generation in the workflow.

**Alternatives considered**:
- Manual deployment: Error-prone and time-consuming
- Other CI services: Adds external dependencies

---

## Data Model

## Project Structure (Astro + Starlight)

### Source Code Layout
```
docs/
├── astro.config.mjs         # Astro configuration with Starlight integration
├── package.json             # Dependencies and scripts
├── tsconfig.json            # TypeScript configuration
├── src/
│   ├── content/
│   │   └── docs/
│   │       ├── index.md     # Homepage
│   │       ├── api/
│   │       │   ├── v1.md   # API v1 documentation
│   │       │   └── ...
│   │       ├── acl.md      # ACL documentation
│   │       ├── tech-stack.md # Technology stack
│   │       └── ...
│   └── env.d.ts            # TypeScript environment
├── public/                  # Static assets (images, etc.)
└── dist/                    # Build output (generated)
```

### Content Organization
- **Root Level**: Main sections in src/content/docs/
- **Subdirectories**: Organized by feature or component
- **Files**: Markdown files with frontmatter for metadata

### Frontmatter Schema
Each markdown file includes:
```yaml
---
title: "Page Title"
description: "Brief description"
sidebar:
  label: "Display Label"
  order: 1
---
```

### Navigation Model
- **Sidebar**: Hierarchical navigation based on directory structure
- **Breadcrumbs**: Path-based navigation
- **Search**: Full-text search across all content
- **Table of Contents**: Auto-generated from headings

### Content Types
- **Reference Docs**: API endpoints, configuration
- **Guides**: Tutorials, setup instructions
- **Examples**: Code samples, use cases

## Build Model

### Static Generation
- Markdown files in src/content/docs/ processed at build time
- HTML output in dist/ with navigation and styling
- Assets (images, CSS, JS) optimized and bundled

### Deployment Model
- Built site in dist/ pushed to GitHub Pages
- PDF version generated and included
- Automatic updates on repository changes

---

## Quickstart

## Prerequisites
- Node.js 18+
- GitHub repository with docs/ directory

## Local Development

### 1. Install Dependencies
```bash
cd docs
npm install
```

### 2. Start Development Server
```bash
npm run dev
```
Open http://localhost:4321 to view the site.

### 3. Build for Production
```bash
npm run build
```

## Deployment

### GitHub Pages Setup
1. Go to repository Settings > Pages
2. Set source to "GitHub Actions"
3. The workflow will automatically deploy on pushes to main

### Manual Deployment
```bash
npm run build
# Copy dist/ contents to GitHub Pages branch
```

## PDF Generation
The PDF is automatically generated during the build process using starlight-to-pdf.

## Testing
- Check all links work
- Verify search functionality
- Test on different browsers
- Validate PDF generation

---

## Tasks

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