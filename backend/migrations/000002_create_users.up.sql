-- Accounts that can authenticate. NOT family members (see persons).
CREATE TABLE users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name          text        NOT NULL,
    email         citext      NOT NULL,
    password_hash text        NOT NULL,
    role          text        NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    deleted_at    timestamptz
);

-- Partial unique index: a soft-deleted account frees its email for re-registration.
CREATE UNIQUE INDEX users_email_active_unique ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX users_deleted_at_idx ON users (deleted_at);
