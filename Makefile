.PHONY: lint test build clean

# Linting and formatting
lint:
	gofmt -l .
	go vet ./...

# Run tests
test:
	go test ./...

# Build the application
build:
	go build -o bin/simple-sync ./src

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	gofmt -w .