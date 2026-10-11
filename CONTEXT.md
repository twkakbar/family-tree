# CONTEXT.md — Business & Product Context

This document explains **what** the application is and **why** it exists.
For **how** it is built, read `ARCHITECTURE.md`.

## 1. Business overview

**Kesultanan Aceh Family Tree** is a public, self-hosted web application that
presents the royal lineage of the Aceh sultanate. Anyone can explore the
canonical tree in an immersive 3D viewer without logging in. Living
descendants who wish to connect their branch to the sultan line register,
request edit access, and — after a super-admin approves them — may add
persons and relationships in a private studio.

This is **not** a multi-tenant genealogy SaaS. There is exactly **one**
canonical family tree for the whole site.

## 2. Goals

- Public, beautiful exploration of the Aceh royal family tree (3D + motion).
- Controlled contribution: only approved editors (and super-admins) mutate data.
- Data ownership: self-hosted, no third-party services for core features.
- Production quality: secure auth for editors, migrations, observability.
- AI-maintainable: documentation-first repo so any session can continue work.

### Non-goals (explicitly out of scope)

- Multi-family creation / collaborator invite-by-email product UX.
- Social features (feeds, comments, likes), DNA integrations, billing,
  mobile native apps, multi-tenant SaaS hosting.
- Email delivery for access requests or password reset (those stay in-app /
  out-of-band). Registration verification email is sent when SMTP is configured.

## 3. Domain glossary (use these terms precisely)

| Term | Meaning |
|------|---------|
| **User** | An account that can log in. Global role: `admin` (super-admin) or `user`. |
| **Super-admin** | Global role `admin`. Seeded from env. Approves edit-access requests and may edit the tree. |
| **Royal family** | The single seeded family container (default name: *Kesultanan Aceh*). |
| **Editor grant** | A `family_collaborators` row with role `editor` (or `owner` for the bootstrap admin) on the royal family. |
| **Access request** | A registered user asking for editor rights; status `pending` / `approved` / `rejected`. |
| **Person** | A node in the royal tree (a family member). Pure data — not a login. |
| **Relationship** | Typed edge: `parent` (from parent → to child) or `spouse` (canonical order). |
| **Tree** | The rendered graph of persons + relationships of the royal family. |

The distinction **User ≠ Person** remains fundamental.

## 4. Features

| # | Feature | Status | Phase |
|---|---------|--------|-------|
| 1 | Authentication (register, login, refresh, logout) | ✅ Implemented | 1 |
| 2 | User profile (`/users/me`) | ✅ Implemented | 1 |
| 3 | Global roles (admin/user) + RBAC middleware | ✅ Implemented | 1 |
| 4 | Singleton royal family (seeded) + public read APIs | ✅ Implemented | R |
| 5 | Edit-access requests + super-admin approve/reject | ✅ Implemented | R |
| 6 | Studio: add/edit persons & relationships (editors) | ✅ Implemented | R |
| 7 | Public 3D tree explorer | ✅ Implemented | R |
| 8 | Ancestors / descendants (API + UI) | ✅ Implemented | 3 |
| 9 | Profile picture upload for persons | ✅ Implemented | 4 |
| 10 | Search family members | ✅ Implemented | 4 |
| 11 | Export / import (JSON) | ✅ Implemented | 5 |
| 12 | Admin user management + password reset (no SMTP) | ✅ Implemented | 6 |
| 13 | Hardening (CSP, E2E smoke, backup/health scripts) | ✅ Implemented | 7 |

Legacy multi-family CRUD / collaborator-invite UX is **retired** as a product
surface (schema retained for the singleton).

## 5. Constraints

**Technical**
- Fixed stack: Go/Gin/GORM/PostgreSQL backend, Next.js/TypeScript frontend
  (see `AGENTS.md`). Deviations require an ADR.
- 3D viewer uses React Three Fiber (ADR-0007).
- Deployment: single Linux VM, Docker Compose, Nginx.
- All configuration via environment variables; secrets never committed.

**Non-functional**
- Expected scale: one tree with up to ~2,000 persons; hundreds of users.
- Security: public reads; mutations require JWT + editor/admin; refresh
  cookies httpOnly; rate limiting at the proxy.

## 6. User flows

**Public exploration**
1. Visitor opens the landing page → taps Explore → immersive 3D tree.
2. Click a person → detail panel (bio, dates, ancestors/descendants).

**Descendant contribution**
1. Register / login → submit an access request (message describing the branch).
2. Super-admin reviews in `/admin/requests` → approve or reject.
3. On approve, user gains editor grant → uses `/studio` to add living
   descendants and link them toward the sultan line.

**Super-admin**
1. Bootstrap admin from `ADMIN_EMAIL` / `ADMIN_PASSWORD` on first startup.
2. Approves requests; may edit the tree like an editor (owner of the seeded family).
3. Manages accounts at `/admin/users` (role, email-verified flag, soft-delete,
   issue password-reset tokens to share out-of-band — no SMTP yet).

## 7. Roadmap

| Phase | Deliverable | Status |
|-------|-------------|--------|
| 1 | Foundation: monorepo, docs, Docker, CI/CD, DB schema, auth | ✅ Done |
| 2–3 | Families/persons/tree (internal building blocks) | ✅ Done |
| **R** | **Product pivot: public Aceh tree, access requests, 3D explore, studio** | ✅ Done |
| 4 | Photo uploads, member search (pg_trgm) | ✅ Done |
| 5 | Export/import JSON (GEDCOM deferred, ADR-0008) | ✅ Done |
| 6 | Broader admin UI, email verification/password reset (ADR-0009) | ✅ Done |
| 7 | Hardening: E2E smoke, CSP, backups script, health monitoring | ✅ Done |

When a phase is completed: update this table, §4, and `AGENTS.md` status.
