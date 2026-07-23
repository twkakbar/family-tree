-- Phase 6: password reset tokens + email verification flag (SMTP deferred).
ALTER TABLE users
    ADD COLUMN email_verified_at timestamptz;

-- Existing accounts are treated as verified so login is unchanged.
UPDATE users SET email_verified_at = created_at WHERE email_verified_at IS NULL;

CREATE TABLE password_reset_tokens (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash text        NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX password_reset_tokens_token_hash_unique ON password_reset_tokens (token_hash);
CREATE INDEX password_reset_tokens_user_id_idx ON password_reset_tokens (user_id);
CREATE INDEX password_reset_tokens_expires_at_idx ON password_reset_tokens (expires_at);
