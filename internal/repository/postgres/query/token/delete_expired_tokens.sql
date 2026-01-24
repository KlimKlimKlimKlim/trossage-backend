DELETE FROM refresh_tokens
WHERE expires_at < $1
    AND revoked_at IS NULL;

