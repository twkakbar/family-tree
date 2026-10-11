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
    "email_verified": true,
    "created_at": "2026-07-06T19:00:00Z"
  }
}
```

---

## POST /api/v1/auth/register

Creates an account. When SMTP is configured, the account stays unverified
and the response is `{ "verification_required": true }` with no session.
Without SMTP, the account is verified immediately and a session is created.

Request:

```json
{ "name": "Jane Doe", "email": "jane@example.com", "password": "hunter2024" }
```

Validation: `name` 2–100 chars; `email` valid + unique among active
accounts; `password` 8–128 chars with at least one letter and one digit.

Responses: `201` AuthResult + refresh cookie when SMTP is off ·
`201` `{ "verification_required": true }` when SMTP is on ·
`409 CONFLICT` (verified email taken) ·
`429 RATE_LIMITED` (more than 3 verification emails in an hour) ·
`422 VALIDATION_ERROR`.

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
`403 EMAIL_NOT_VERIFIED` when the password is correct but the address is
not confirmed yet · `422 VALIDATION_ERROR`.

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

## POST /api/v1/auth/verify-email

```json
{ "token": "…" }
```

Responses: `200` `{ "verified": true }` · `401 TOKEN_INVALID | TOKEN_EXPIRED`.

## POST /api/v1/auth/resend-verification

```json
{ "email": "jane@example.com" }
```

Sends another link when the account exists and is still unverified.
Unknown or already-verified addresses return the same `200`. The 3-per-hour
cap still applies to real sends.

## POST /api/v1/auth/forgot-password

Creates a password-reset token when the email exists. Always returns `200`
with the same shape (no account enumeration). See ADR-0009.

Request:

```json
{ "email": "jane@example.com" }
```

Response `data` (production):

```json
{
  "requested": true,
  "message": "If that email is registered, a reset token was created."
}
```

In **non-production**, when a token was created, `reset_token` may also be
present so local testing works without an admin.

## POST /api/v1/auth/reset-password

Consumes a one-time token and sets a new password. Revokes all sessions.

Request:

```json
{ "token": "…", "new_password": "hunter2024" }
```

Validation: same strong password rules as register. Responses: `200`
`{ "reset": true }` · `401 TOKEN_INVALID | TOKEN_EXPIRED` · `422`.

---

## Client integration (frontend reference)

`frontend/src/lib/api-client.ts`: attach access token from memory → on 401
run one single-flight `/auth/refresh` → retry the original request once →
on refresh failure clear state and redirect to `/login`.

User object now includes `email_verified` (boolean).
