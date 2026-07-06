# docker/

Infrastructure images that are not part of an application.

- `nginx/` — reverse proxy image. Routes `/api/*` to the Go backend,
  `/uploads/*` to the shared uploads volume, everything else to Next.js.
  Applies gzip, security headers, body-size limits and rate limiting
  (general `/api/` zone + stricter `/api/v1/auth/` zone).

Application images live next to their code: `backend/Dockerfile`,
`frontend/Dockerfile` (multi-stage, with `dev` and `runtime` targets).

Compose notes:

- `docker-compose.yml` is the **production** base; `docker-compose.dev.yml`
  adds development overrides. `.env`'s `COMPOSE_FILE` makes plain
  `docker compose up` use both (see `.env.example`).
- Frontend dependencies live in the `frontend_node_modules` named volume in
  dev. After changing `package.json`, refresh it with:
  `docker compose build frontend && docker compose up -V frontend`.
- If nginx ever caches a stale upstream after container recreation:
  `docker compose restart nginx` (mitigated by the `resolver` directive).
