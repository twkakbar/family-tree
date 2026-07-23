# AGENTS.md — Instructions for AI Agents

> **READ THIS FILE BEFORE MODIFYING ANY CODE.** It is the entry point for every
> AI session. Conversation context is lost between sessions; this repository's
> documentation is the single source of truth.

## Required reading order

| Order | File | What it tells you | Read when |
|-------|------|-------------------|-----------|
| 1 | `AGENTS.md` (this file) | Rules, repo map, current status | Always |
| 2 | `CONTEXT.md` | Business domain, features, roadmap, glossary | Always |
| 3 | `ARCHITECTURE.md` | System design, layers, data flow, auth, deployment | Before structural changes |
| 4 | `DEVELOPMENT_RULES.md` | Coding standards, layer rules, definition of done | Before writing code |
| 5 | `docs/` | API contracts, database schema, deployment, ADRs | When touching that area |

## Non-negotiable rules

1. **Never violate `DEVELOPMENT_RULES.md`.** If a requested change conflicts
   with it, say so and propose a compliant alternative.
2. **Never rename or move folders** unless explicitly requested.
3. **Never introduce a new library without written justification** (what it
   does, why the stdlib/existing deps are insufficient). Record it in the PR
   description or the relevant doc.
4. **Never hardcode secrets or configuration.** Everything configurable goes
   through environment variables (see `.env.example` files).
5. **Never mix layers.** Business logic lives in `backend/internal/service`,
   never in HTTP handlers or repositories. Frontend components never call
   axios directly — they use hooks → services → `lib/api-client.ts`.
6. **Never edit an applied database migration.** Schema changes = new
   migration files in `backend/migrations/`. GORM `AutoMigrate` is forbidden.
7. **Keep the architecture clean**: SOLID, Clean Architecture, no duplicate
   code, reusable components, separation of concerns.
8. **Documentation is part of the change.** See the update triggers below.
9. **Database conventions**: UUID primary keys, `timestamptz` timestamps,
   soft delete (`deleted_at`) for user-facing entities, snake_case names.
10. **API conventions**: all endpoints under `/api/v1`, JSON envelope
    `{success, data, error, meta}` (see `docs/api/README.md`).

## Documentation update triggers

| If you change… | You MUST update… |
|----------------|------------------|
| Features, scope, roadmap, domain concepts | `CONTEXT.md` |
| Architecture, layers, infra, auth design | `ARCHITECTURE.md` (+ an ADR in `docs/adr/` for significant decisions) |
| Coding standards, workflow, tooling | `DEVELOPMENT_RULES.md` |
| API endpoints (add/modify) | `docs/api/*` |
| Database schema (new migration) | `docs/DATABASE.md` |
| Deployment, CI/CD | `docs/DEPLOYMENT.md`, workflow comments |
| Implementation progress | **Status section at the bottom of this file** |

## Repository map

```
family-tree/
├── AGENTS.md                  ← you are here
├── CONTEXT.md                 ← business context, features, roadmap
├── ARCHITECTURE.md            ← system design (backend, frontend, DB, deploy)
├── DEVELOPMENT_RULES.md       ← coding standards and project rules
├── README.md                  ← human-facing quickstart
├── Makefile                   ← common commands (Linux/macOS/CI)
├── .env.example               ← compose-level env template (copy to .env)
├── docker-compose.yml         ← PRODUCTION stack (base)
├── docker-compose.dev.yml     ← DEV overrides (hot reload, exposed ports)
├── docker/
│   └── nginx/                 ← reverse proxy image (routes /, /api, /uploads)
├── .github/workflows/         ← ci.yml (lint+test+build), deploy.yml (Docker Hub+SSH)
├── scripts/                   ← backup.sh, health-monitor.sh (Phase 7 ops)
├── e2e/                       ← Playwright smoke tests against Compose
├── docs/
│   ├── api/                   ← API reference (implemented + planned contracts)
│   ├── adr/                   ← Architecture Decision Records
│   ├── DATABASE.md            ← ERD, schema rationale, migration how-to
│   └── DEPLOYMENT.md          ← VM provisioning, TLS, backups, CD
├── backend/                   ← Go + Gin + GORM + PostgreSQL (Clean Architecture)
│   ├── cmd/api/               ← entry point (DI wiring, graceful shutdown)
│   ├── cmd/migrate/           ← migration CLI (up/down/version/new)
│   ├── internal/domain/       ← entities, repo+service interfaces, domain errors
│   ├── internal/repository/   ← GORM implementations, DB connection, migrator
│   ├── internal/service/      ← business logic (the ONLY place for it)
│   ├── internal/handler/      ← Gin handlers, DTOs, error→HTTP mapping
│   ├── internal/middleware/   ← auth, CORS, logging, request-id, recovery…
│   ├── internal/routes/       ← route registration
│   ├── internal/config/       ← env-based configuration
│   ├── internal/validation/   ← custom validator rules
│   ├── pkg/                   ← reusable, domain-agnostic packages
│   └── migrations/            ← golang-migrate SQL files (embedded)
└── frontend/                  ← Next.js App Router + TypeScript + Tailwind
    └── src/
        ├── app/               ← routes: (auth) login/register, (dashboard)
        ├── components/        ← ui/ (design system), layout/
        ├── hooks/             ← React Query hooks (use-auth, …)
        ├── services/          ← API service layer (the ONLY axios callers)
        ├── lib/               ← api-client, query-client, validations, env
        ├── types/             ← shared TS types (API envelope, entities)
        └── utils/             ← small pure helpers (cn, …)
```

## Tech stack (do not deviate without an ADR)

- **Backend**: Go 1.25, Gin, GORM, PostgreSQL 16, golang-migrate, JWT
  (`golang-jwt/v5`), bcrypt, slog (stdlib logging).
- **Frontend**: Next.js 15 (App Router), TypeScript (strict), TailwindCSS v4,
  TanStack React Query v5, Axios, Zod, React Hook Form.
- **Infra**: Docker multi-stage builds, Docker Compose, Nginx reverse proxy,
  GitHub Actions, single-VM deployment.

## Commands cheat sheet

```bash
# Dev stack (after: cp .env.example .env)
docker compose up --build            # dev (COMPOSE_FILE in .env merges dev overrides)

# Production stack (on the VM)
docker compose -f docker-compose.yml up -d --build

# Backend (from backend/)
go run ./cmd/api                     # run natively (needs backend/.env + Postgres)
go test -race -cover ./...           # tests
go vet ./... && gofmt -l .            # lint
go run ./cmd/migrate up|down 1|version|new <name>

# Frontend (from frontend/)
npm run dev | build | lint | typecheck
```

## How to add a feature (recipe)

Backend (follow the auth slice as the reference implementation):
1. Domain: entity + repository interface + service interface + errors in
   `internal/domain/`.
2. Migration: `go run ./cmd/migrate new <name>`, write `.up.sql`/`.down.sql`.
3. Repository: GORM implementation in `internal/repository/`, translate GORM
   errors to domain errors.
4. Service: business logic in `internal/service/` + unit tests with fake repos.
5. Handler: DTOs + validation in `internal/handler/`, map domain errors.
6. Routes: register in `internal/routes/routes.go`.
7. Docs: add/extend `docs/api/<feature>.md`, update `docs/DATABASE.md`.

Frontend:
1. Types in `types/`, Zod schemas in `lib/validations/`.
2. Service functions in `services/<feature>.service.ts`.
3. React Query hooks in `hooks/use-<feature>.ts` (register keys in
   `lib/query-keys.ts`).
4. UI from `components/ui/*`; forms with React Hook Form + zodResolver.
5. Pages under `src/app/`.

## Current status (KEEP THIS UP TO DATE)

**Phase 1 — Foundation: COMPLETE**

**Phase 2–3 — Families / persons / tree building blocks: COMPLETE**

**Phase R — Aceh royal public pivot: COMPLETE**

**Phase 4 — Photos & search: COMPLETE**

**Phase 5 — JSON export/import: COMPLETE**

**Phase 6 — Admin users & password reset: COMPLETE**

**Phase 7 — Hardening: COMPLETE**

Implemented:
- Monorepo, docs, Docker, CI/CD, auth slice, person/relationship APIs,
  recursive CTE ancestors/descendants.
- Product pivot (ADR-0007): singleton royal family seed, public read APIs,
  edit-access request approve flow, Aceh landing + R3F tree, studio,
  admin approvals.
- Phase 4: photo upload + trigram name search.
- Phase 5: JSON export/import into the royal family (ADR-0008); GEDCOM deferred.
- Phase 6: admin user list/patch/soft-delete, password-reset tokens without
  SMTP (ADR-0009), `/admin/users` UI, forgot/reset password pages.
- Phase 7: nginx CSP + security headers (ADR-0010), `scripts/backup.sh`,
  `scripts/health-monitor.sh`, Playwright smoke E2E in CI against Compose.

Released:
- `v0.1.0` — foundation images on Docker Hub; SSH deploy still optional.

Not implemented yet (see `CONTEXT.md` / deferred):
- SMTP email delivery (verification + reset emails).
- GEDCOM export/import.
- Nonce-based stricter CSP (follow-up to ADR-0010).

API docs: `docs/api/public.md`, `docs/api/access-requests.md`,
`docs/api/persons.md`, `docs/api/export-import.md`, `docs/api/admin-users.md`.
Ops: `docs/DEPLOYMENT.md` (backups, monitoring, security headers).
