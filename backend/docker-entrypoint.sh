#!/bin/sh
# The uploads volume is created by Docker as root. Fix ownership, then
# drop to the app user before starting the API.
set -eu
mkdir -p /app/uploads
chown -R app:app /app/uploads
chmod 755 /app/uploads
exec su-exec app "$@"
