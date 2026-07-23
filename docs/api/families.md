# Families API

> **Product note (Phase R):** The public product exposes a **single seeded royal
> family**. Multi-family create/list and collaborator invite-by-email are no
> longer product surfaces — editors resolve the family id via
> `GET /studio/can-edit`. The endpoints below remain for studio mutations and
> internal use.

Family containers and collaborator management. Authenticated endpoints use
`Authorization: Bearer <access_token>`.

Per-family roles (most to least privileged): `owner`, `editor`, `viewer`.
The bootstrap admin is `owner` of the seeded royal family; approved descendants
receive `editor` via access-request approval.

## Family object

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "The Does",
  "description": "Paternal branch",
  "role": "owner",
  "persons_count": 0,
  "created_at": "2026-07-06T12:00:00Z",
  "updated_at": "2026-07-06T12:00:00Z"
}
```

`role` is the **caller's** role in that family. `persons_count` counts
non-deleted persons (Phase 3).

## Collaborator object

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "Jane Doe",
  "email": "jane@example.com",
  "role": "editor",
  "created_at": "2026-07-06T12:05:00Z"
}
```

Invites target **existing user accounts** by email. There is no email-delivery
flow yet; the user must already be registered.

---

## `POST /families`

Create a family. The caller becomes `owner`.

**Auth:** any authenticated user

**Body**

```json
{ "name": "The Does", "description": "Optional note" }
```

| Field | Rules |
|-------|-------|
| `name` | required, 2–100 chars |
| `description` | optional, max 1000 chars |

**Response `201`** — Family object (see above).

---

## `GET /families`

List families the caller collaborates on, newest first.

**Auth:** any authenticated user

**Query:** `page` (default 1), `per_page` (default 20, max 100)

**Response `200`**

```json
{
  "success": true,
  "data": [ /* Family objects */ ],
  "meta": { "page": 1, "per_page": 20, "total": 2, "total_pages": 1 }
}
```

---

## `GET /families/:id`

Family details including the caller's role.

**Auth:** `viewer+` on the family

**Response `200`** — Family object.

**Errors:** `404 NOT_FOUND` if the family does not exist or the caller is not
a collaborator.

---

## `PATCH /families/:id`

Update name and description.

**Auth:** `owner`

**Body**

```json
{ "name": "Renamed Family", "description": "Updated note" }
```

**Response `200`** — updated Family object.

**Errors:** `403 FORBIDDEN` if the caller is not an owner.

---

## `DELETE /families/:id`

Soft-delete the family and cascade-delete collaborators, persons, and
relationships (via FK constraints).

**Auth:** `owner`

**Response `200`**

```json
{ "success": true, "data": { "deleted": true } }
```

---

## `GET /families/:id/collaborators`

List all collaborators with public user profile fields.

**Auth:** `viewer+`

**Response `200`** — array of Collaborator objects.

---

## `POST /families/:id/collaborators`

Invite an existing user by email.

**Auth:** `owner`

**Body**

```json
{ "email": "collaborator@example.com", "role": "editor" }
```

| Field | Rules |
|-------|-------|
| `email` | required, valid email of an existing account |
| `role` | required, `editor` or `viewer` (`owner` cannot be assigned via invite) |

**Response `201`** — Collaborator object.

**Errors:**
- `404 NOT_FOUND` — email not registered
- `409 CONFLICT` — user is already a collaborator

---

## `PATCH /families/:id/collaborators/:userId`

Change a collaborator's role.

**Auth:** `owner`

**Body**

```json
{ "role": "viewer" }
```

Allowed values: `owner`, `editor`, `viewer`.

**Response `200`** — updated Collaborator object.

**Errors:**
- `409 CONFLICT` — would remove or demote the last owner

---

## `DELETE /families/:id/collaborators/:userId`

Revoke a user's access to the family.

**Auth:** `owner`

**Response `200`**

```json
{ "success": true, "data": { "removed": true } }
```

**Errors:**
- `409 CONFLICT` — cannot remove the last owner
