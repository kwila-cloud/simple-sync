# Changelog

## [0.3.0] - unreleased
- [#28](https://github.com/kwila-cloud/simple-sync/pull/28): Fixed healthcheck URL in docker compose file
- [#27](https://github.com/kwila-cloud/simple-sync/pull/27): Fixed image URL in docker compose file

## [0.2.0] - 2025-09-22
- [#25](https://github.com/kwila-cloud/simple-sync/pull/25): Add v1 API prefix for future versioning
  - Updated all API endpoints to start with /api/v1/
- [#24](https://github.com/kwila-cloud/simple-sync/pull/24): Docker configuration for easy deployment
  - Added multi-stage Dockerfile with Go 1.25 and Alpine runtime
  - Created docker-compose.yml for local development with health checks
  - Implemented /health endpoint with service status, version, and uptime
  - Added environment configuration validation (JWT_SECRET, PORT, ENVIRONMENT)
  - Updated CI/CD pipeline for automated Docker image builds on GitHub Container Registry
- [#15](https://github.com/kwila-cloud/simple-sync/pull/15): Basic JWT authentication
  - Added JWT token generation, validation, and user authentication
  - Protected all /events endpoints with authentication

## [0.1.0] - 2025-09-20
- [#14](https://github.com/kwila-cloud/simple-sync/pull/14): Fixed GitHub release workflow access issue
- [#13](https://github.com/kwila-cloud/simple-sync/pull/13): Fixed GitHub workflow issues
  - Added contents write permission to version-update workflow for committing and tagging
  - Removed unnecessary Go cache step from CI workflow
- [#11](https://github.com/kwila-cloud/simple-sync/pull/11): Enhanced test coverage for timestamp filtering and concurrency
  - Added TestGetEventsWithTimestampFiltering to verify filtering with actual data
  - Added TestConcurrentPostEvents to test thread safety with multiple goroutines
  - Fixed data race in concurrent test by adding proper mutex synchronization
  - Enabled race detection in CI pipeline
- [#9](https://github.com/kwila-cloud/simple-sync/pull/9): Initial MVP implementation with basic event storage REST API
  - Core Go application with Gin framework
  - GET/POST /events endpoints with timestamp filtering
  - Contract and unit test suites
 - CI/CD pipeline with GitHub Actions
