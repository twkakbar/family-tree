# ADR-0006: TailwindCSS v4 CSS-first configuration

**Status**: Accepted (2026-07)

## Context

Tailwind v4 (stable since January 2025) replaces `tailwind.config.js` with
CSS-first configuration (`@import "tailwindcss"` + `@theme` tokens) and a
dedicated PostCSS plugin, with significantly faster builds.

## Decision

Use Tailwind v4 with all design tokens defined in
`frontend/src/app/globals.css` under `@theme`. No `tailwind.config` file.
Class merging via `cn()` (clsx + tailwind-merge) in every UI primitive.

## Consequences

- One place (globals.css) defines the design system; primitives in
  `components/ui/` consume it.
- Contributors must not add a `tailwind.config.js` — customization happens
  in CSS. Most online snippets still show v3 syntax; translate accordingly.
