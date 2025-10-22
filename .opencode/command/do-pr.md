---
description: Develop the changes for a new pull request
agent: build
---

If the user has NOT previously ran the /start-pr command, exit early and ask them to run that command.

1. Use the `todoread` tool to get the TODO list for the current issue
1. For each item in the TODO list:
   1. Carefully implement change
   1. Run relevant automated tests
   1. Refine based on test results
   1. Check off corresponding item in the issue spec
   1. Commit
1. After all items on the TODO list are complete:
   1. Push
   1. Create pull request
   1. Update changelog
   1. Push
