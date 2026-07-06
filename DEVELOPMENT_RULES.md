# DEVELOPMENT_RULES.md — Coding Standards & Project Rules

These rules are binding for every contributor, human or AI. If a change
cannot be made without breaking a rule, stop and update this document first
(with justification) in the same PR.

## 1. General principles

- SOLID, Clean Architecture, Clean Code. Small units, clear names, single
  responsibility.
- No duplicate code: extract into `backend/pkg`, `internal/<layer>` helpers,
  or `frontend/src/{components,utils,lib}` once used twice.
- YAGNI: build what the current phase needs; pre-design (docs/contracts) is
  welcome, speculative code is not.
- Every dependency addition must be justified (purpose + why existing deps
  are insufficient) in the PR/commit description.
- Secrets and configuration only via environment variables. `.env` files are
  never committed; `.env.example` files are always kept complete and current.
- Documentation updates are part of the change, not an afterthought
  (triggers table in `AGENTS.md`).

## 2. Git workflow

- Branches: `feature/<slug>`, `fix/<slug>`, `chore/<slug>`.
- Conventional Commits: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`,
  `chore:`, `ci:` (+ optional scope, e.g. `feat(backend): …`).
- A PR/commit must leave the repo green: backend `gofmt`/`go vet`/tests pass,
  frontend lint/typecheck/build pass.
- Never force-push shared branches. Never commit generated artifacts
  (`bin/`, `.next/`, `node_modules/`, coverage files).

## 3. Backend rules (Go)

**Layering (violations are rejected):**
- HTTP handlers: bind/validate DTOs, call one service method, map errors,
  shape the response. NO business logic, NO DB access.
- Services: ALL business logic and transaction boundaries. NO gin/HTTP
  types, NO SQL/GORM.
- Repositories: persistence only. NO business rules. Translate GORM/driver
  errors into `domain` errors — callers must never see `gorm.*` errors.
- `domain` imports nothing from other app layers. `pkg/` never imports
  `internal/`.

**Code style:**
- `gofmt`-clean, `go vet`-clean, race-free tests (`go test -race ./...`).
- Exported identifiers have doc comments. Package names are short, lower
  case, no underscores.
- Errors: wrap with `fmt.Errorf("context: %w", err)`; compare with
  `errors.Is/As`; sentinel domain errors live in `internal/domain/errors.go`.
- `context.Context` is the first parameter of every service/repository
  method. Honor cancellation.
- Constructor injection only (`NewXxx(deps…)`); no globals, no `init()` side
  effects, no service locators.
- DTOs (request/response structs) live in `internal/handler`; domain
  entities are never serialized to clients directly.
- Validation at the edge: Gin binding tags + custom rules in
  `internal/validation`. Services still enforce invariants they own.
- Logging: `slog` only. No `fmt.Println`/`log.Printf`. Never log secrets,
  passwords, or tokens.

**Testing:**
- Services: unit tests with in-memory fakes of repository interfaces
  (pattern: `internal/service/auth_service_test.go`).
- `pkg/`: pure unit tests.
- Handlers/repos: integration tests are a Phase 7 goal; until then keep
  handlers thin enough to not need them.

## 4. Frontend rules (Next.js / TypeScript)

- TypeScript strict; the `any` type is forbidden (use `unknown` + narrowing).
- **No `fetch`/axios in components or pages.** Components → hooks →
  services → `lib/api-client.ts`. No exceptions.
- All server state through React Query; query keys only from
  `lib/query-keys.ts`.
- Forms: React Hook Form + `zodResolver` with schemas in `lib/validations/`.
  Map API 422 `details` onto form fields via `setError`.
- Server Components by default; add `"use client"` only for interactivity
  (forms, hooks, event handlers).
- Reuse `components/ui/*`; new primitives go there, feature-specific
  components in `components/<feature>/`.
- Styling: Tailwind utility classes + `cn()`; no inline `style`, no CSS
  modules, no styled-components. Design tokens in `globals.css` `@theme`.
- Never store tokens in `localStorage`/`sessionStorage`. Access token stays
  in `lib/auth-token.ts` memory store; refresh token stays in the httpOnly
  cookie (browser-managed).
- Env vars: browser-visible values must be prefixed `NEXT_PUBLIC_` and read
  only in `lib/env.ts`.

## 5. Database rules

- Schema changes ONLY via migration files in `backend/migrations/`
  (`go run ./cmd/migrate new <name>`), always with a working `.down.sql`.
- Never edit a migration that may have been applied anywhere. Never run
  GORM AutoMigrate. Never edit production schema by hand.
- Conventions: UUID PKs, `timestamptz`, `created_at`/`updated_at` on every
  table, `deleted_at` soft delete on user-facing entities (users, families,
  persons), snake_case identifiers, plural table names.
- Every FK declares `ON DELETE` behavior explicitly. Add indexes for FK
  columns and known query patterns; justify unusual indexes in
  `docs/DATABASE.md`.
- Queries must respect soft delete (GORM does this automatically; raw SQL
  must filter `deleted_at IS NULL` explicitly).

## 6. API design rules

- Versioned base path `/api/v1`. Nouns for resources, plural
  (`/families/:id/persons`). No verbs in paths except sub-actions that are
  not CRUD (e.g. `/export`).
- Methods: POST create, GET read, PATCH partial update, DELETE remove.
  Status codes: 200/201/204 success; 400 malformed; 401 unauthenticated;
  403 forbidden; 404 not found; 409 conflict; 422 validation; 429 rate
  limited; 500 internal.
- Every response uses the envelope (`pkg/response`). Error `code` values are
  stable, SCREAMING_SNAKE_CASE, and registered in `docs/api/README.md`.
- List endpoints are paginated by default.
- Every new/changed endpoint is documented in `docs/api/` in the same PR.

## 7. Security rules

- Never hardcode or log secrets/tokens/password hashes.
- bcrypt for passwords (cost from config, default 12). Passwords: min 8
  chars, require letters + digits (validator `strong_password`).
- JWT: HS256 only, `JWT_SECRET` ≥ 32 chars, short-lived access tokens.
  Refresh tokens: opaque, hashed at rest, rotating, reuse detection.
- Cookies: `HttpOnly`, `SameSite=Lax`, `Secure` in production.
- Validate ALL client input at the API boundary; enforce ownership checks in
  services (a user may only touch families they collaborate on).
- Uploads (Phase 4): validate MIME + size server-side, randomize stored
  filenames, never serve from the API process path directly (nginx serves
  the volume).
- Keep rate limiting on auth endpoints (nginx `limit_req`).

## 8. Definition of Done (checklist)

- [ ] Follows layer rules (§3/§4) and API/DB conventions (§5/§6).
- [ ] Backend: `gofmt -l .` empty, `go vet ./...` clean,
      `go test -race ./...` green.
- [ ] Frontend: `npm run lint`, `npm run typecheck`, `npm run build` green.
- [ ] New behavior covered by service-level unit tests (backend) where
      applicable.
- [ ] `.env.example` files updated if config changed.
- [ ] Documentation updated per the triggers table in `AGENTS.md`
      (including `AGENTS.md` Status section when a feature/phase lands).
- [ ] No secrets, no debug leftovers, no commented-out code.
