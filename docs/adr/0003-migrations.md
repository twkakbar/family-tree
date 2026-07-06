# ADR-0003: golang-migrate with embedded SQL (AutoMigrate forbidden)

**Status**: Accepted (2026-07)

## Context

GORM AutoMigrate mutates schema implicitly from struct tags: no review, no
rollback, no partial-index/CHECK/extension support, and drift between
environments. Production requires deterministic, versioned schema changes.

## Decision

Plain SQL migrations in `backend/migrations/` managed by golang-migrate,
embedded into the binary (`embed.FS`) and applied on startup when
`RUN_MIGRATIONS=true`, plus a `cmd/migrate` CLI (`up/down/version/force/new`).
AutoMigrate is forbidden (DEVELOPMENT_RULES §5).

## Consequences

- Reviewable, reversible, environment-identical schema; full Postgres
  feature access (citext, partial indexes, CHECKs, extensions).
- Entities and migrations must be kept in sync manually — the migration is
  the source of truth.
- Startup auto-migration suits the single-instance VM; multi-instance
  rollouts must run `cmd/migrate` as a separate step instead.
