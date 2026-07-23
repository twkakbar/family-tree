# frontend/

Next.js 15 (App Router) client for the **Kesultanan Aceh** public tree.
Read `../ARCHITECTURE.md` §4 and `../DEVELOPMENT_RULES.md` §4 before
changing code. Prefer existing feature slices (auth, public, studio, admin)
as patterns for new UI.

## Layout

```
src/app/               routes: / (3D explore), /about, (auth)/*, (dashboard)/studio|admin|…
src/components/ui/     design-system primitives (Button, Input, Dialog, …)
src/components/layout/ shell (Header, Logo)
src/components/tree/   public R3F canvas + studio 2D diagram
src/components/persons/ studio forms, detail popup, export/import
src/hooks/             React Query hooks (use-auth, use-public, use-persons, …)
src/services/          API service layer — the ONLY axios callers
src/lib/               api-client, query-keys, validations, env, auth-token
src/types/             API envelope + entity types
src/utils/             pure helpers (cn, tree-layout*)
src/middleware.ts      cookie-presence route guard (UX only)
```

Data flow (strict): component → hook → service → api-client. Forms use
React Hook Form + Zod (`lib/validations/`).

## Run

```bash
cp .env.example .env.local
npm install
npm run dev        # http://localhost:3000 (proxies /api to localhost:8080)
```

Checks: `npm run lint` · `npm run typecheck` · `npm run build`.

## Dependencies (why each exists)

| Package | Purpose |
|---------|---------|
| `next`, `react`, `react-dom` | Framework — project requirement |
| `@tanstack/react-query` (+devtools) | Server-state cache — project requirement |
| `axios` | HTTP client with refresh interceptors — project requirement |
| `zod`, `react-hook-form`, `@hookform/resolvers` | Forms + validation — project requirement |
| `tailwind-merge`, `clsx` | ClassName helpers for UI primitives |
| `@react-three/fiber`, `@react-three/drei`, `three` | Public 3D tree (ADR-0007) |
