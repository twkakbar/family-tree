-- Phase 4: trigram GIN index for person name search within a family.
-- pg_trgm was enabled in 000001; this index makes ILIKE / % lookups fast.
CREATE INDEX persons_name_trgm_idx ON persons
    USING gin ((first_name || ' ' || last_name) gin_trgm_ops)
    WHERE deleted_at IS NULL;
