# CONTEXT.md — Business & Product Context

This document explains **what** the Family Tree application is and **why** it
exists. For **how** it is built, read `ARCHITECTURE.md`.

## 1. Business overview

Family Tree is a self-hosted web application for creating, managing and
exploring family trees. A user creates one or more *families*, adds *persons*
(family members) to them, connects persons with *relationships*
(parent/child, spouse), and views the result as an interactive tree. Families
can be shared with other users for collaboration, and trees can be exported
and imported for backup or migration.

The product is aimed at individuals and small groups (a family) who want to
preserve their family history on infrastructure they control — not a social
network or a genealogy research platform.

## 2. Goals

- Simple, modern UX for building and browsing family trees.
- Multi-user: accounts, shared families with per-family permissions.
- Data ownership: self-hosted, exportable data, no third-party services.
- Production quality: secure auth, backups, migrations, observability.
- AI-maintainable: documentation-first repo so any AI/dev session can
  continue the work without prior conversation context.

### Non-goals (explicitly out of scope)

- Social features (feeds, comments, likes), DNA/records integrations,
  billing/subscriptions, mobile native apps, multi-tenant SaaS hosting.

## 3. Domain glossary (use these terms precisely)

| Term | Meaning |
|------|---------|
| **User** | An account that can log in. Has a global role: `admin` or `user`. |
| **Family** | A container for one family tree. Created and owned by a user. |
| **Collaborator** | A user granted access to a family with a per-family role: `owner`, `editor` or `viewer`. The creator is automatically `owner`. |
| **Person** | A node in a family tree (a family member). Persons are data records, NOT user accounts. A person belongs to exactly one family. |
| **Relationship** | A typed edge between two persons in the same family: `parent` (directed: from parent → to child) or `spouse` (undirected, stored once in canonical order). |
| **Tree** | The rendered graph of persons + relationships of a family. |

The distinction **User ≠ Person** is fundamental: deleting a user never
deletes persons; persons have no login.

## 4. Features

| # | Feature | Status | Phase |
|---|---------|--------|-------|
| 1 | Authentication (register, login, refresh, logout) | ✅ Implemented | 1 |
| 2 | User profile (`/users/me` view + edit) | ✅ Implemented | 1 |
| 3 | Global roles (admin/user) + RBAC middleware | ✅ Implemented | 1 |
| 4 | Create/manage families | 🔲 Planned | 2 |
| 5 | Share family with collaborators (owner/editor/viewer) | 🔲 Planned | 2 |
| 6 | Add/edit family members (persons) | 🔲 Planned | 3 |
| 7 | Define & edit relationships (parent, spouse) | 🔲 Planned | 3 |
| 8 | Interactive family tree view | 🔲 Planned | 3 |
| 9 | View ancestors / descendants of a person | 🔲 Planned | 3 |
| 10 | Profile picture upload for persons | 🔲 Planned | 4 |
| 11 | Search family members | 🔲 Planned | 4 |
| 12 | Export family tree (JSON; GEDCOM later) | 🔲 Planned | 5 |
| 13 | Import family tree (JSON; GEDCOM later) | 🔲 Planned | 5 |
| 14 | Admin user management UI | 🔲 Planned | 6 |

The database schema for features 4–9 already exists (see `docs/DATABASE.md`);
API contracts for them are drafted in `docs/api/planned.md`.

## 5. Constraints

**Technical**
- Fixed stack: Go/Gin/GORM/PostgreSQL backend, Next.js/TypeScript frontend
  (see `AGENTS.md`). Deviations require an ADR.
- Deployment target: a single Linux VM running Docker Compose behind Nginx.
  No managed cloud services; no Kubernetes.
- Uploaded images are stored on a Docker volume on the VM (S3-compatible
  storage is a possible future ADR).
- All configuration via environment variables; secrets are never committed.

**Non-functional**
- Expected scale: hundreds of users, families with up to ~2,000 persons.
  Optimize for correctness and maintainability first; recursive CTEs and
  proper indexes cover this scale comfortably.
- Availability: single node; brief downtime during deploys is acceptable.
- Security: OWASP-conscious auth (short-lived JWT access tokens, rotating
  httpOnly refresh cookies, bcrypt, rate limiting at the proxy).

## 6. User flows

**Onboarding**
1. Visitor opens the landing page → registers with name/email/password.
2. Registration logs the user in (access token + refresh cookie) and
   redirects to the dashboard.

**Create a family & build a tree** *(Phases 2–3)*
1. User creates a family (name, description) → becomes its `owner`.
2. User adds persons (name, gender, birth/death dates, bio, photo).
3. User connects persons: select person → add relationship (parent/spouse).
4. Tree view renders the family; user can pan/zoom, focus a person, and view
   their ancestors or descendants.

**Collaboration** *(Phase 2)*
1. Owner invites an existing user by email to a family with a role.
2. `editor` can modify persons/relationships; `viewer` is read-only;
   `owner` can also manage collaborators and delete the family.

**Export / import** *(Phase 5)*
1. Owner exports a family as JSON (all persons + relationships + metadata).
2. Import creates a new family from a previously exported file, validating
   the payload and remapping IDs.

**Administration** *(Phase 6)*
1. `admin` users can list/deactivate users and change global roles via the
   admin UI. Bootstrap admin is created from env on first startup.

## 7. Roadmap

| Phase | Deliverable | Status |
|-------|-------------|--------|
| 1 | Foundation: monorepo, docs, Docker, CI/CD, DB schema, auth slice, frontend shell | ✅ Done |
| 2 | Families CRUD + collaborators (API + UI) | Next |
| 3 | Persons + relationships CRUD, interactive tree, ancestors/descendants | |
| 4 | Photo uploads (persons), member search | |
| 5 | Export/import JSON; GEDCOM evaluation | |
| 6 | Admin UI, email verification/password reset, notification polish | |
| 7 | Hardening: integration/E2E tests, CSP, backups automation, monitoring | |

When a phase is completed: update the table above, the feature table in §4,
and the Status section in `AGENTS.md`.
