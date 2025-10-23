#!/usr/bin/env bash
set -euo pipefail

# Simple Sync backup script
# Usage: ./scripts/backup.sh [--stop] [--dir <backup-dir>] [path-to-db]
# If --stop is provided, the script will stop the docker-compose service before copying and restart it afterwards.
# If --dir is provided, it specifies the directory to store backups (defaults to ./backups).

STOP=false
BACKUP_DIR="./backups"
DB_PATH="./data/simple-sync.db"
DB_PATH_PROVIDED=false

while [ "$#" -gt 0 ]; do
  case "$1" in
    --stop)
      STOP=true
      shift
      ;;
    --dir)
      if [ -z "${2:-}" ]; then
        echo "Error: --dir requires a directory argument"
        exit 2
      fi
      BACKUP_DIR="$2"
      shift 2
      ;;
    --dir=*)
      BACKUP_DIR="${1#--dir=}"
      shift
      ;;
    --help|-h)
      echo "Usage: $0 [--stop] [--dir <backup-dir>] [path-to-db]"
      exit 0
      ;;
    *)
      if [ "$DB_PATH_PROVIDED" = false ]; then
        DB_PATH="$1"
        DB_PATH_PROVIDED=true
        shift
      else
        echo "Unknown argument: $1"
        exit 2
      fi
      ;;
  esac
done

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
