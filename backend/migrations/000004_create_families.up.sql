-- A family is a container for one family tree.
-- created_by is RESTRICT: user accounts must be soft-deleted (kept as rows),
-- so hard-deleting a user with families is an application-level error.
CREATE TABLE families (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        text        NOT NULL,
    description text        NOT NULL DEFAULT '',
    created_by  uuid        NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    deleted_at  timestamptz
);

CREATE INDEX families_created_by_idx ON families (created_by);
CREATE INDEX families_deleted_at_idx ON families (deleted_at);
