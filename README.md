# Simple Sync
A simple sync system for local-first apps.

Built with [Go](https://go.dev/), [Gin](https://github.com/gin-gonic/gin), and [SQLite](https://www.sqlite.org/index.html). See the [Tech Stack](docs/tech-stack.md) document for details on the technologies used in this project and the rationale behind those choices.

**NOTE** - This project is in the alpha stage. Many of the things documented here and elsewhere in this repo do not actually exist yet.

## Documentation

📚 [View Documentation](https://kwila-cloud.github.io/simple-sync/) - Complete API reference, ACL system, and technical guides.

## Quick Start

To run Simple Sync using Docker Compose:

### Prerequisites
- Docker and Docker Compose installed
- Git (to clone the repository)

### Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/kwila-cloud/simple-sync.git
   cd simple-sync
   ```

2. Start the services:
   ```bash
   docker compose up -d
   ```

4. Verify the service is running:
   ```bash
    curl http://localhost:8080/api/v1/health
   ```
   You should see a JSON response with status "healthy".

### Optional: Add Frontend
To run with a frontend application, add it as an additional service in `docker-compose.yml`. For example:
```yaml
services:
  frontend:
    image: your-frontend-image
    ports:
      - "3000:3000"
    depends_on:
      - simple-sync
```

## Development

### Building

To build the application:

```bash
go build -o simple-sync ./src
```

This will make a `simple-sync` executable file.

### Running Locally

To run the application locally:

```bash
# Run the server
go run ./src

The server will start on port 8080 by default.

### Running Tests

To run the test suite:

Run unit, contract, and integration tests with race detection:
```bash
go test -race ./tests/unit ./tests/contract ./tests/integration
```

Run performance tests (without race detection):
```bash
go test ./tests/performance
```

This will run all tests including:
- **Contract tests** (`tests/contract/`) - API contract validation
- **Integration tests** (`tests/integration/`) - Full workflow testing
- **Unit tests** (`tests/unit/`) - Individual component testing
- **Performance tests** (`tests/performance/`) - Response time validation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

