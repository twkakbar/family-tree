# Access requests & studio

Registered descendants request edit access; super-admins approve or reject.

## `POST /access-requests`

**Auth:** any authenticated user who is not already an editor/admin.

**Body**

```json
{ "message": "I descend from … and wish to connect my living relatives." }
```

`message` required, 10–2000 characters.

**Response `201`** — access request object.

**Errors:** `409` if already pending (`CONFLICT`) or already editor.

## `GET /access-requests/me`

**Auth:** authenticated

**Response `200`** — latest request for the caller, or `404` if none.

```json
{
  "id": "…",
  "status": "pending",
  "message": "…",
  "created_at": "…",
  "reviewed_at": null
}
```

## `GET /studio/can-edit`

**Auth:** authenticated

```json
{ "can_edit": true, "family_id": "…" }
```

`family_id` is the singleton royal family id used for mutation endpoints.

## Admin

All under `/api/v1/admin`, require global role `admin`.

### `GET /admin/access-requests`

Lists **pending** requests with applicant name/email.

### `POST /admin/access-requests/:id/approve`

Creates an `editor` collaborator grant on the royal family and marks the
request `approved`.

### `POST /admin/access-requests/:id/reject`

Marks the request `rejected` without granting access.
