# Auth API

Implemented by `backend/internal/handler/auth_handler.go`. Design rationale:
ADR-0004. All bodies use the [response envelope](README.md); examples below
show only `data`.

Session model recap: the response body carries a short-lived JWT **access
token**; a rotating opaque **refresh token** is set as an httpOnly cookie
(`ft_refresh`, `SameSite=Lax`, `Secure` in production, path `/`). The
refresh token never appears in JSON.

### `AuthResult` (shared response shape)

```json
{
  "access_token": "eyJhbGciOi…",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "id": "6f1e…",
    "name": "Jane Doe",
    "email": "jane@example.com",
    "role": "user",
    "created_at": "2026-07-06T19:00:00Z"
  }
}
```

---

## POST /api/v1/auth/register

Creates an account and immediately establishes a session (auto-login).

Request:

```json
{ "name": "Jane Doe", "email": "jane@example.com", "password": "hunter2024" }
```

Validation: `name` 2–100 chars; `email` valid + unique among active
accounts; `password` 8–128 chars with at least one letter and one digit.

Responses: `201` AuthResult + refresh cookie · `422 VALIDATION_ERROR` ·
`409 CONFLICT` (email taken).

```bash
curl -i -X POST http://localhost/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Jane Doe","email":"jane@example.com","password":"hunter2024"}'
```

## POST /api/v1/auth/login

Request:

```json
{ "email": "jane@example.com", "password": "hunter2024" }
```

Responses: `200` AuthResult + refresh cookie · `401 INVALID_CREDENTIALS`
(identical for unknown email and wrong password — no account enumeration) ·
`422 VALIDATION_ERROR`.

## POST /api/v1/auth/refresh

No body; the browser sends the `ft_refresh` cookie automatically
(`withCredentials: true`).

Behavior: **rotation** — the presented token is revoked, a new pair is
issued. **Reuse detection** — presenting an already-revoked token revokes
all of the user's refresh tokens and fails.

Responses: `200` AuthResult + NEW refresh cookie ·
`401 TOKEN_INVALID | TOKEN_EXPIRED` (cookie cleared).

## POST /api/v1/auth/logout

Revokes the presented refresh token and clears the cookie. Idempotent —
succeeds even without a valid session. The frontend also drops the
in-memory access token (it expires within minutes regardless).

Response: `204 No Content`.

---

## Client integration (frontend reference)

`frontend/src/lib/api-client.ts`: attach access token from memory → on 401
run one single-flight `/auth/refresh` → retry the original request once →
on refresh failure clear state and redirect to `/login`.
