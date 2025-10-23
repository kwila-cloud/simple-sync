#!/usr/bin/env bash
set -euo pipefail

# Simple Sync restore script
# Usage: ./scripts/restore.sh <backup-file> [--stop]
# Restores the given backup file into ./data/simple-sync.db. If --stop is provided, the service will be stopped before restore and started after.

if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <backup-file> [--stop]"
  exit 2
fi

BACKUP_FILE="$1"
shift

STOP=false
if [ "${1:-}" = "--stop" ]; then
  STOP=true
fi

DB_DIR="./data"
DB_PATH="$DB_DIR/simple-sync.db"

if [ ! -f "$BACKUP_FILE" ]; then
  echo "Error: backup file not found: $BACKUP_FILE"
  exit 1
fi

mkdir -p "$DB_DIR"

if [ "$STOP" = true ]; then
  echo "Stopping simple-sync service..."
  docker compose stop simple-sync || true
fi

# Safety: move existing DB aside
if [ -f "$DB_PATH" ]; then
  OLD_TS=$(date +%Y%m%d%H%M%S)
  mv -- "$DB_PATH" "$DB_DIR/simple-sync.db.bak.$OLD_TS"
  echo "Moved existing DB to $DB_DIR/simple-sync.db.bak.$OLD_TS"
fi

cp -- "$BACKUP_FILE" "$DB_PATH"

if [ "$STOP" = true ]; then
  echo "Starting simple-sync service..."
  docker compose start simple-sync || true
fi

echo "Restore complete: $DB_PATH"
exit 0
