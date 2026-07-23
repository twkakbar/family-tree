# Tree export / import (JSON)

Phase 5 backup and bulk-load for the royal family tree. See ADR-0008 for why
import targets an **existing** family (append) instead of creating a new one.

All endpoints require authentication.

## Document shape (`format_version: 1`)

```json
{
  "format_version": 1,
  "exported_at": "2026-07-22T00:00:00Z",
  "family": {
    "name": "Kesultanan Aceh",
    "description": "…"
  },
  "persons": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "first_name": "Cut Nyak Dhien",
      "last_name": "",
      "gender": "female",
      "birth_date": "1848-01-01",
      "death_date": "1908-11-06",
      "bio": ""
    }
  ],
  "relationships": [
    {
      "from_person_id": "…",
      "to_person_id": "…",
      "type": "parent"
    }
  ]
}
```

- Person `id` values are export-time UUIDs used only to remap relationships on
  import (new UUIDs are always assigned).
- Dates are `YYYY-MM-DD` or `null`.
- Photos are **not** included in the JSON.

---

## `GET /families/:id/export?format=json`

Download the current tree.

**Auth:** `viewer+` on the family

**Query:** `format` — only `json` is supported (default).

**Response `200`** — TreeTransfer document (also sets `Content-Disposition`
attachment filename).

**Errors:** `422` if `format` is not `json`.

---

## `POST /families/:id/import`

Append persons and relationships from a TreeTransfer document into the family.

**Auth:** `editor+`

**Body:** TreeTransfer JSON (`format_version` must be `1`).

**Behavior:**
- Creates new person rows (remapped IDs).
- Creates relationships among those new rows (same invariants as normal create:
  ≤2 parents, no cycles, spouse canonicalization).
- Does **not** create a new family; does **not** delete existing members.
- Max 2000 persons per import.

**Response `200`**

```json
{
  "success": true,
  "data": {
    "persons_created": 12,
    "relationships_created": 10
  }
}
```

**Errors:**
- `422 VALIDATION_ERROR` — bad/unsupported document
- `403` — not an editor
- Relationship invariant failures surface as `422` / `409` as usual

## GEDCOM

Not implemented. Deferred per ADR-0008.
