# backend/

Go API for the **Kesultanan Aceh** family tree. Gin + GORM + PostgreSQL,
Clean Architecture. Read `../ARCHITECTURE.md` §2–§3 and
`../DEVELOPMENT_RULES.md` §3 before changing code. The **auth slice** and
later feature services (persons, public tree, access requests, admin) are
the reference patterns — copy their structure for new features.

## Layout

```
cmd/api/          entry point: wiring, HTTP server, graceful shutdown
cmd/migrate/      migration CLI (up / down / version / force / new)
internal/domain/      entities, repository+service interfaces, domain errors
internal/repository/  GORM implementations, DB connection, migration runner
internal/service/     business logic + unit tests (in-memory fakes)
internal/handler/     Gin handlers, DTOs, domain-error → HTTP mapping
internal/middleware/  request-id, logging, recovery, CORS, auth (JWT), RBAC
internal/routes/      route registration
internal/config/      env-based configuration (the only env reader)
internal/validation/  custom validator tags (strong_password)
pkg/                  reusable: token, password, response, pagination, upload
migrations/           embedded SQL migrations (golang-migrate)
```

## Run natively (hybrid dev)

```bash
docker compose up postgres      # from repo root: database only
cp .env.example .env            # adjust if needed
go run ./cmd/api                # http://localhost:8080/healthz
```

## Test & lint

```bash
go test -race -cover ./...
go vet ./... && gofmt -l .
golangci-lint run               # optional extended linting
```

## Dependencies (why each exists)

| Module | Purpose |
|--------|---------|
| `gin-gonic/gin` | HTTP framework (routing, binding, middleware) — project requirement |
| `gorm.io/gorm` + `driver/postgres` | ORM over PostgreSQL — project requirement |
| `golang-migrate/migrate/v4` | Versioned SQL schema migrations (ADR-0003) |
| `golang-jwt/jwt/v5` | Sign/verify JWT access tokens (ADR-0004) |
| `google/uuid` | UUID primary keys and request ids |
| `golang.org/x/crypto` | bcrypt password hashing |
| `jackc/pgx/v5` | Postgres driver (transitively via GORM); used directly only to detect unique-violation errors |
| `joho/godotenv` | Loads `.env` in non-production for developer convenience |

Adding a dependency requires justification per `AGENTS.md` rule 3.
