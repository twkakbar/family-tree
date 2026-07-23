# ADR-0010: Content-Security-Policy and hardened security headers

- **Status:** Accepted
- **Date:** 2026-07-23

## Context

Phase 7 requires CSP and broader hardening. The public UI uses Next.js App
Router and React Three Fiber (WebGL). Strict nonce-only CSP for Next.js is
possible but fragile without a dedicated nonce pipeline in the App Router
layout; shipping zero CSP leaves XSS blast radius unbounded.

Nginx already set `X-Content-Type-Options`, `X-Frame-Options`, and
`Referrer-Policy`. Uploads location overrides `add_header` and previously
dropped parent security headers (nginx inheritance quirk).

## Decision

1. Centralize headers in `docker/nginx/conf.d/security-headers.inc` and
   `include` them from the server block and from `/uploads/`.
2. Enforce a **pragmatic CSP** that allows Next + Three.js to run:
   - `script-src 'self' 'unsafe-inline' 'unsafe-eval' 'wasm-unsafe-eval'`
   - `style-src 'self' 'unsafe-inline'`
   - `worker-src 'self' blob:`
   - `img-src 'self' data: blob:`
   - `connect-src 'self'` (same-origin API; ADR-0005)
   - `frame-ancestors 'none'`, `object-src 'none'`, `base-uri 'self'`
3. Add `Permissions-Policy`, `Cross-Origin-Opener-Policy`, and
   `Cross-Origin-Resource-Policy` at the edge.
4. Keep Gin `SecureHeaders` as defense-in-depth for direct API access (no
   full CSP on JSON API responses — browsers ignore CSP for XHR/fetch JSON).
5. Document HSTS (and optional `upgrade-insecure-requests`) only on the TLS
   server block — not on plain HTTP, so local Docker/CI keep working.

## Consequences

- XSS impact is reduced versus no CSP; script allowances remain broad until
  a nonce-based Next CSP lands.
- Changing frontend CDNs / analytics requires CSP updates (currently none).
- `/uploads/` responses again carry security headers alongside cache headers.
