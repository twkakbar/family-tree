-- Persons are family tree nodes (family members). Pure data records:
-- no login, unrelated to users. Scoped to exactly one family.
CREATE TABLE persons (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id  uuid        NOT NULL REFERENCES families (id) ON DELETE CASCADE,
    first_name text        NOT NULL,
    last_name  text        NOT NULL DEFAULT '',
    gender     text        NOT NULL DEFAULT 'unknown' CHECK (gender IN ('male', 'female', 'other', 'unknown')),
    birth_date date,
    death_date date,
    bio        text        NOT NULL DEFAULT '',
    photo_path text        NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,
    CONSTRAINT persons_death_after_birth CHECK (
        death_date IS NULL OR birth_date IS NULL OR death_date >= birth_date
    )
);

CREATE INDEX persons_family_id_idx ON persons (family_id);
CREATE INDEX persons_family_name_idx ON persons (family_id, last_name, first_name);
CREATE INDEX persons_deleted_at_idx ON persons (deleted_at);
