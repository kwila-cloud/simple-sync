#!/usr/bin/env bash
set -euo pipefail

PORT=8080
START_TIMEOUT=15

# Start the server in the background
PORT=$PORT go run ./src >/dev/null 2>&1 &
PID=$!
PGID=$(ps -o pgid= "$PID" | tr -d ' ' || echo "")

cleanup() {
  rc=${1:-0}
  trap - INT TERM EXIT
  if [ -n "$PGID" ]; then
    kill -- -"$PGID" 2>/dev/null || true
  fi
  kill "$PID" 2>/dev/null || true
  wait "$PID" 2>/dev/null || true
  exit $rc
}

trap 'cleanup $?' INT TERM EXIT

# Wait for readiness
for i in $(seq 1 $START_TIMEOUT); do
  if curl -sS --fail "http://localhost:${PORT}/api/v1/health" >/dev/null 2>&1; then
    break
  fi
  sleep 1
  if [ "$i" -eq "$START_TIMEOUT" ]; then
    cleanup 1
  fi
done

set +e
uvx schemathesis run specs/openapi.yaml --url "http://localhost:${PORT}"
SC=$?
set -e

cleanup $SC
