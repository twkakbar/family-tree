#!/usr/bin/env bash
# Backup Postgres + uploads for the Family Tree stack.
#
# Install on the VM:
#   sudo cp scripts/backup.sh /opt/family-tree/backup.sh
#   sudo chmod +x /opt/family-tree/backup.sh
#   crontab -e  →  0 3 * * * /opt/family-tree/backup.sh >> /var/log/family-tree-backup.log 2>&1
#
# Override paths with env vars if needed:
#   APP_DIR=/opt/family-tree BACKUP_DIR=/opt/family-tree/backups RETENTION_DAYS=14

set -euo pipefail

APP_DIR="${APP_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
BACKUP_DIR="${BACKUP_DIR:-${APP_DIR}/backups}"
RETENTION_DAYS="${RETENTION_DAYS:-14}"
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yml}"
STAMP="$(date +%F)"

cd "${APP_DIR}"
mkdir -p "${BACKUP_DIR}"

if [[ ! -f .env ]]; then
  echo "error: ${APP_DIR}/.env not found" >&2
  exit 1
fi

# shellcheck disable=SC1091
set -a
# Read only the keys we need (avoid sourcing unknown shell syntax).
POSTGRES_USER="$(grep -E '^POSTGRES_USER=' .env | head -n1 | cut -d= -f2-)"
POSTGRES_DB="$(grep -E '^POSTGRES_DB=' .env | head -n1 | cut -d= -f2-)"
set +a

if [[ -z "${POSTGRES_USER}" || -z "${POSTGRES_DB}" ]]; then
  echo "error: POSTGRES_USER / POSTGRES_DB missing from .env" >&2
  exit 1
fi

echo "[$(date -Is)] dumping database ${POSTGRES_DB}…"
docker compose -f "${COMPOSE_FILE}" exec -T postgres \
  pg_dump -U "${POSTGRES_USER}" "${POSTGRES_DB}" \
  | gzip > "${BACKUP_DIR}/familytree-${STAMP}.sql.gz"

echo "[$(date -Is)] archiving uploads volume…"
UPLOADS_MOUNT="$(
  docker volume inspect family-tree_uploads_data -f '{{.Mountpoint}}' 2>/dev/null || true
)"
if [[ -n "${UPLOADS_MOUNT}" && -d "${UPLOADS_MOUNT}" ]]; then
  tar -czf "${BACKUP_DIR}/uploads-${STAMP}.tar.gz" -C "${UPLOADS_MOUNT}" .
else
  echo "warning: uploads volume mountpoint not found; skipping uploads archive" >&2
fi

echo "[$(date -Is)] pruning backups older than ${RETENTION_DAYS} days…"
find "${BACKUP_DIR}" -type f \( -name 'familytree-*.sql.gz' -o -name 'uploads-*.tar.gz' \) \
  -mtime "+${RETENTION_DAYS}" -delete

echo "[$(date -Is)] backup complete → ${BACKUP_DIR}"
ls -lh "${BACKUP_DIR}" | tail -n +1
