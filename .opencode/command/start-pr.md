---
description: Start a new pull request by checking specs and creating the appropriate branch
agent: build
---

Automates the process of starting work on a new pull request for the given issue.

The user gave the input: "$ARGUMENTS"

Use the user input as the issue number.

If the user input is empty or invalid, prompt the user for the issue number.

Required behavior and confirmation flow

1. Read the spec for the given issue in `specs/` and determine the next incomplete section from the Task List.
2. Branch creation rules:
   - Create a new branch only when the current branch name does **not** already match the desired `{issue-number}-{section-name}` for the section.
   - If the current branch already matches the section, do not create or switch branches.
   - If the user explicitly requests to stay on the current branch, do not create a branch.
3. Research the codebase to gather information about the change.
4. Ask the user clarifying questions.
   - Clearly number the questions.
   - Clearly letter the options for each question.
5. Update the Task List section with any new updates based on your research and the user's answers.
6. Explain the current Task List section to the user.
7. When the Task List section is approved by the user, instruct the user to run `/do-pr` to begin implementing the changes. Do not use the word “proceed” as the final prompt — always reference `/do-pr`.
