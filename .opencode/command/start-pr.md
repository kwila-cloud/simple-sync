---
description: Start a new pull request by checking specs and creating the appropriate branch
agent: build
---

Automates the process of starting work on a new pull request:

1. Ask user for the issue number
1. Check the spec for that issue in `specs/`
1. Determine the next incomplete section from the Task List
1. Create branch using format `{issue-number}-{section-name}`
1. Research codebase to understand what's needed for the section
1. Ask the user any clarifying questions needed to implement the section
1. Explain the plan to the user for completing the section
1. Refine plan based on user feedback
1. Start implementing the plan when the user approves

Example usage:
- User: `/start-pr`
- Agent: "What issue number would you like to work on?"
- User: "7"
- Agent: Checks `specs/7-data-persistence.md`, finds next incomplete section, creates branch like `7-storage-interface-updates`, researches requirements, explains plan
