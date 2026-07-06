# ADR-0005: Same-origin routing via nginx (prod) / Next rewrites (dev)

**Status**: Accepted (2026-07)

## Context

Cookie-based refresh tokens (ADR-0004) work best first-party. Cross-origin
setups (app on one domain, API on another) force `SameSite=None`, broad CORS
with credentials, and third-party-cookie headaches.

## Decision

The browser always talks to ONE origin:

- **Production**: nginx routes `/api/*` → Go backend, `/uploads/*` → shared
  volume, everything else → Next.js.
- **Docker dev**: same nginx (parity), plus direct ports for convenience.
- **Native/hybrid dev**: Next.js `rewrites()` proxies `/api/*` and
  `/uploads/*` to `API_PROXY_TARGET`.

The frontend therefore uses the relative base URL `/api/v1` everywhere.

## Consequences

- Cookies are first-party; CORS is effectively unused in production (a
  minimal allowlist middleware remains for direct-to-API dev calls).
- Frontend needs no per-environment API URL configuration.
- Nginx is a hard dependency of the production topology (it also serves
  uploads and applies rate limits) — documented in docs/DEPLOYMENT.md.
