# Code Review Instructions for Simple-Sync

This is a Go REST API project using Gin framework for event storage and access control.

## General Guidelines
- Follow Go best practices and idioms
- Write clean, readable, and maintainable code
- Use proper error handling throughout
- Ensure thread safety with appropriate mutex usage
- Follow existing code patterns and structure

## Pull Request Requirements
- **CHANGELOG.md MUST be updated** for every pull request that introduces user-facing changes
- Include version numbers and dates for releases
- Follow the existing CHANGELOG format with pull request links (NOT issue links)
- Keep entries less than five lines

## Project Structure
- `handlers/`: HTTP endpoint handlers
- `models/`: Data structures for events, users, ACL
- `middleware/`: Authentication and CORS middleware
- `storage/`: SQLite persistence layer
- `main.go`: Application entry point

## Technical Requirements
- Use Gin web framework for HTTP handling
- Implement JWT authentication
- Use SQLite for data storage
- Support CORS for web clients
- Return appropriate HTTP status codes
- Handle malformed requests gracefully

## Code Quality
- Use proper JSON serialization/deserialization
- Implement proper logging
- Add unit tests for new functionality
- Ensure data validation
- Follow API specifications in `docs/api.md`

## Security
- Never expose sensitive information
- Validate all inputs
- Use secure JWT handling
- Implement proper access control via ACL

## Performance
- Optimize for concurrent requests
- Use efficient data structures
- Minimize file I/O operations where possible
