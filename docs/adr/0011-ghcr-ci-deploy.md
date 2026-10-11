# ADR-0011: Build images in GitHub Actions and publish to GHCR

- **Status:** Accepted
- **Date:** 2026-10-11

## Context

Production was first started by building images on the VPS (`docker compose
up --build`). That uses the small VM's CPU and RAM for compiles, and it
drifts from what CI tested.

`deploy.yml` originally pushed to Docker Hub (`twkakbar/family-tree-*`) on
git tags. No Docker Hub token is available to Actions. GHCR accepts
`GITHUB_TOKEN` from the same repository.

## Decision

1. **Build only in GitHub Actions**, after the `CI` workflow succeeds on
   `main` (or when someone runs Deploy manually).
2. **Publish to GHCR**: `ghcr.io/twkakbar/family-tree-backend|frontend|nginx`,
   tagged `sha-<7>` and `latest`. Packages are public so the VM can pull
   without a personal access token.
3. **The VM only pulls.** Deploy SSHes with a dedicated key, fast-forwards
   `/opt/family-tree` to `origin/main`, points `.env` image names at GHCR,
   then `docker compose -f docker-compose.yml -f docker-compose.prod.yml pull`
   and `up -d`. No `--build` on the server.
4. Docker Hub remains unused. A later move back needs a new ADR.

## Consequences

- A failed CI run does not deploy.
- The first deploy after this ADR waits for CI, including the Compose e2e job.
- Public images expose the build output, not secrets. Secrets stay in the
  VM `.env`, which deploy does not overwrite except the three `*_IMAGE` lines.
- Renewing or rotating the deploy key is a GitHub secret plus
  `authorized_keys` change, not a password in the workflow.
