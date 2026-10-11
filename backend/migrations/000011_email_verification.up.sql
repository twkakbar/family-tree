-- Registration email verification tokens and a send log for the hourly cap.
CREATE TABLE email_verification_tokens (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash text        NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX email_verification_tokens_token_hash_unique
    ON email_verification_tokens (token_hash);
CREATE INDEX email_verification_tokens_user_id_idx
    ON email_verification_tokens (user_id);

CREATE TABLE email_dispatches (
    id       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email    citext      NOT NULL,
    purpose  text        NOT NULL,
    sent_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX email_dispatches_email_purpose_sent_idx
    ON email_dispatches (email, purpose, sent_at);
