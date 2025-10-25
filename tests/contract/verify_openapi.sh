#!/usr/bin/env bash
set -euo pipefail

PORT=8080
START_TIMEOUT=15

# Start the server in the background on the configured port
PORT=$PORT go run ./src &
PID=$!

# Ensure we kill the server process group on exit/interrupt
trap 'rc=$?; echo "Cleaning up server (PID=$PID)"; kill -- -"$PID" 2>/dev/null || kill "$PID" 2>/dev/null || true; wait "$PID" 2>/dev/null || true; exit $rc' INT TERM EXIT

# Wait for the server to be ready (health endpoint)
echo "Waiting up to ${START_TIMEOUT}s for server to become ready on port ${PORT}..."
for i in $(seq 1 $START_TIMEOUT); do
  if curl -sS --fail "http://localhost:${PORT}/api/v1/health" >/dev/null 2>&1; then
    echo "Server is ready."
    break
  fi
  sleep 1
  if [ "$i" -eq "$START_TIMEOUT" ]; then
    echo "Server did not become ready within ${START_TIMEOUT}s; aborting."
    exit 1
  fi
done

# Run Schemathesis (uvx wrapper)
uvx schemathesis run specs/openapi.yaml --url "http://localhost:${PORT}"
