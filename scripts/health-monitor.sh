#!/usr/bin/env bash
# Lightweight uptime check for healthz + readyz.
#
# Cron example (every 5 minutes):
#   */5 * * * * BASE_URL=https://example.com /opt/family-tree/scripts/health-monitor.sh
#
# Exit 0 = healthy. Non-zero = probe failed (use with monitoring / alerting).

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost}"
BASE_URL="${BASE_URL%/}"
TIMEOUT_SECS="${TIMEOUT_SECS:-10}"

check() {
  local path="$1"
  local body
  body="$(curl -fsS --max-time "${TIMEOUT_SECS}" "${BASE_URL}${path}")"
  echo "${body}" | grep -q '"status":"ok"'
}

echo "[$(date -Is)] probing ${BASE_URL}…"
check /healthz
check /readyz
echo "[$(date -Is)] healthy"
