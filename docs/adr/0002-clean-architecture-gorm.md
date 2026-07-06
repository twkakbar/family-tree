# ADR-0002: Clean Architecture with GORM tags on domain entities

**Status**: Accepted (2026-07)

## Context

Textbook Clean Architecture keeps the domain 100% framework-free, which
requires a separate persistence model per entity plus mapping code in both
directions. For a CRUD-heavy application of this size that duplication buys
little and costs a lot of boilerplate.

## Decision

Keep the Clean Architecture *dependency rule* (handlers/repositories depend
on domain; services contain all business logic; repository interfaces live
in domain) but allow GORM struct tags and `gorm.DeletedAt` on domain
entities, which double as persistence models.

## Consequences

- One struct per entity; repositories stay thin; soft delete is automatic.
- The domain package imports `gorm.io/gorm` (types only). Accepted trade-off;
  services still never touch SQL/GORM APIs, so swapping persistence remains
  a repository-layer change.
- HTTP DTOs remain separate (`internal/handler`) — entities are never
  serialized to clients.
