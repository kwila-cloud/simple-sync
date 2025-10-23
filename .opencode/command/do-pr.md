---
description: Develop the changes for a new pull request
agent: build
---

IMPORTANT: This command must run the full non-interactive flow for creating a PR. That means it MUST run the test suite(s), commit any changes, push the branch, create the GitHub pull request, update `CHANGELOG.md` with the PR number, and push the changelog — all without asking the user for additional input.

If the user has NOT previously run the `/start-pr` command, exit early and ask them to run that command.

Required behavior (non-interactive flow)

1. Use the `todoread` tool to get the TODO list for the current issue.
2. For each item in the TODO list:
   - Implement the change as described in the TODO item.
   - Run the relevant automated tests immediately after implementing the change. Tests must be run and pass before committing. Typical commands to run are:
     - Unit/contract/integration (with race):
       - `go test -race ./tests/unit ./tests/contract ./tests/integration`
     - Performance tests (optional, without race):
       - `go test ./tests/performance`
     - If a change only affects unit tests, run the narrower set of packages to save time.
   - If tests fail, refine the code until tests pass. Do not proceed to committing that TODO item until its tests pass.
   - Once tests pass, update the spec (check off corresponding item) and commit the change locally using a descriptive conventional commit message (example `feat(7): add backup script`).
     - Use: `git add -A && git commit -m "<scope>: <short description>"`
3. After all TODO items are completed and committed locally:
   - Push the branch to the remote:
     - `git push -u origin "$(git rev-parse --abbrev-ref HEAD)"`
   - Create the pull request non-interactively using `gh` (GitHub CLI). Provide a clear title and a PR body via a HEREDOC to avoid shell quoting issues. Example:
     - `gh pr create --title "<PR title>" --body "$(cat <<'EOF'\n<PR body>\nEOF\n)"`
   - Retrieve the created PR number and URL to update `CHANGELOG.md`.
     - Example to get the PR number: `gh pr view --json number,url --jq '.number'`
   - Update `CHANGELOG.md` at the top (under the current unreleased section) with a single entry referencing the PR number in the repository style, for example:
     - `- [#NN](https://github.com/kwila-cloud/simple-sync/pull/NN): Short description`
   - Commit the `CHANGELOG.md` update and push the commit (it must be on the same branch so the changelog change is included in the PR):
     - `git add CHANGELOG.md && git commit -m "chore: add changelog entry for PR #NN" && git push`

Error handling and constraints

- The command must NOT prompt the user for extra confirmation during the flow. If an operation would normally require input (for example, `gh pr create` in interactive mode), invoke the non-interactive flags and provide the input programmatically (HEREDOC or CLI flags).
- If network push or GH CLI operations fail, surface the error and abort; do not attempt destructive recovery automatically.

Examples

- Commit and run tests for a single TODO item:

  ```bash
  # implement change (edit files)
  go test -race ./tests/unit ./tests/contract ./tests/integration
  git add -A
  git commit -m "feat(storage): add backup script"
  ```

- Create PR non-interactively and update changelog (example body uses HEREDOC):

  ```bash
  git push -u origin "$(git rev-parse --abbrev-ref HEAD)"
  gh pr create --title "docs: documentation and configuration updates" --body "$(cat <<'EOF'\n## Summary\n- Add backup/restore scripts and docs\n\n## Files\n- scripts/backup.sh\n- scripts/restore.sh\n- docs/data-persistence.mdx\nEOF\n)"
  PR_NUMBER=$(gh pr view --json number --jq '.number')
  sed -i "1i- [#$PR_NUMBER](https://github.com/kwila-cloud/simple-sync/pull/$PR_NUMBER): Documentation and configuration updates" CHANGELOG.md
  git add CHANGELOG.md && git commit -m "chore: add changelog entry for PR #$PR_NUMBER" && git push
  ```

Notes

- Use HEREDOC with a single-quoted delimiter (`<<'EOF'`) when the PR body contains backticks or other shell-sensitive characters to avoid unintended shell expansion.
- Prefer running the narrowest test set that verifies the change to save CI time, but ensure integration/contract tests run when changes affect those areas.

