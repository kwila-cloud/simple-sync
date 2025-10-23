#!/usr/bin/env bash
set -euo pipefail

# Simple Sync backup script
# Usage: ./scripts/backup.sh [--stop] [path-to-db]
# If --stop is provided, the script will stop the docker-compose service before copying and restart it afterwards.

STOP=false
DB_PATH="./data/simple-sync.db"

if [ "${1:-}" = "--stop" ]; then
  STOP=true
  shift
fi

if [ "$#" -ge 1 ]; then
  DB_PATH="$1"
fi

BACKUP_DIR="./backups"
mkdir -p "$BACKUP_DIR"

if [ ! -f "$DB_PATH" ]; then
  echo "Error: database file not found at $DB_PATH"
  exit 1
fi

TIMESTAMP=$(date +%Y%m%d%H%M%S)
BACKUP_FILE="$BACKUP_DIR/simple-sync-$TIMESTAMP.db"

if [ "$STOP" = true ]; then
  echo "Stopping simple-sync service..."
  docker compose stop simple-sync || true
fi

echo "Copying $DB_PATH -> $BACKUP_FILE"
cp -- "$DB_PATH" "$BACKUP_FILE"

if [ "$STOP" = true ]; then
  echo "Starting simple-sync service..."
  docker compose start simple-sync || true
fi

echo "Backup complete: $BACKUP_FILE"
exit 0
