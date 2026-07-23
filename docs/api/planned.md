# Planned API Contracts (Phases 2–6)

Pre-designed so future sessions implement against agreed shapes. These are
**drafts**: adjust freely while implementing, then move the endpoint into
its own `docs/api/<feature>.md` and update the index in `README.md`.
All endpoints require authentication; family-scoped endpoints additionally
enforce the caller's per-family role (owner/editor/viewer).

> **Families (Phase 2)** — implemented; see [`families.md`](families.md).
> **Persons, relationships, tree (Phase 3)** — implemented; see [`persons.md`](persons.md).
> **Photo upload + name search (Phase 4)** — implemented; see [`persons.md`](persons.md)
> (`PUT /persons/:id/photo`, `GET .../persons?search=`).
> **Export / import JSON (Phase 5)** — implemented; see [`export-import.md`](export-import.md)
> and ADR-0008 (import appends into the royal family; GEDCOM deferred).
> **Admin users + password reset (Phase 6)** — implemented; see [`admin-users.md`](admin-users.md),
> [`auth.md`](auth.md) (forgot/reset), and ADR-0009.

## Deferred / stretch

| Topic | Notes |
|-------|-------|
| GEDCOM export/import | Deferred (ADR-0008) |
| SMTP email delivery | Deferred; reset tokens are out-of-band (ADR-0009) |
| Email verification send | Flag exists; no verification email until SMTP |
| Stricter nonce-based CSP | Follow-up to ADR-0010 when Next nonce pipeline is practical |
