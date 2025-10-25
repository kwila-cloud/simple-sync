#!/usr/bin/env bash
set -euo pipefail

PORT=8080
# Start the server in the background on the configured port
PORT=$PORT go run ./src &
PID=$!
trap 'kill "$PID" 2>/dev/null || true' EXIT

uvx schemathesis run specs/openapi.yaml --url http://localhost:8080
