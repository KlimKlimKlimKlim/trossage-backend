DELETE FROM refresh_tokens
WHERE revoked_at IS NOT NULL
    AND revoked_at < $1;

