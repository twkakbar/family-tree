# Planned API Contracts (Phases 2–6)

Pre-designed so future sessions implement against agreed shapes. These are
**drafts**: adjust freely while implementing, then move the endpoint into
its own `docs/api/<feature>.md` and update the index in `README.md`.
All endpoints require authentication; family-scoped endpoints additionally
enforce the caller's per-family role (owner/editor/viewer).

## Families (Phase 2)

| Method & path | Role | Purpose |
|---------------|------|---------|
| `POST /families` | any user | Create family; creator becomes `owner` collaborator |
| `GET /families` | any user | List families the caller collaborates on (paginated) |
| `GET /families/:id` | viewer+ | Family details incl. caller's role |
| `PATCH /families/:id` | owner | Update name/description |
| `DELETE /families/:id` | owner | Soft delete |

```json
// POST /families request / Family response shape
{ "name": "The Does", "description": "Paternal branch" }
{ "id": "…", "name": "The Does", "description": "…", "role": "owner",
  "persons_count": 0, "created_at": "…", "updated_at": "…" }
```

### Collaborators

| Method & path | Role | Purpose |
|---------------|------|---------|
| `GET /families/:id/collaborators` | viewer+ | List members with roles |
| `POST /families/:id/collaborators` | owner | Invite existing user `{ "email": "…", "role": "editor" }` |
| `PATCH /families/:id/collaborators/:userId` | owner | Change role |
| `DELETE /families/:id/collaborators/:userId` | owner | Revoke access (owner cannot remove the last owner) |

## Persons (Phase 3, photo in Phase 4)

| Method & path | Role | Purpose |
|---------------|------|---------|
| `GET /families/:id/persons?search=&page=` | viewer+ | List/search (ILIKE on names; pg_trgm index) |
| `POST /families/:id/persons` | editor+ | Create person |
| `GET /persons/:id` | viewer+ | Person details incl. relationships |
| `PATCH /persons/:id` | editor+ | Update fields |
| `DELETE /persons/:id` | editor+ | Soft delete person + hard delete its edges (one transaction) |
| `PUT /persons/:id/photo` | editor+ | multipart upload (jpeg/png/webp ≤ 10 MB) → stores under `/uploads`, returns `photo_url` |

```json
// Person response shape
{ "id": "…", "family_id": "…", "first_name": "John", "last_name": "Doe",
  "gender": "male", "birth_date": "1950-02-01", "death_date": null,
  "bio": "…", "photo_url": "/uploads/persons/<file>.webp",
  "created_at": "…", "updated_at": "…" }
```

## Relationships (Phase 3)

| Method & path | Role | Purpose |
|---------------|------|---------|
| `POST /families/:id/relationships` | editor+ | `{ "from_person_id": "…", "to_person_id": "…", "type": "parent" \| "spouse" }` |
| `DELETE /relationships/:id` | editor+ | Remove edge |

Service-layer invariants (beyond DB constraints): both persons belong to the
family; a child has ≤ 2 `parent` edges; no parent-cycle (reject if the child
is an ancestor of the parent); spouse pairs normalized to canonical order.

## Tree & traversal (Phase 3)

| Method & path | Role | Purpose |
|---------------|------|---------|
| `GET /families/:id/tree` | viewer+ | Whole tree: `{ "persons": [...], "relationships": [...] }` — client renders |
| `GET /persons/:id/ancestors?depth=5` | viewer+ | Recursive CTE upward (docs/DATABASE.md) |
| `GET /persons/:id/descendants?depth=5` | viewer+ | Recursive CTE downward |

## Export / import (Phase 5)

| Method & path | Role | Purpose |
|---------------|------|---------|
| `GET /families/:id/export?format=json` | viewer+ | Download `{format_version, family, persons, relationships}` |
| `POST /families/import` | any user | Upload an export → creates a NEW family, remaps all IDs |

GEDCOM support is a stretch goal — evaluate before Phase 5 and record an ADR.

## Admin (Phase 6) — under `/api/v1/admin`, global role `admin`

| Method & path | Purpose |
|---------------|---------|
| `GET /admin/users?search=&page=` | List accounts (paginated) |
| `PATCH /admin/users/:id` | Change global role / name |
| `DELETE /admin/users/:id` | Soft-delete account (revokes sessions) |
