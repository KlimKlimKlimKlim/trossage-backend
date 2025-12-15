CREATE TABLE refresh_tokens (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash varchar(64) NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz DEFAULT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_token ON refresh_tokens (user_id, token_hash);

CREATE INDEX idx_refresh_tokens_user_active ON refresh_tokens (user_id)
WHERE
    revoked_at IS NULL;

-- CREATE INDEX idx_refresh_tokens_cleanup ON refresh_tokens (expires_at, revoked_at);
