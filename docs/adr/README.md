# Architecture Decision Records

Short, immutable records of significant decisions. When a decision changes,
write a new ADR that supersedes the old one — do not rewrite history.

Format: Status / Context / Decision / Consequences. Keep them under a page.

| # | Title | Status |
|---|-------|--------|
| [0001](0001-monorepo.md) | Monorepo with root-level AI context documents | Accepted |
| [0002](0002-clean-architecture-gorm.md) | Clean Architecture with GORM tags on domain entities | Accepted |
| [0003](0003-migrations.md) | golang-migrate with embedded SQL (AutoMigrate forbidden) | Accepted |
| [0004](0004-auth-tokens.md) | JWT access token + rotating opaque refresh cookie | Accepted |
| [0005](0005-same-origin.md) | Same-origin routing via nginx (prod) / Next rewrites (dev) | Accepted |
| [0006](0006-tailwind-v4.md) | TailwindCSS v4 CSS-first configuration | Accepted |
