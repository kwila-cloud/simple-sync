# Build stage
FROM golang:1.25-bookworm AS builder

# Install build dependencies (needed for CGO + sqlite3)
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    ca-certificates \
    build-essential \
    libsqlite3-dev \
    pkg-config \
  && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy the full repository
COPY . .

# Build the application with CGO enabled (sqlite requires cgo)
RUN CGO_ENABLED=1 GOOS=linux go build -v -a -installsuffix cgo -o main ./src

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies (sqlite library, certificates, curl)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libsqlite3-0 \
    curl \
  && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -r appgroup && useradd -r -g appgroup appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Create data directory for SQLite database
RUN mkdir -p /app/data && chown -R appuser:appgroup /app/data

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
