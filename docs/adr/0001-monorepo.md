# ADR-0001: Monorepo with root-level AI context documents

**Status**: Accepted (2026-07)

## Context

The project is developed AI-assisted across many sessions that lose
conversation context. Backend, frontend and infrastructure evolve together
and are deployed together to a single VM.

## Decision

One repository containing `backend/`, `frontend/`, `docker/`, `docs/` and
compose files, with four canonical context documents at the root
(`AGENTS.md`, `CONTEXT.md`, `ARCHITECTURE.md`, `DEVELOPMENT_RULES.md`).
`AGENTS.md` carries a maintained "current status" section as the session
memory anchor.

## Consequences

- Atomic cross-stack changes (one PR touches API + client + docs).
- Single CI pipeline with per-area jobs; single deploy artifact set.
- Docs must be updated with the code they describe (enforced by
  DEVELOPMENT_RULES §8); stale docs are treated as bugs.
