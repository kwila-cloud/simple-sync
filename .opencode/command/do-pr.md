---
description: Develop the changes for a new pull request
agent: build
---

IMPORTANT: This command must run the full non-interactive flow for creating a PR. That means it MUST run the test suite(s), commit any changes, push the branch, create the GitHub pull request, update `CHANGELOG.md` with the PR number, and push the changelog — all without asking the user for additional input.

The user gave the input: "$ARGUMENTS"

Use the user input as the issue number.

If the user input is empty or invalid, use the previously entered issue number from `/start-pr` (but if `/start-pr` was not previously ran, prompt the user for the issue number).

Required behavior (non-interactive flow)

1. Read the spec for the given issue in `specs/` and determine the next incomplete section from the Task List.
2. For each task in the incomplete section:
   - Implement the task.
   - Run the relevant automated tests immediately after implementing the change. Tests must be run and pass before committing. Typical commands to run are:
     - Unit/contract/integration (with race):
       - `go test -race ./tests/unit ./tests/contract ./tests/integration`
     - Performance tests (optional, without race):
       - `go test ./tests/performance`
     - If a change only affects unit tests, run the narrower set of packages to save time.
   - If tests fail, refine the code until tests pass. Do not proceed to committing that TODO item until its tests pass.
   - Once tests pass, update the spec (check off corresponding item) and commit the change locally using a descriptive conventional commit message (example `feat(7): add backup script`).
     - Use: `git add -A && git commit -m "<scope>: <short description>"`
3. After all task items for the current section are completed and committed locally:
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
   - Ask the user for code review feedback on the new pull request.
   - DO NOT suggest starting on the next section of the task list.

Error handling and constraints

- The command must NOT prompt the user for extra confirmation during the flow. If an operation would normally require input (for example, `gh pr create` in interactive mode), invoke the non-interactive flags and provide the input programmatically (HEREDOC or CLI flags).
- If network push or GH CLI operations fail, surface the error and abort; do not attempt destructive recovery automatically.

