# ADR-0004: JWT access token + rotating opaque refresh cookie

**Status**: Accepted (2026-07)

## Context

SPA authentication must balance statelessness, XSS/CSRF exposure and
revocability. Storing long-lived tokens in localStorage is the main
anti-pattern to avoid (trivially exfiltrated by XSS).

## Decision

- **Access token**: HS256 JWT, ~15 min TTL, claims `sub/role/jti/iss`.
  Kept in frontend **memory only**; sent as `Authorization: Bearer`.
- **Refresh token**: 256-bit opaque random value; SHA-256 hash stored in
  `refresh_tokens`. Delivered as `HttpOnly; SameSite=Lax` cookie
  (`Secure` in production). **Rotated on every use**; reuse of a consumed
  token revokes all of the user's sessions (theft signal).
- Revocation latency for role changes/bans = access token TTL (≤ 15 min),
  acceptable per CONTEXT.md constraints.

## Consequences

- XSS cannot steal the refresh token (httpOnly) and the access token is
  short-lived and never persisted. SameSite=Lax + same-origin routing
  (ADR-0005) keeps CSRF exposure minimal for the cookie-consuming endpoints.
- Page reloads require a silent `/auth/refresh` round-trip on first API 401
  (implemented in `api-client.ts` as a single-flight).
- HS256 symmetric signing is sufficient for a single-service backend;
  moving to multiple token consumers would motivate RS256/JWKS (new ADR).
