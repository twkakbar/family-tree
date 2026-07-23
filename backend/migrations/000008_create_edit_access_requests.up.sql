-- Edit-access requests: registered users ask to become editors of the
-- singleton royal family. Super-admin approves/rejects in-app.
CREATE TABLE edit_access_requests (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    message     text        NOT NULL DEFAULT '',
    status      text        NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewed_by uuid        REFERENCES users (id) ON DELETE SET NULL,
    reviewed_at timestamptz,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

-- At most one pending request per user.
CREATE UNIQUE INDEX edit_access_requests_one_pending_per_user
    ON edit_access_requests (user_id)
    WHERE status = 'pending';

CREATE INDEX edit_access_requests_status_idx ON edit_access_requests (status);
CREATE INDEX edit_access_requests_user_id_idx ON edit_access_requests (user_id);
