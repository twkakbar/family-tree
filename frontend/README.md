# frontend/

Next.js 15 (App Router) client for Family Tree. Read `../ARCHITECTURE.md` §4
and `../DEVELOPMENT_RULES.md` §4 before changing code. The **auth flow is the
reference implementation** — copy its structure for new features.

## Layout

```
src/app/               routes: / (landing), (auth)/login|register, (dashboard)/dashboard
src/components/ui/     design-system primitives (Button, Input, Card, …)
src/components/layout/ shell components (Header, Logo)
src/hooks/             React Query hooks (use-auth)
src/services/          API service layer — the ONLY axios callers
src/lib/               api-client (interceptors + silent refresh), query-client,
                       query-keys, validations (zod), env, auth-token store
src/types/             API envelope + entity types
src/utils/             pure helpers (cn)
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
| `axios` | HTTP client with interceptors for auth/refresh — project requirement |
| `zod` | Schema validation for forms — project requirement |
| `react-hook-form` + `@hookform/resolvers` | Form state + zod bridge — project requirement |
| `tailwindcss` + `@tailwindcss/postcss` | Styling — project requirement |
| `clsx` + `tailwind-merge` | Conditional classes with Tailwind conflict resolution (`cn()`) |
| `eslint`, `eslint-config-next`, `@eslint/eslintrc` | Linting (flat config) |
| `prettier` + `prettier-plugin-tailwindcss` | Formatting with sorted classes |

Adding a dependency requires justification per `AGENTS.md` rule 3.
