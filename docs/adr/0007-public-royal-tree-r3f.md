# ADR-0007: Public Aceh royal tree + React Three Fiber

- **Status:** Accepted
- **Date:** 2026-07-22

## Context

The product pivoted from a private multi-family genealogy app to a **public
Aceh royal family tree**: anyone may view; only approved editors and
super-admins may mutate. The primary experience is an immersive **3D**
viewer with animation. The previous 2D SVG canvas cannot deliver camera
orbit, depth, or the intended cinematic motion without reinventing WebGL.

## Decision

1. **Singleton family** — keep the existing `families` / `persons` /
   `relationships` / `family_collaborators` schema; seed one royal family
   at startup; public APIs always resolve that family. Multi-family product
   UX is retired.
2. **Public read APIs** under `/api/v1/public/*` (no JWT). Mutations remain
   authenticated and require global `admin` or an `editor`/`owner`
   collaborator grant on the royal family.
3. **Edit access requests** — registered users submit a request; super-admin
   approves (creates `editor` collaborator) or rejects. No email delivery
   in this phase.
4. **3D frontend** — add `three`, `@react-three/fiber`, and
   `@react-three/drei`. Justification: stdlib and existing deps cannot
   render interactive WebGL trees; R3F is the maintained React binding for
   Three.js and keeps the viewer inside the React tree (hooks, Suspense).

## Consequences

- New env: `ROYAL_FAMILY_NAME`, `ROYAL_FAMILY_DESCRIPTION` (optional).
- Frontend landing + `/explore` are public; `/studio` and `/admin/*` are
  protected.
- Bundle size increases for the explore route; use dynamic import of the
  canvas so the landing page stays light.
- If WebGL is unavailable, show a simple non-3D fallback message (full 2D
  fallback can come later).
