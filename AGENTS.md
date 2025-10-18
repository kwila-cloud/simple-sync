# AGENTS.md - AI Development Guide for Simple-Sync

## Project Knowledge

Simple-sync is a lightweight REST API built in Go that provides event storage and access control functionality. The system allows users to authenticate via setup tokens exchanged for API keys, store timestamped events for specific items, and manage permissions through Access Control Lists (ACLs).

**Technology Stack:**
- Go 1.25 with Gin web framework
- SQLite database storage
- CORS support for web clients

**Core Features:**
- User authentication with API keys
- Event storage with timestamps and metadata
- ACL-based permission system (read/write permissions)

### GitHub Data Access

- **ALWAYS use GitHub CLI for GitHub data** - NEVER use webfetch for issues, PRs, or other GitHub information
- **Issue information**: Use `gh issue view <number>`
- **PR information**: Use `gh pr view <number>`
- **PR diff**: Use `gh pr diff <number>`
- **Examples**:
  - ✅ Good: `gh issue view 7` (clean, structured output)
  - ❌ Bad: `webfetch` with GitHub URL (scraped HTML, verbose output)

### Issue Specifications

- **File Naming**: Use format `{issue-number}-{brief-description}.md` in `specs/`
- **Structure**: 
  - Title with issue link: `# Title\n\nhttps://github.com/kwila-cloud/simple-sync/issues/{number}`
  - Brief plan description
  - Design decisions section (if applicable)
  - Task List with sections corresponding to atomic pull requests
    - Each section header represents one PR
    - Items within section are changes included in that PR (not necessarily atomic)
    - Use `[ ]` for pending and `[x]` for completed
- **Style**: 
  - ✅ Good: Simple, scannable checklist format
  - ✅ Good: Group related items logically
  - ❌ Avoid: Verbose descriptions, detailed explanations, multiple sections
- **TDD Approach**: Each task item should include tests first, then implementation
  - ✅ Good: "Add tests for X", "Implement X"
  - ❌ Bad: Separate testing section at the end
- **Task List Structure**: 
  - Each section = one atomic pull request
  - Items within section = changes included in that PR (can be multiple related changes)
  - ✅ Good: Section with multiple related implementation items
  - ❌ Bad: Each individual item as separate PR

### Specification Development Process

- **CRITICAL**: DO NOT make code changes while working on specifications
- **Specification phase**: Focus only on planning, design, and Task List creation
- **Implementation phase**: After spec is finalized, then proceed with code changes
- **Examples**:
  - ✅ Good: Discussing TDD approaches, test strategies, implementation order
  - ❌ Bad: Reading files, writing tests, or implementing code during spec phase
- **When in doubt**: If you're about to read/write code files, stop - you're still in spec phase

#### Example: Issue #7 Data Persistence
See `specs/7-data-persistence.md` for a well-structured specification that:
- Includes design decisions section explaining SQLite vs Go marshaling choice
- Uses TDD approach with tests first for each implementation item
- Task List sections correspond to atomic pull requests
- Items within sections are related changes for that PR
- Groups related functionality logically
- Maintains focus without excessive detail

### TDD Implementation Process

**CRITICAL**: Always follow Test-Driven Development (TDD) when implementing features:

1. **Tests First**: Write tests BEFORE implementing any code
   - ✅ Good: "Add tests for X", then "Implement X" 
   - ❌ Bad: "Implement X", then "Add tests for X"
   - ✅ Good: Plan implementation order based on test requirements
   - ❌ Bad: Plan implementation without considering test structure

2. **Implementation Order**: 
   - Read the spec task list carefully - tests are always listed before implementation
   - Write failing tests first
   - Implement minimal code to make tests pass
   - Refactor if needed

3. **Examples**:
   - ✅ Good: "First I'll add tests for the new ACL storage methods, then implement the interface methods"
   - ❌ Bad: "I'll implement the ACL storage interface, then add tests for it"

4. **When in doubt**: If you're about to implement code without writing tests first, stop - you're violating TDD principles

### Git Workflow

- Feature branches for issues (e.g., `63-new-setting`)
- Use GitHub CLI for PR creation: `gh pr create`
- Commit messages follow conventional format: `feat:`, `refactor:`, `chore:`, `fix:`, etc.

### Changelog
- **ALWAYS update add a new line to CHANGELOG.md for each new pull request.**
- Document new features, enhancements, bug fixes, and breaking changes
- Follow the existing format with PR links and clear descriptions
- Keep entries concise but descriptive for users and maintainers
- **IMPORTANT**: Always verify the actual PR content before updating the changelog. Use `gh pr view <PR-number>` to check the PR title, body, and changed files to ensure accurate changelog entries.
- **CRITICAL**: Add exactly ONE entry per PR. Never add multiple entries for the same pull request, even if the PR contains multiple types of changes. Combine all changes into a single, concise description.

### Opencode Commands

Opencode commands are defined as separate markdown files in `.opencode/command/`. Each command file should:
- Have a `description` and `agent` in the frontmatter
- Contain the command logic and usage instructions
- Follow the naming convention `{command-name}.md`

See `.opencode/command/` directory for examples.

**Common mistake to avoid:** Do NOT add command documentation to AGENTS.md. Commands belong in their own files in `.opencode/command/`.

### Naming Scheme

- **Code Names (variables, methods, functions)**: Don't put acronyms as all-caps in names. For example, use `AclRule` rather than `ACLRule` for an ACL rule.
  - This way, words are clearly differentiated by case.
  - Examples: `CreateApiKey`, `GetAclRules`, `UpdateAclRule` (NOT: `CreateAPIKey`, `GetACLRules`)
- **Documentation**: Use normal capitalization for acronyms in plain text, comments, and documentation.
  - Examples: "API key", "ACL rule", "REST API" (NOT: "ApiKey", "AclRule" in documentation)
