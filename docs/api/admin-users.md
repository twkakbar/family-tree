# Admin users API

Implemented by `backend/internal/handler/admin_user_handler.go`.
Requires global role `admin`. Design notes: ADR-0009 (password reset without SMTP).

All bodies use the [response envelope](README.md); examples show only `data`.

### User (admin list/detail shape)

```json
{
  "id": "6f1e…",
  "name": "Jane Doe",
  "email": "jane@example.com",
  "role": "user",
  "email_verified": true,
  "created_at": "2026-07-06T19:00:00Z"
}
```

---

## GET /api/v1/admin/users

Paginated account list. Query: `?search=&page=&per_page=`.

`search` matches name or email (case-insensitive substring).

Response: `200` array of User + `meta` · `401` · `403`.

## PATCH /api/v1/admin/users/:id

Partial update. Body (all optional):

```json
{ "name": "New Name", "role": "admin", "email_verified": true }
```

Rules:

- `role` must be `admin` or `user`.
- An admin cannot demote themselves.
- `email_verified: true` sets `email_verified_at` to now; `false` clears it.
  There is no SMTP verification flow yet (ADR-0009).

Responses: `200` User · `403` · `404` · `422`.

## DELETE /api/v1/admin/users/:id

Soft-deletes the account and revokes refresh + password-reset tokens.
Cannot delete yourself.

Responses: `200` `{ "deleted": true }` · `403` · `404`.

## POST /api/v1/admin/users/:id/password-reset

Creates a one-time reset token (1 hour TTL) and returns the **raw** token so
the admin can share it out-of-band (chat, phone). No email is sent.

```json
{
  "reset_token": "…",
  "expires_in": 3600,
  "hint": "Share /reset-password?token=… with the user out-of-band (no email delivery yet)."
}
```

Responses: `200` · `403` · `404`.
