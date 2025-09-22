# Build stage
FROM golang:1.25-alpine AS builder

# Install git and ca-certificates (needed for Go modules)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY src/ ./src/

    # Build the application
    RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src

# Runtime stage
FROM alpine:latest

# Install ca-certificates and curl for HTTPS requests and health checks
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

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