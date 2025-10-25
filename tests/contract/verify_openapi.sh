#!/usr/bin/env bash
set -euo pipefail

PORT=8080
START_TIMEOUT=15

# Start the server in the background on the configured port in a new session
# Use setsid so the server becomes leader of a new process group we can kill reliably
if command -v setsid >/dev/null 2>&1; then
  PORT=$PORT setsid go run ./src &
else
  PORT=$PORT go run ./src &
fi
PID=$!

# Capture the process group id (may equal PID if started with setsid)
PGID=$(ps -o pgid= "$PID" | tr -d ' ' || true)

cleanup() {
  rc=${1:-0}
  echo "Cleaning up server (PID=$PID, PGID=$PGID)"
  # Disable traps while cleaning up to avoid recursion
  trap - INT TERM EXIT
  if [ -n "$PGID" ] && [ "$PGID" -gt 0 ] 2>/dev/null; then
    kill -- -"$PGID" 2>/dev/null || true
  fi
  kill "$PID" 2>/dev/null || true
  # Kill any child processes of PID as a fallback
  if command -v pkill >/dev/null 2>&1; then
    pkill -P "$PID" 2>/dev/null || true
  fi
  wait "$PID" 2>/dev/null || true
  return $rc
}

# Ensure we clean up on exit/interrupt
trap 'cleanup $?' INT TERM EXIT

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
    cleanup 1
    exit 1
  fi
done

# Run Schemathesis (uvx wrapper)
set +e
uvx schemathesis run specs/openapi.yaml --url "http://localhost:${PORT}"
SC_RESULT=$?
set -e

# Run cleanup explicitly and exit with schemathesis' exit code
cleanup $SC_RESULT
exit $SC_RESULT
