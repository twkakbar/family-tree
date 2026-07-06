# migrations/

Sequential SQL migrations, applied by [golang-migrate](https://github.com/golang-migrate/migrate).
The files are embedded into the API binary (`migrations.go`) and applied on
startup when `RUN_MIGRATIONS=true`, or manually via `cmd/migrate`.

## Rules

- Never edit a migration that may have been applied anywhere. Fix mistakes
  with a NEW migration.
- Every `.up.sql` has a working `.down.sql`.
- Follow the conventions in `DEVELOPMENT_RULES.md` §5 (UUID PKs,
  timestamptz, soft delete on user-facing entities, snake_case).
- Document schema changes in `docs/DATABASE.md` in the same PR.

## Commands (from `backend/`)

```bash
go run ./cmd/migrate new <snake_case_name>   # create numbered .up/.down pair
go run ./cmd/migrate up                      # apply all pending
go run ./cmd/migrate down 1                  # roll back N migrations
go run ./cmd/migrate version                 # show current version
go run ./cmd/migrate force <version>         # repair a dirty state (careful)
```

## Dirty state

If a migration fails halfway, golang-migrate marks the schema "dirty" and
refuses to continue. Inspect the database, fix the issue manually, then
`go run ./cmd/migrate force <last-good-version>` and re-run `up`.
See `docs/DATABASE.md` for details.
