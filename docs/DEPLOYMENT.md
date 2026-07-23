# Deployment Runbook — Linux VM

Target topology (ARCHITECTURE.md §7): one VM, Docker Compose, nginx as the
only published port, images built locally or pulled from Docker Hub
(`twkakbar/family-tree-*`).

## 1. Provision the VM

Any modern Debian/Ubuntu VM (1 vCPU / 1 GB RAM is enough to start).

```bash
# as root or with sudo
apt-get update && apt-get install -y ca-certificates curl git
curl -fsSL https://get.docker.com | sh          # Docker Engine + Compose v2
useradd -m -s /bin/bash deploy && usermod -aG docker deploy
```

Firewall: allow 22 (SSH), 80 (HTTP), 443 (HTTPS). Nothing else — Postgres
and the apps live on the internal Docker network.

## 2. First deployment

```bash
sudo mkdir -p /opt/family-tree && sudo chown deploy: /opt/family-tree
su - deploy
git clone https://github.com/twkakbar/family-tree.git /opt/family-tree && cd /opt/family-tree

cp .env.example .env
nano .env
```

Production `.env` checklist:

- **Remove** (or comment) the `COMPOSE_FILE` / `COMPOSE_PATH_SEPARATOR`
  lines — plain `docker compose` must use only `docker-compose.yml`.
- `POSTGRES_PASSWORD`: strong random (`openssl rand -base64 32`).
- `JWT_SECRET`: strong random, ≥ 32 chars (`openssl rand -base64 48`).
- `COOKIE_SECURE=true` once TLS is enabled (step 4).
- Optionally `ADMIN_EMAIL` / `ADMIN_PASSWORD` for the bootstrap admin.

Start:

```bash
docker compose -f docker-compose.yml up -d --build
docker compose ps                          # all services healthy?
curl -s http://localhost/healthz           # {"status":"ok",...}
```

Migrations run automatically on backend startup (`RUN_MIGRATIONS=true`).

## 3. Updating

```bash
cd /opt/family-tree
git pull
docker compose -f docker-compose.yml up -d --build
docker image prune -f
```

Zero-ish downtime is not a goal (CONTEXT.md constraints); expect a few
seconds of restart.

## 4. TLS (Let's Encrypt)

1. Point your domain's A record at the VM.
2. Obtain certificates on the host:
   `apt-get install -y certbot && certbot certonly --standalone -d example.com`
   (stop nginx briefly, or use the webroot method later).
3. In `docker-compose.yml` (nginx service): publish `443:443` and mount
   `/etc/letsencrypt:/etc/letsencrypt:ro`.
4. In `docker/nginx/conf.d/default.conf`: enable the commented TLS server
   block, set `server_name`, and redirect HTTP→HTTPS.
5. Set `COOKIE_SECURE=true` in `.env`, then
   `docker compose -f docker-compose.yml up -d --build`.
6. Auto-renewal: certbot installs a systemd timer; add a deploy hook that
   runs `docker compose -f /opt/family-tree/docker-compose.yml restart nginx`.

## 5. CI/CD deploys (GitHub Actions)

`deploy.yml` pushes images to Docker Hub on `v*` tags and can deploy over
SSH.

Setup:

1. Repository **secrets** for the registry: `DOCKERHUB_USERNAME`
   (`twkakbar`) and `DOCKERHUB_TOKEN` (Docker Hub → Account Settings →
   Personal access tokens, with Read/Write scope).
2. Repository **variables** to enable deploys: `DEPLOY_ENABLED=true`,
   `DEPLOY_PATH=/opt/family-tree`.
3. Repository **secrets** for SSH: `DEPLOY_HOST`, `DEPLOY_USER` (e.g.
   `deploy`), `DEPLOY_SSH_KEY` (private key whose public half is in the
   VM's `~/.ssh/authorized_keys`), optional `DEPLOY_PORT`.
4. On the VM, switch `.env` to registry images so `docker compose pull`
   works:

```bash
BACKEND_IMAGE=twkakbar/family-tree-backend:latest
FRONTEND_IMAGE=twkakbar/family-tree-frontend:latest
NGINX_IMAGE=twkakbar/family-tree-nginx:latest
```

   If the Docker Hub repositories are private, run a one-time
   `docker login -u twkakbar` on the VM.
5. Release: `git tag v1.0.0 && git push --tags`.

## 6. Backups

Use the repo script (keeps dump + uploads + retention in one place):

```bash
# on the VM
sudo cp /opt/family-tree/scripts/backup.sh /opt/family-tree/backup.sh
sudo chmod +x /opt/family-tree/backup.sh
# optional: APP_DIR=/opt/family-tree BACKUP_DIR=/var/backups/family-tree
```

`crontab -e`:

```cron
0 3 * * * /opt/family-tree/backup.sh >> /var/log/family-tree-backup.log 2>&1
```

Copy `backups/` off the VM regularly (object storage, rsync, …).

Restore DB:

```bash
gunzip -c backups/familytree-YYYY-MM-DD.sql.gz \
  | docker compose -f docker-compose.yml exec -T postgres \
      psql -U <user> <db>
```

## 7. Monitoring

Probes (also used by Docker healthchecks):

| URL | Meaning |
|-----|---------|
| `GET /healthz` | Process up (`{"status":"ok","version":…}`) |
| `GET /readyz` | DB reachable (`{"status":"ok","version":…}`) |

Cron / external uptime (every 5 minutes):

```cron
*/5 * * * * BASE_URL=https://example.com /opt/family-tree/scripts/health-monitor.sh \
  >> /var/log/family-tree-health.log 2>&1
```

Point an uptime service (UptimeRobot, Better Stack, …) at `/readyz`.
Alert on non-200 or missing `"status":"ok"`.

## 8. Security headers (Phase 7)

Nginx applies CSP and related headers via
`docker/nginx/conf.d/security-headers.inc` (ADR-0010). After enabling TLS,
add HSTS (+ optional `upgrade-insecure-requests`) in the TLS server block
template in `default.conf`.

## 9. Operations cheat sheet

| Task | Command |
|------|---------|
| Logs (all / one service) | `docker compose logs -f --tail=100 [backend]` |
| Service status + health | `docker compose ps` |
| Manual health probe | `make health-check` / `BASE_URL=… bash scripts/health-monitor.sh` |
| Manual backup | `make backup` / `bash scripts/backup.sh` |
| Restart one service | `docker compose restart backend` |
| DB shell | `docker compose exec postgres psql -U familytree familytree` |
| Migration status | `docker compose exec backend /app/migrate version` |
| Roll back last migration | `docker compose exec backend /app/migrate down 1` |
| Disk usage cleanup | `docker system prune -f` |

## Troubleshooting

| Symptom | Likely cause / fix |
|---------|--------------------|
| `502` from nginx | Backend/frontend unhealthy → `docker compose ps`, `logs backend`. After container recreation nginx usually re-resolves (10 s DNS TTL); worst case `docker compose restart nginx`. |
| Backend exits: `JWT_SECRET is required…` | `.env` missing/short secret (min 32 chars). |
| Backend exits: `database is dirty` | Failed migration → see docs/DATABASE.md dirty-state recovery. |
| Login works, session lost after refresh | `COOKIE_SECURE=true` while serving plain HTTP — browsers drop the cookie. Enable TLS or set it to false temporarily. |
| Uploads 404 in production | `uploads_data` volume not mounted in nginx, or backend wrote outside `/app/uploads`. |
| Blank page / CSP console errors | See ADR-0010; report blocked directives before tightening further. |
