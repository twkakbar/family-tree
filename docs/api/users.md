# Users API

Implemented by `backend/internal/handler/user_handler.go`. Requires
`Authorization: Bearer <access token>`.

### `User` (response shape)

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

## GET /api/v1/users/me

Returns the authenticated user's profile.

Responses: `200` User · `401` when unauthenticated.

## PATCH /api/v1/users/me

Updates the caller's profile. Email changes stay out of scope until SMTP
verification exists. Password changes use forgot/reset (ADR-0009).

Request:

```json
{ "name": "Jane D. Doe" }
```

Validation: `name` required, 2–100 chars.

Responses: `200` updated User · `422 VALIDATION_ERROR` · `401`.

---

Admin user management (`GET /users`, role changes, deactivation) is a
Phase 6 deliverable under `/api/v1/admin/*` — see [`planned.md`](planned.md).
