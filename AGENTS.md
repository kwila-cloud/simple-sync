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
  - Checklist format using `[ ]` for pending and `[x]` for completed
- **Style**: 
  - ✅ Good: Simple, scannable checklist format
  - ✅ Good: Group related items logically
  - ❌ Avoid: Verbose descriptions, detailed explanations, multiple sections

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

### Naming Scheme

- **Code Names (variables, methods, functions)**: Don't put acronyms as all-caps in names. For example, use `AclRule` rather than `ACLRule` for an ACL rule.
  - This way, words are clearly differentiated by case.
  - Examples: `CreateApiKey`, `GetAclRules`, `UpdateAclRule` (NOT: `CreateAPIKey`, `GetACLRules`)
- **Documentation**: Use normal capitalization for acronyms in plain text, comments, and documentation.
  - Examples: "API key", "ACL rule", "REST API" (NOT: "ApiKey", "AclRule" in documentation)
