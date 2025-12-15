UPDATE
    refresh_tokens
SET
    revoked_at = NOW()
WHERE
    id = $1
    AND revoked_at IS NULL;

