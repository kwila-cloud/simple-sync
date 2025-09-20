# Changelog

## [0.1.0] - 2025-09-20
- [#13](https://github.com/kwila-cloud/simple-sync/pull/13): Fixed GitHub workflow issues
  - Added contents write permission to version-update workflow for committing and tagging
  - Removed unnecessary Go cache step from CI workflow
  - Ensured Go version consistency across workflows
- [#11](https://github.com/kwila-cloud/simple-sync/pull/11): Enhanced test coverage for timestamp filtering and concurrency
  - Added TestGetEventsWithTimestampFiltering to verify filtering with actual data
  - Added TestConcurrentPostEvents to test thread safety with multiple goroutines
  - Fixed data race in concurrent test by adding proper mutex synchronization
  - Enabled race detection in CI pipeline
- [#9](https://github.com/kwila-cloud/simple-sync/pull/9): Initial MVP implementation with basic event storage REST API
  - Core Go application with Gin framework
  - GET/POST /events endpoints with timestamp filtering
  - Thread-safe in-memory storage
  - Contract and unit test suites
  - CI/CD pipeline with GitHub Actions
  - Complete project documentation and specifications
  - Docker development environment setup
