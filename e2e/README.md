# E2E smoke tests (Playwright)

Phase 7 smoke suite against the full Docker stack (nginx entrypoint).

## Prerequisites

```bash
# from repo root — production compose (or your running stack)
docker compose -f docker-compose.yml up -d --build
# wait until: curl -fsS http://127.0.0.1/readyz

cd e2e
npm ci
npx playwright install chromium
```

## Run

```bash
# default BASE_URL=http://127.0.0.1
npm test

# custom port / host
BASE_URL=http://127.0.0.1:8088 npm test
```

Or from root: `make e2e-install` then `BASE_URL=… make e2e`.

CI runs the same suite after bringing Compose up (see `.github/workflows/ci.yml`).
