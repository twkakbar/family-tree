-- pgcrypto: gen_random_uuid() for UUID primary key defaults.
-- citext:   case-insensitive emails.
-- pg_trgm:  trigram indexes for person name search (used from Phase 4).
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
