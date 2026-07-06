-- Access grants: which user can do what within a family.
-- The family creator receives an 'owner' row (service-layer rule).
-- Grants are hard-deleted when access is revoked.
CREATE TABLE family_collaborators (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id  uuid        NOT NULL REFERENCES families (id) ON DELETE CASCADE,
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role       text        NOT NULL DEFAULT 'viewer' CHECK (role IN ('owner', 'editor', 'viewer')),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT family_collaborators_unique UNIQUE (family_id, user_id)
);

CREATE INDEX family_collaborators_user_id_idx ON family_collaborators (user_id);
