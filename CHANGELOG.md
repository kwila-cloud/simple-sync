# Release History

## [0.4.0] - unreleased
- [#60](https://github.com/kwila-cloud/simple-sync/pull/60): Implement SQLite setup token and API key storage
- [#59](https://github.com/kwila-cloud/simple-sync/pull/59): Implement SQLite ACL rule storage
- [#58](https://github.com/kwila-cloud/simple-sync/pull/58): Implement SQLite user storage
- [#57](https://github.com/kwila-cloud/simple-sync/pull/57): Implement SQLite event storage
- [#56](https://github.com/kwila-cloud/simple-sync/pull/56): Add database schema migrations
- [#55](https://github.com/kwila-cloud/simple-sync/pull/55): Add SQLite storage foundation (initialization, pragmas, pooling, and tests)
- [#61](https://github.com/kwila-cloud/simple-sync/pull/61): Add performance and concurrency tests for SQLite storage
- [#54](https://github.com/kwila-cloud/simple-sync/pull/54): Add model validation and database compatibility for SQLite storage
- [#53](https://github.com/kwila-cloud/simple-sync/pull/53): Add storage factory and error types
- [#52](https://github.com/kwila-cloud/simple-sync/pull/52): Refactor ACL service to use storage interface directly with proper error handling
- [#51](https://github.com/kwila-cloud/simple-sync/pull/51): Add ACL storage interface methods with comprehensive test coverage
- [#49](https://github.com/kwila-cloud/simple-sync/pull/49): Added data persistence implementation specification

## [0.3.0] - 2025-10-18
- [#45](https://github.com/kwila-cloud/simple-sync/pull/45): Implement ACL endpoint
- [#43](https://github.com/kwila-cloud/simple-sync/pull/43): Changed API authentication from Authorization: Bearer to X-API-Key header
- [#40](https://github.com/kwila-cloud/simple-sync/pull/40): Implemented ACL system
- [#37](https://github.com/kwila-cloud/simple-sync/pull/37): Replaced JWT authentication with API key system
- [#36](https://github.com/kwila-cloud/simple-sync/pull/36): More documentation improvements
- [#34](https://github.com/kwila-cloud/simple-sync/pull/34): Updated ACL documentation to event-based system with specificity evaluation
- [#32](https://github.com/kwila-cloud/simple-sync/pull/32): Fixed links on docs home page
- [#31](https://github.com/kwila-cloud/simple-sync/pull/31): Fixed docs deployment
- [#30](https://github.com/kwila-cloud/simple-sync/pull/30): Added docs site
- [#28](https://github.com/kwila-cloud/simple-sync/pull/28): Fixed healthcheck URL in docker compose file
- [#27](https://github.com/kwila-cloud/simple-sync/pull/27): Fixed image URL in docker compose file

## [0.2.0] - 2025-09-22
- [#25](https://github.com/kwila-cloud/simple-sync/pull/25): Added v1 API prefix for future versioning
- [#24](https://github.com/kwila-cloud/simple-sync/pull/24): Docker configuration for easy deployment
- [#15](https://github.com/kwila-cloud/simple-sync/pull/15): Basic JWT authentication

## [0.1.0] - 2025-09-20
- [#14](https://github.com/kwila-cloud/simple-sync/pull/14): Fixed GitHub release workflow access issue
- [#13](https://github.com/kwila-cloud/simple-sync/pull/13): Fixed GitHub workflow issues
- [#11](https://github.com/kwila-cloud/simple-sync/pull/11): Enhanced test coverage for timestamp filtering and concurrency
- [#9](https://github.com/kwila-cloud/simple-sync/pull/9): Initial MVP implementation with basic event storage REST API
