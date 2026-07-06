-- Typed edges between persons of the same family.
--   parent: directed, from = parent, to = child.
--   spouse: undirected, stored once in canonical order (from < to).
-- Edges are hard-deleted: persons are soft-deleted, so the service layer
-- must delete a person's edges in the same transaction that soft-deletes
-- the person (FK CASCADE does not fire on soft delete).
CREATE TABLE relationships (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id      uuid        NOT NULL REFERENCES families (id) ON DELETE CASCADE,
    from_person_id uuid        NOT NULL REFERENCES persons (id) ON DELETE CASCADE,
    to_person_id   uuid        NOT NULL REFERENCES persons (id) ON DELETE CASCADE,
    type           text        NOT NULL CHECK (type IN ('parent', 'spouse')),
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT relationships_no_self_edge CHECK (from_person_id <> to_person_id),
    CONSTRAINT relationships_spouse_canonical CHECK (
        type <> 'spouse' OR from_person_id < to_person_id
    ),
    CONSTRAINT relationships_unique UNIQUE (from_person_id, to_person_id, type)
);

CREATE INDEX relationships_family_id_idx ON relationships (family_id);
CREATE INDEX relationships_from_person_idx ON relationships (from_person_id);
CREATE INDEX relationships_to_person_idx ON relationships (to_person_id);
