-- Stores SHA-256 hashes of opaque refresh tokens (never the raw token).
-- Rows are revoked (revoked_at) on rotation/logout and may be purged by a
-- future cleanup job. Hard delete only; no soft-delete semantics.
CREATE TABLE refresh_tokens (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash text        NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX refresh_tokens_token_hash_unique ON refresh_tokens (token_hash);
CREATE INDEX refresh_tokens_user_id_idx ON refresh_tokens (user_id);
CREATE INDEX refresh_tokens_expires_at_idx ON refresh_tokens (expires_at);
