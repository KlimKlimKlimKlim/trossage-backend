SELECT
    id,
    user_login,
    password_hash,
    display_name,
    created_at,
    updated_at,
    COALESCE(deleted_at, '0001-01-01'::timestamptz) AS deleted_at
FROM
    users
WHERE
    id = $1
    AND deleted_at IS NULL;

