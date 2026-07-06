# Family Tree

A modern, self-hosted web application for creating and managing family
trees: families, members, relationships, interactive tree views, sharing,
and import/export — built with Go and Next.js.

> **AI-assisted repository.** Start with [`AGENTS.md`](AGENTS.md) (agent
> rules & status), then [`CONTEXT.md`](CONTEXT.md) (product),
> [`ARCHITECTURE.md`](ARCHITECTURE.md) (design) and
> [`DEVELOPMENT_RULES.md`](DEVELOPMENT_RULES.md) (standards).

## Tech stack

| Area | Technology |
|------|-----------|
| Backend | Go 1.25, Gin, GORM, PostgreSQL 16, golang-migrate, JWT, slog |
| Frontend | Next.js 15 (App Router), TypeScript, TailwindCSS v4, React Query v5, Axios, Zod, React Hook Form |
| Infrastructure | Docker, Docker Compose, Nginx reverse proxy, GitHub Actions, single-VM deployment |

## Quickstart (development)

Prerequisites: Docker Desktop (or Docker Engine + Compose v2), Git.

```bash
git clone https://github.com/twkakbar/family-tree.git
cd family-tree
cp .env.example .env      # defaults work for local development
docker compose up --build
```

The `.env` file sets `COMPOSE_FILE=docker-compose.yml;docker-compose.dev.yml`,
so plain `docker compose up` starts the **development** stack:

| URL | Service |
|-----|---------|
| http://localhost | Full app through Nginx (production-like) |
| http://localhost:3000 | Next.js dev server (HMR) |
| http://localhost:8080/healthz | Go API direct |
| localhost:5432 | PostgreSQL (`familytree` / password from `.env`) |

Database migrations run automatically on backend startup
(`RUN_MIGRATIONS=true`). To get an admin account, set `ADMIN_EMAIL` and
`ADMIN_PASSWORD` in `.env` before the first start — or just register a
normal user via the UI.

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
| Frontend dev | `cd frontend && npm run dev` |
| Frontend checks | `cd frontend && npm run lint && npm run typecheck` |

## Production (Linux VM)

```bash
cp .env.example .env   # then: remove the COMPOSE_FILE lines, set real secrets
docker compose -f docker-compose.yml up -d --build
```

Full runbook (provisioning, TLS, backups, CI/CD deploys):
[`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md).

## Repository layout

```
backend/    Go API — Clean Architecture (cmd, internal/{domain,service,…}, pkg, migrations)
frontend/   Next.js app — src/{app,components,hooks,services,lib,types,utils}
docker/     Nginx reverse proxy image
docs/       API reference, database docs, deployment runbook, ADRs
```

## Project status

Phase 1 (foundation: architecture, Docker, CI/CD, database schema, full
authentication slice, frontend shell) is complete. See the roadmap in
[`CONTEXT.md`](CONTEXT.md) and the status section in [`AGENTS.md`](AGENTS.md).

## Documentation

- API reference: [`docs/api/`](docs/api/README.md)
- Database schema & migrations: [`docs/DATABASE.md`](docs/DATABASE.md)
- Deployment runbook: [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md)
- Architecture decisions: [`docs/adr/`](docs/adr/README.md)
