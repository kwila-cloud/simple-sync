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
2. Before making any repository changes (branch creation, writing TODOs), ask the user a short set of clarifying questions and require explicit answers. These MUST include at minimum:
   - Whether to create a new branch or stay on the current branch.
   - If creating a branch, confirm the branch name format: `{issue-number}-{section-name}`.
   - Whether to create the TODO list now or wait for additional instructions.
3. Branch creation rules:
   - Create a new branch only when the current branch name does **not** already match the desired `{issue-number}-{section-name}` for the section.
   - If the current branch already matches the section, do not create or switch branches.
   - If the user explicitly requests to stay on the current branch, do not create a branch.
4. After the user confirms the clarifying questions, then (and only then):
   - Create the branch if needed (see rule #3).
   - Research the codebase to gather implementation notes for the section.
   - Ask any additional clarifying questions that arise while researching (require answers before proceeding further).
5. Create a TODO list with the `todowrite` tool only after the user approved creating it.
6. Explain the TODO list to the user and allow them to request refinements.
7. When the plan/TODO list is approved by the user, instruct the user to run `/do-pr` to begin implementing the changes. Do not use the word “proceed” as the final prompt — always reference `/do-pr`.

Example usage and flows

- Typical create-branch flow:
  - Agent: "Next incomplete section: `Documentation and Configuration Updates`. Create branch `7-documentation-and-configuration-updates` and prepare TODOs? (Y/n)"
  - User: "Y"
  - Agent: Creates branch, researches, creates TODOs, explains plan, then: "When ready to implement, run `/do-pr`."

- Stay-on-current-branch flow:
  - Agent: "Next incomplete section: `Docs`. Current branch is `7-documentation-and-configuration-updates`. Create a new branch `7-docs` or stay on current branch? (create/stay)"
  - User: "stay"
  - Agent: Does not create a branch, asks whether to create TODOs now, and if confirmed, creates TODOs only.

Notes

- The command MUST NOT make repository modifications before receiving the user's answers to the clarifying questions.
- All user-facing prompts produced by this command must end by instructing the user to run `/do-pr` when they want the agent to start implementation.
