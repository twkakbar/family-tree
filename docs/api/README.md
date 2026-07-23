# API Conventions

Base URL: `/api/v1` (behind nginx in production, Next rewrites in dev).
Content type: `application/json`. All timestamps are RFC 3339 UTC.

## Response envelope

Every endpoint returns the same envelope:

```json
// success
{ "success": true, "data": { }, "meta": { } }

// failure
{ "success": false, "error": { "code": "CONFLICT", "message": "email is already registered", "details": null } }
```

`meta` appears only on paginated lists:

```json
{ "page": 1, "per_page": 20, "total": 45, "total_pages": 3 }
```

Pagination query parameters: `?page=1&per_page=20` (max `per_page` = 100;
invalid values fall back to defaults).

## Authentication

- Access token: JWT in `Authorization: Bearer <token>`, ~15 min TTL.
- Refresh token: httpOnly cookie `ft_refresh`, rotated on every refresh.
- Flow details: [`auth.md`](auth.md) and ADR-0004.

## Error code registry

Stable, machine-readable codes (add new ones here in the same PR):

| Code | HTTP | Meaning |
|------|------|---------|
| `BAD_REQUEST` | 400 | Malformed body/query (unparseable JSON, wrong types) |
| `VALIDATION_ERROR` | 422 | Failed validation; `details` = `{field: message}` |
| `UNAUTHORIZED` | 401 | Missing/invalid authentication |
| `INVALID_CREDENTIALS` | 401 | Login failed (email or password wrong) |
| `TOKEN_EXPIRED` | 401 | Access/refresh token expired |
| `TOKEN_INVALID` | 401 | Token malformed, unknown, revoked or reused |
| `FORBIDDEN` | 403 | Authenticated but not allowed |
| `NOT_FOUND` | 404 | Resource does not exist (or is soft-deleted) |
| `METHOD_NOT_ALLOWED` | 405 | Wrong HTTP method on a known route |
| `CONFLICT` | 409 | Uniqueness/state conflict (e.g. email taken) |
| `INTERNAL_ERROR` | 500 | Unexpected failure; correlate via `X-Request-ID` |

Rate limiting is enforced by nginx (HTTP 429 with an nginx body, not the
JSON envelope).

## Health probes (unversioned)

| Endpoint | Purpose |
|----------|---------|
| `GET /healthz` | Liveness: process is running (returns version) |
| `GET /readyz` | Readiness: database reachable |

## Endpoint index

| Area | Endpoints | Status |
|------|-----------|--------|
| Health | `GET /healthz`, `GET /readyz` (version in both) | ✅ implemented |
| Auth | `POST /auth/register`, `/login`, `/refresh`, `/logout`, `/forgot-password`, `/reset-password` | ✅ implemented |
| Users | `GET /users/me`, `PATCH /users/me` | ✅ implemented |
| Families | Singleton royal family (seeded); legacy CRUD retired from product | ✅ pivot |
| Persons, relationships, tree | see [`persons.md`](persons.md) — mutations for editors | ✅ implemented |
| Public tree | see [`public.md`](public.md) | ✅ implemented |
| Access requests | see [`access-requests.md`](access-requests.md) | ✅ implemented |
| Export / import (JSON) | see [`export-import.md`](export-import.md) | ✅ implemented |
| Admin users | see [`admin-users.md`](admin-users.md) | ✅ implemented |
