# Kesultanan Aceh Family Tree

A public, self-hosted web app for exploring the **Aceh royal family tree**.
Anyone can browse the lineage in an immersive 3D viewer; approved editors
add members in a private studio. Built with Go and Next.js.

This is **one** canonical tree for the whole site — not a multi-family
genealogy SaaS.

> **AI-assisted repository.** Start with [`AGENTS.md`](AGENTS.md) (agent
> rules & status), then [`CONTEXT.md`](CONTEXT.md) (product),
> [`ARCHITECTURE.md`](ARCHITECTURE.md) (design) and
> [`DEVELOPMENT_RULES.md`](DEVELOPMENT_RULES.md) (standards).

## What you get

| Surface | Purpose |
|---------|---------|
| `/` | Public 3D explore (portrait medallions, sea-depth backdrop) |
| `/about` | Brand / story |
| `/studio` | Editors: add members, photos, relationships, JSON backup |
| `/request-access` | Descendants ask for edit rights |
| `/admin/requests` · `/admin/users` | Super-admin approvals and accounts |
| `/login` · `/register` · forgot/reset password | Editor auth (not marketed on the public tree) |

## Tech stack

| Area | Technology |
|------|-----------|
| Backend | Go 1.25, Gin, GORM, PostgreSQL 16, golang-migrate, JWT, slog |
| Frontend | Next.js 15 (App Router), TypeScript, TailwindCSS v4, React Query v5, React Three Fiber, Axios, Zod, React Hook Form |
| Infrastructure | Docker, Docker Compose, Nginx reverse proxy, GitHub Actions, single-VM deployment |

## Quickstart (development)

Prerequisites: Docker Desktop (or Docker Engine + Compose v2), Git.

```bash
git clone https://github.com/twkakbar/family-tree.git
cd family-tree
cp .env.example .env      # defaults work for local development
# optional: set ADMIN_EMAIL + ADMIN_PASSWORD for a bootstrap admin
docker compose up --build
```

The `.env` file sets `COMPOSE_FILE=docker-compose.yml;docker-compose.dev.yml`,
so plain `docker compose up` starts the **development** stack:

| URL | Service |
|-----|---------|
| http://localhost | Full app through Nginx (public tree + API) |
| http://localhost:3000 | Next.js dev server (HMR) |
| http://localhost:8080/healthz | Go API direct |
| localhost:5432 | PostgreSQL (`familytree` / password from `.env`) |

Database migrations run automatically on backend startup
(`RUN_MIGRATIONS=true`). The royal family is seeded on first boot.

### Hybrid development (Go/Node natively, Postgres in Docker)

```bash
docker compose up postgres           # database only
cd backend  && cp .env.example .env  && go run ./cmd/api
cd frontend && cp .env.example .env.local && npm install && npm run dev
```

## Common commands

With GNU make: `make help`. Raw equivalents:

| Task | Command |
|------|---------|
| Dev stack | `docker compose up --build` |
| Stop dev stack | `docker compose down` |
| Backend tests | `cd backend && go test -race -cover ./...` |
| Backend lint | `cd backend && go vet ./... && gofmt -l .` |
| Migrations | `cd backend && go run ./cmd/migrate up\|version\|new <name>` |
| Frontend checks | `cd frontend && npm run lint && npm run typecheck` |
| Health probe | `make health-check` / `bash scripts/health-monitor.sh` |
| Backup (Compose) | `make backup` / `bash scripts/backup.sh` |
| E2E smoke | `make e2e-install` then `BASE_URL=http://127.0.0.1 make e2e` |

## Production (Linux VM)

```bash
cp .env.example .env   # then: remove the COMPOSE_FILE lines, set real secrets
docker compose -f docker-compose.yml up -d --build
```

Full runbook (provisioning, TLS, backups, monitoring, CSP, CI/CD):
[`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md).

## Repository layout

```
backend/    Go API — Clean Architecture (cmd, internal/{domain,service,…}, pkg, migrations)
frontend/   Next.js app — src/{app,components,hooks,services,lib,types,utils}
docker/     Nginx reverse proxy image (CSP + security headers)
e2e/        Playwright smoke tests against Compose
scripts/    backup.sh, health-monitor.sh
docs/       API reference, database docs, deployment runbook, ADRs
```

## Project status

**Phases 1–7 complete** (foundation → Aceh pivot → photos/search →
export/import → admin/password reset → hardening).

Released: `v0.1.0` images on Docker Hub; SSH deploy is optional until a VM
is configured.

Deferred: SMTP email, GEDCOM, nonce-strict CSP.

Roadmap & feature table: [`CONTEXT.md`](CONTEXT.md).  
Implementation checklist: [`AGENTS.md`](AGENTS.md).

## Documentation

- Product context: [`CONTEXT.md`](CONTEXT.md)
- Architecture: [`ARCHITECTURE.md`](ARCHITECTURE.md)
- API reference: [`docs/api/`](docs/api/README.md)
- Database schema & migrations: [`docs/DATABASE.md`](docs/DATABASE.md)
- Deployment & ops: [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md)
- Architecture decisions: [`docs/adr/`](docs/adr/README.md)
