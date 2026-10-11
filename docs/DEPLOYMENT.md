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
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
docker compose -f docker-compose.yml -f docker-compose.prod.yml ps
curl -s http://localhost/healthz           # {"status":"ok",...}
```

Migrations run automatically on backend startup (`RUN_MIGRATIONS=true`).

## 3. Updating

Normal updates are deployed by GitHub Actions after CI passes on `main`
(ADR-0011). The VM pulls images; it does not compile.

Manual fallback on the VM:

```bash
cd /opt/family-tree
git fetch --depth 1 origin main && git reset --hard origin/main
docker compose -f docker-compose.yml -f docker-compose.prod.yml pull
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
docker image prune -f
```

Expect a few seconds of restart. `.env` secrets stay on the VM. Deploy only
rewrites the three `*_IMAGE` lines.

## 4. TLS (Let's Encrypt)

Production host: `twkakbar.my.id` → this VM. `www` is not configured.
Nginx listens on 443, redirects port 80 to HTTPS, and reads the certificate
from `/etc/letsencrypt` (mounted read-only). ACME challenges are served from
`/var/www/certbot`. Login cookies require `COOKIE_SECURE=true`.

First issue (stop nginx so standalone can bind port 80):

```bash
sudo apt-get install -y certbot
sudo mkdir -p /var/www/certbot
cd /opt/family-tree
docker compose -f docker-compose.yml -f docker-compose.prod.yml stop nginx
sudo certbot certonly --standalone -d twkakbar.my.id \
  --non-interactive --agree-tos -m twkakbar@gmail.com
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build nginx
```

Renewal uses webroot so nginx can stay up. Certbot's systemd timer calls
`certbot renew`. A deploy hook must reload nginx, which caches certificates
until restart:

```bash
# /etc/letsencrypt/renewal-hooks/deploy/restart-nginx.sh
#!/bin/sh
docker compose -f /opt/family-tree/docker-compose.yml -f /opt/family-tree/docker-compose.prod.yml restart nginx
```

Switch an existing standalone certificate to webroot after the HTTP
challenge location is live:

```bash
sudo certbot certonly --webroot -w /var/www/certbot -d twkakbar.my.id \
  --non-interactive --agree-tos -m twkakbar@gmail.com --keep-until-expiring
```

## 5. CI/CD deploys (GitHub Actions)

`deploy.yml` runs after the `CI` workflow succeeds on `main`, and can also
be started by hand. It builds images in Actions, pushes them to GHCR, then
SSHes into this VM (ADR-0011).

Images:

```text
ghcr.io/twkakbar/family-tree-backend:latest
ghcr.io/twkakbar/family-tree-frontend:latest
ghcr.io/twkakbar/family-tree-nginx:latest
```

Setup (already applied for `43.173.10.141`):

1. Repository **variables**: `DEPLOY_ENABLED=true`,
   `DEPLOY_PATH=/opt/family-tree`.
2. Repository **secrets**: `DEPLOY_HOST`, `DEPLOY_USER` (`ubuntu`),
   `DEPLOY_SSH_KEY` (private key). The matching public key is in
   `/home/ubuntu/.ssh/authorized_keys` on the VM. Do not put the VM
   password in GitHub.
3. The VM `.env` image names are rewritten on each deploy to the GHCR
   `:latest` tags above. Other secrets in `.env` are left alone.
4. Images stay private. The deploy job logs into GHCR with the Actions
   token for the pull, then logs out. Do not commit that token.

There is no Docker Hub token in this setup. Do not point `*_IMAGE` back at
Docker Hub unless a new ADR says so.

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
