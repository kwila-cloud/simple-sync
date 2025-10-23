# Simple Sync
A simple sync system for local-first apps.

Built with [Go](https://go.dev/), [Gin](https://github.com/gin-gonic/gin), and [SQLite](https://www.sqlite.org/index.html).

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

3. Verify the service is running:
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

### Persistent Data (Docker Compose)
The Docker Compose configuration mounts a local `./data` directory into the container (`./data:/app/data`) by default. This is the recommended setup for development and simple deployments because it keeps the SQLite database file on the host, making backups and inspection straightforward.

If you prefer Docker-managed storage, a named volume `simple-sync-data` is declared in `docker-compose.yml`; you can switch to it by uncommenting the named volume line and removing the `./data` bind mount.

Backup and restore helper scripts are provided in `./scripts`:
- `./scripts/backup.sh [--stop] [path-to-db]` — copy the DB file to `./backups/` (use `--stop` to stop the container during backup)
- `./scripts/restore.sh <backup-file> [--stop]` — restore a backup into `./data/simple-sync.db` (moves the existing DB aside first)

Example (take a backup and then start):

```bash
# create a backup (stop service during the copy)
./scripts/backup.sh --stop
# start services
docker compose up -d
```

Developer note: the app uses `github.com/mattn/go-sqlite3` which requires `libsqlite3-dev` and `CGO_ENABLED=1` when building locally or in CI.


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
```

### Database configuration (SQLite)

The application uses SQLite for persistent storage. Configure the database path with the `DB_PATH` environment variable. Defaults to `./data/simple-sync.db`.

Notes:
- The project currently uses `github.com/mattn/go-sqlite3`, which requires a C toolchain and the system SQLite development headers (`libsqlite3-dev`) to build. Ensure `CGO_ENABLED=1` when building a release binary.
- For testing, the code uses an in-memory SQLite database (`file::memory:?cache=shared`) where appropriate.

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

