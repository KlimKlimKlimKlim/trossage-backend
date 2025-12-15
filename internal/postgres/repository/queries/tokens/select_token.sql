SELECT
    id,
    user_id,
    token_hash,
    expires_at,
    COALESCE(revoked_at, '0001-01-01'::timestamptz) AS revoked_at,
    created_at
FROM
    refresh_tokens
WHERE
    user_id = $1
    AND token_hash = $2;

