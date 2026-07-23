# ADR-0008: JSON tree export/import into the singleton royal family

- **Status:** Accepted
- **Date:** 2026-07-22

## Context

`docs/api/planned.md` drafted Phase 5 as:

- `GET /families/:id/export?format=json`
- `POST /families/import` → **creates a new family**

Phase R (ADR-0007) made the product a **single** public Aceh tree. Creating a
new family on import would produce data invisible to `/public/*` and revive
multi-family product surface we retired.

GEDCOM remains a stretch goal; evaluating it before shipping JSON would block
backup/restore for editors.

## Decision

1. **JSON only for Phase 5.** GEDCOM is deferred until there is a clear need
   and a dedicated ADR.
2. **Export:** `GET /families/:id/export?format=json` (viewer+). Payload
   `{format_version, exported_at, family, persons, relationships}`.
3. **Import:** `POST /families/:id/import` (editor+) **appends** into the
   target family with new UUIDs (remap). Does **not** create a new family.
   Photos are not embedded (paths ignored on import).
4. Document the deviation from the original planned “create NEW family”
   contract in `docs/api/export-import.md`.

## Consequences

- Editors can backup/restore or bulk-load a branch into the royal tree.
- Import is additive; replacing the whole tree requires manual cleanup first
  (or a future owner-only replace mode).
- GEDCOM can be added later without changing the JSON `format_version` story.
