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

- **Storage & Docker notes (for future agents):**
  - The project uses SQLite for persistent storage; production/dev processes expect a local file under `./data` by default (`./data/simple-sync.db`).
  - Docker Compose is configured to bind-mount `./data` into the container (`./data:/app/data`) so the DB file is stored on the host. This is the recommended setup for development and simple deployments because it makes backups and inspection straightforward.
  - If you need a Docker-managed volume, a named volume `simple-sync-data` exists in `docker-compose.yml` (commented). Using a named volume is fine for environments where host access is not required, but it makes manual backups/restores less obvious.
  - Building the Go binary with `github.com/mattn/go-sqlite3` requires a system C toolchain and headers (`libsqlite3-dev`) and `CGO_ENABLED=1` during build. CI workflows already install `libsqlite3-dev` where needed.
  - For backups, there are helper scripts in `./scripts/backup.sh` and `./scripts/restore.sh` that copy the DB file from `./data` to `./backups` and restore it safely; future agents should prefer these scripts for simple backup/restore workflows.

### Project Knowledge

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

#### Tests Location
- **Policy:** Tests MUST be placed under `tests/` organized by type (`tests/unit/`, `tests/integration/`, `tests/contract/`, `tests/performance/`).
  - DO NOT place tests in `src/` OR any subdirectory of `src/`
- **Why:** Keeps code vs test separation clear; matches existing repo layout and CI patterns.

### Git Workflow

- Feature branches for issues (e.g., `63-new-setting`)
- Use GitHub CLI for PR creation: `gh pr create`

### Shell quoting when using backticks

- **Problem:** Unescaped backticks in bash commands are interpreted as command substitution, causing the shell to execute the content between backticks instead of passing it as literal text (this can break `gh` calls or insert unintended output).
- **Rule:** When passing text that contains backticks to shell commands, avoid unescaped backticks. Prefer one of these safe patterns:
  - Single-quoted argument (simple cases):
    - `gh pr edit 56 --body 'Tables: `users`, `events`'`
  - Escape backticks inside double quotes:
    - `gh pr edit 56 --body "Tables: \`users\`, \`events\`"`
  - HEREDOC with single-quoted delimiter (recommended for multi-line bodies or complex content):
    - `gh pr edit 56 --body "$(cat <<'EOF'\nTables: `users`, `events`\nEOF\n)"`

- **Examples:**
  - Bad: `gh pr edit 56 --body "Tables: `users`, `events`"`  (backticks executed by shell)
  - Good (HEREDOC): `gh pr edit 56 --body "$(cat <<'EOF'\nTables: `users`, `events`\nEOF\n)"`
  - Good (single quotes): `gh pr edit 56 --body 'Tables: `users`, `events`'`
  - Good (escaped): `gh pr edit 56 --body "Tables: \`users\`, \`events\`"`

- **Recommendation:** Prefer the HEREDOC pattern when generating multi-line PR bodies that include code formatting or backticks. It avoids shell expansion and is easy to read and maintain.

- Commit messages follow conventional format: `feat:`, `refactor:`, `chore:`, `fix:`, etc.

### Grep / Ripgrep Patterns

- **Problem:** The repository's search tooling (e.g. `rg`/`ripgrep`) uses regular expressions by default. Supplying an invalid regex (for example, an unclosed grouping like `AddUser(`) causes an immediate error and interrupts automated search tools.
- **Rule:** Always provide a valid regular expression to regex-based search tools, or use fixed-string mode when searching for literal text containing regex metacharacters.
- **How to search safely:**
  - Use fixed-string mode for literals: `rg -F 'AddUser('` or `rg --fixed-strings 'AddUser('
  - Escape regex metacharacters: `rg 'AddUser\('` (escape `(` with `\` in single-quoted shell strings)
  - Search the identifier only (no parens): `rg 'AddUser'`
  - Prefer single quotes around patterns to avoid shell interpolation: `rg 'GetUserById\('`
  - When using the assistant `functions.grep` tool, pass a syntactically valid regex (escape metacharacters) or a simple identifier-only pattern.
- **Examples:**
  - Bad: `rg "AddUser("` → causes ripgrep regex parse error (unclosed group)
  - Good (escape): `rg 'AddUser\('`
  - Good (fixed-string): `rg -F 'AddUser('
  - Good (identifier only): `rg 'AddUser'`

- **Recommendation:** When programmatically constructing search patterns, either validate the regex before use or default to fixed-string searches. If you are unsure whether a pattern contains regex metacharacters, use `-F` to avoid surprises.

### PR Title & Description Rules
- **Always inspect the full diff for the branch before creating a PR.** Use Git to view changes against the base branch and confirm the final, combined diff that will become the PR.

  - View commits on your branch relative to `main`:
    - `git fetch origin && git log --oneline origin/main...HEAD`
  - View a concise file/status diff against `main`:
    - `git fetch origin && git diff --name-status origin/main...HEAD`
  - View a human-readable patch summary before creating the PR:
    - `git fetch origin && git diff --stat origin/main...HEAD`
  - View the full patch if needed:
    - `git fetch origin && git diff origin/main...HEAD`

- **Title rules:**
  - The PR title must be a short, descriptive summary of the change and MUST NOT include the issue number.
  - Example — incorrect: `docs(7): documentation and restore script fixes for data persistence`
  - Example — correct: `docs: documentation and restore script fixes for data persistence`

- **Description rules:**
  - The PR body/description MUST include the issue number (and link) that the PR addresses. Prefer an explicit `Fixes #<issue>` or `Related to #<issue>` line.
  - Example body snippet:

    Issue: https://github.com/kwila-cloud/simple-sync/issues/7

    This PR moves the data-persistence doc into the site content and adds a checklist to verify `restore.sh` behavior.

    Fixes #7

- **How to get the final PR number non-interactively:**
  - After creating the PR, capture the assigned number:
    - `gh pr view --json number,url --jq '.number'` (use the branch name if needed)

### Changelog
- **ALWAYS add a new line to `CHANGELOG.md` for each new *pull request* (PR).** Do not include or link to the issue number — reference the PR number only.
- Document new features, enhancements, bug fixes, and breaking changes.
- Follow the existing format with PR links and clear descriptions (see examples below).
- Keep entries concise but descriptive for users and maintainers.
- **IMPORTANT**: Always verify the actual PR details **and the PR title** before updating the changelog. Use `gh pr view <PR-number>` or `gh pr view <branch>` to confirm the PR number, title, and changed files.

- **CRITICAL**: Add exactly ONE entry per PR. Never add multiple entries for the same pull request, even if the PR contains multiple types of changes. Combine all changes into a single, concise description.

Examples (incorrect vs correct):

- Incorrect (links to issue rather than PR, or includes issue numbers in the text):
  - `- [#7](https://github.com/kwila-cloud/simple-sync/issues/7): Implement ACL rule storage (CreateAclRule, GetAclRules)`

- Correct (match existing style — use PR link and number, no issue numbers in entry):
  - `- [#59](https://github.com/kwila-cloud/simple-sync/pull/59): Implement ACL rule storage (CreateAclRule, GetAclRules)`

Recommended workflow to avoid mistakes:

1. Inspect the full diff using the `git` commands above and confirm the intended changes.
2. Create the PR non-interactively using `gh pr create` with a title (no issue number) and a body that includes the issue link/number.
3. Immediately run `gh pr view --json number,url,title --jq '.number'` to get the PR number.
4. Update `CHANGELOG.md` at the top (under the current unreleased version) with a single entry referencing the **PR number only**.
5. Commit the changelog update to the same branch so the changelog change is included in the PR.

Notes:
- The changelog should always reference the PR number (not the issue number) because the PR is the canonical unit that contains the implemented changes and the exact diff reviewers will see.
- If the PR number is not yet available (for example, you plan to open the PR later), delay updating `CHANGELOG.md` until the PR is created so you can reference the correct PR number.
- Keep the entry format consistent with existing history: `- [#NN](https://github.com/kwila-cloud/simple-sync/pull/NN): Short description`

### Opencode Commands

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
- **Database Tables**: Use singular form for database tables
  - Example: "user" rather than "users"

### Idiomtic Go


#### Looping

Prefer `for i := range n` (Go 1.22+) over `for i := 0; i < n; i++` when iterating a fixed integer count; prefer `for i := range slice` when iterating slices.

#### Standard Library

**CRITICAL**: Always use functions from the standard library when possible, rather than creating custom helper functions.

1. **Prefer Standard Library**: Check if Go's standard library already provides the functionality you need
   - ✅ Good: `strings.Contains()`, `strings.HasPrefix()`, `strings.Split()`
   - ✅ Good: `fmt.Sprintf()`, `strconv.Atoi()`, `time.Parse()`
   - ❌ Bad: Creating custom `contains()`, `split()`, `parseTime()` functions

2. **Common Standard Library Functions to Use**:
   - **Strings**: `strings.Contains()`, `strings.HasPrefix()`, `strings.HasSuffix()`, `strings.Split()`, `strings.Join()`
   - **Formatting**: `fmt.Sprintf()`, `fmt.Errorf()`
   - **Conversions**: `strconv.Atoi()`, `strconv.Itoa()`, `strconv.ParseBool()`
   - **Time**: `time.Now()`, `time.Parse()`, `time.Format()`
   - **Slices**: `sort.Slice()`, `append()`, `copy()`

3. **Examples**:
   - ✅ Good: `if strings.Contains(err.Error(), "malformed") { ... }`
   - ❌ Bad: Creating custom `contains(s, substr string) bool` function
   - ✅ Good: `result := fmt.Sprintf("User %s has %d items", user, count)`
   - ❌ Bad: Creating custom string formatting helpers

4. **When in doubt**: If you're about to write a helper function, first check the Go standard library documentation for existing solutions
