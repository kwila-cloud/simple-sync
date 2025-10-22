---
description: Start a new pull request by checking specs and creating the appropriate branch
agent: build
---

Automates the process of starting work on a new pull request for the given issue.

The user gave the input: "$ARGUMENTS"

Use the user input as the issue number.

If the user input is empty or invalid, prompt the user for the issue number.

1. Check the spec for the given issue in `specs/`
1. Determine the next incomplete section from the Task List
1. Create branch using format `{issue-number}-{section-name}`
   - DO NOT create branch if already in the correct branch for the issue number and section name
1. Research codebase to understand what's needed for the section
1. Ask the user any clarifying questions needed to implement the section
1. Create a TODO list with the `todowrite` tool
1. Explain the TODO list to the user 
1. Refine TODO list based on user feedback
1. Prompt the user to run `/do-pr` when they are ready

Example usage:
- User: `/start-pr 7`
- Agent: Checks `specs/7-data-persistence.md`, finds next incomplete section, creates branch like `7-storage-interface-updates`, researches requirements, explains plan
