<!--
Sync Impact Report:
- Version change: none → 1.0.0
- List of modified principles: Initial principles added
- Added sections: Technology Stack, Development Workflow
- Removed sections: none
- Templates requiring updates: .specify/templates/plan-template.md (Constitution Check section), .specify/templates/tasks-template.md (version reference), .specify/templates/agent-file-template.md ([PROJECT NAME])
- Follow-up TODOs: none
-->
# simple-sync Constitution

## Core Principles

### I. RESTful API Design
All API endpoints MUST follow REST principles, using appropriate HTTP methods (GET, POST, PUT, DELETE) and status codes. Endpoints MUST be resource-oriented, stateless, and provide consistent JSON responses. Rationale: Ensures predictable and maintainable API interactions for clients.

### II. Event-Driven Architecture
The system MUST store and manage data as a sequence of timestamped events with user, item, and action metadata. All data operations MUST be append-only to the event history. Rationale: Supports local-first synchronization and audit trails.

### III. Authentication and Authorization
All endpoints except public ones MUST require JWT-based authentication. Access control MUST be enforced via ACL rules that define user permissions on items and actions. Rationale: Protects data integrity and user privacy in a multi-user environment.

### IV. Data Persistence
Data MUST be persisted to file-based storage (JSON format) for simplicity and portability. The system MUST maintain data integrity across restarts. Rationale: Aligns with local-first philosophy and reduces deployment complexity.

### V. Security and Access Control
ACL rules MUST be evaluated in order, with deny-by-default behavior. Wildcard support MUST be provided for flexible permission management. Rationale: Ensures fine-grained control over data access while maintaining security.

## Technology Stack
The project MUST use Go with Gin web framework, SQLite for data storage, JWT for authentication, and TOML for configuration. All dependencies MUST be justified for simplicity, performance, and maintainability. Rationale: Chosen stack optimizes for the project's goals of simple code and high maintainability.

## Development Workflow
Development MUST follow an issue-driven workflow using GitHub CLI for tracking. Features MUST be implemented incrementally with testing. Code MUST be committed with descriptive messages referencing issues. Rationale: Ensures structured progress and accountability.

## Governance
Amendments to this constitution require consensus among maintainers and MUST be documented with rationale. Versioning follows semantic rules: MAJOR for breaking changes, MINOR for additions, PATCH for clarifications. All changes MUST be reviewed for compliance. Rationale: Maintains project integrity and guides decision-making.

**Version**: 1.0.0 | **Ratified**: 2025-09-20 | **Last Amended**: 2025-09-20