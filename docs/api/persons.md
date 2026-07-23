# Persons, Relationships & Tree API

Person records, relationship edges, and tree traversal (Phase 3), plus
photo upload and trigram name search (Phase 4).

Family-scoped endpoints enforce the caller's per-family role:
`viewer` (read), `editor` (mutate persons/relationships).

## Person object

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "family_id": "550e8400-e29b-41d4-a716-446655440001",
  "first_name": "John",
  "last_name": "Doe",
  "gender": "male",
  "birth_date": "1950-02-01",
  "death_date": null,
  "bio": "Family patriarch",
  "photo_url": "/uploads/persons/<uuid>.jpg",
  "created_at": "2026-07-06T12:00:00Z",
  "updated_at": "2026-07-06T12:00:00Z"
}
```

Dates are `YYYY-MM-DD` or `null`. `photo_url` is empty until a photo is
uploaded; otherwise it is a same-origin path under `/uploads/`.

## Relationship object

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "family_id": "550e8400-e29b-41d4-a716-446655440001",
  "from_person_id": "550e8400-e29b-41d4-a716-446655440000",
  "to_person_id": "550e8400-e29b-41d4-a716-446655440003",
  "type": "parent",
  "created_at": "2026-07-06T12:05:00Z"
}
```

- `parent`: directed — `from` is the parent, `to` is the child.
- `spouse`: undirected — stored once in canonical UUID order.

---

## `GET /families/:id/persons`

List/search persons in a family.

**Auth:** `viewer+`

**Query:** `search` (optional; trigram / ILIKE on full name), `page`, `per_page`

When `search` is set, results are ordered by name similarity (pg_trgm).

**Response `200`** — paginated array of Person objects.

---

## `PUT /persons/:id/photo`

Upload or replace a person's profile photo.

**Auth:** `editor+` on the person's family

**Body:** `multipart/form-data` with field `photo` (file).

| Rule | Value |
|------|-------|
| MIME | `image/jpeg`, `image/png`, or `image/webp` (sniffed server-side) |
| Size | ≤ `MAX_UPLOAD_SIZE_MB` (default 10) |

Stored as `persons/<person_id>.<ext>` under `UPLOAD_DIR`. Previous file is
removed on replace. Served publicly at `/uploads/...` (nginx / API static).

**Response `200`** — updated Person object (`photo_url` set).

**Errors:**
- `400 BAD_REQUEST` — missing file
- `422 VALIDATION_ERROR` — unsupported type or oversize

---

## `POST /families/:id/persons`

Create a person.

**Auth:** `editor+`

**Body**

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "gender": "male",
  "birth_date": "1950-02-01",
  "death_date": null,
  "bio": ""
}
```

| Field | Rules |
|-------|-------|
| `first_name` | required, 1–100 chars |
| `last_name` | optional, max 100 |
| `gender` | required, `male` \| `female` \| `other` \| `unknown` |
| `birth_date`, `death_date` | optional `YYYY-MM-DD`; death ≥ birth |
| `bio` | optional, max 5000 |

**Response `201`** — Person object.

---

## `GET /persons/:id`

Person details including incident relationships.

**Auth:** `viewer+` on the person's family

**Response `200`**

```json
{
  "success": true,
  "data": {
    "id": "…",
    "first_name": "John",
    "relationships": [ /* Relationship objects */ ]
  }
}
```

---

## `PATCH /persons/:id`

Update person fields (same body shape as create).

**Auth:** `editor+`

**Response `200`** — updated Person object.

---

## `DELETE /persons/:id`

Soft-delete the person and hard-delete all of its relationship edges in one
transaction.

**Auth:** `editor+`

**Response `200`**

```json
{ "success": true, "data": { "deleted": true } }
```

---

## `POST /families/:id/relationships`

Create a relationship edge.

**Auth:** `editor+`

**Body**

```json
{
  "from_person_id": "…",
  "to_person_id": "…",
  "type": "parent"
}
```

**Service invariants:**
- Both persons belong to the family.
- `parent`: child has ≤ 2 parent edges; no cycles (parent cannot be a
  descendant of the child).
- `spouse`: pair normalized to canonical order automatically.

**Response `201`** — Relationship object.

**Errors:**
- `422 VALIDATION_ERROR` — invariant violation
- `409 CONFLICT` — duplicate edge

---

## `DELETE /relationships/:id`

Remove a relationship edge (hard delete).

**Auth:** `editor+` on the relationship's family

**Response `200`**

```json
{ "success": true, "data": { "deleted": true } }
```

---

## `GET /families/:id/tree`

Whole-family graph for client rendering.

**Auth:** `viewer+`

**Response `200`**

```json
{
  "success": true,
  "data": {
    "persons": [ /* Person objects */ ],
    "relationships": [ /* Relationship objects */ ]
  }
}
```

---

## `GET /persons/:id/ancestors`

Recursive CTE upward through `parent` edges.

**Auth:** `viewer+`

**Query:** `depth` (default 5, max 20)

**Response `200`** — array of Person objects with extra `depth` field (1 =
direct parent generation).

---

## `GET /persons/:id/descendants`

Recursive CTE downward through `parent` edges.

**Auth:** `viewer+`

**Query:** `depth` (default 5, max 20)

**Response `200`** — array of Person objects with `depth` field.
