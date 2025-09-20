# Code Review Instructions for Simple-Sync

This is a Go REST API project using Gin framework for event storage and access control.

## General Guidelines
- Follow Go best practices and idioms
- Write clean, readable, and maintainable code
- Use proper error handling throughout
- Ensure thread safety with appropriate mutex usage
- Follow existing code patterns and structure

## Project Structure
- `handlers/`: HTTP endpoint handlers
- `models/`: Data structures for events, users, ACL
- `middleware/`: Authentication and CORS middleware
- `storage/`: File-based persistence layer
- `main.go`: Application entry point

## Technical Requirements
- Use Gin web framework for HTTP handling
- Implement JWT authentication
- Use file-based JSON storage for persistence
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