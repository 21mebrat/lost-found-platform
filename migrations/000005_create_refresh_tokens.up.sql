CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token  TEXT NOT NULL UNIQUE,

    device_info TEXT,
    ip_address  INET,
    revoked     BOOLEAN NOT NULL DEFAULT false,
    revoked_at  TIMESTAMPTZ,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_revoked_pair CHECK (
        (revoked = false AND revoked_at IS NULL) OR
        (revoked = true AND revoked_at IS NOT NULL)
    )
);

CREATE INDEX idx_refresh_tokens_user_active ON refresh_tokens(user_id) WHERE revoked = false;

CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);