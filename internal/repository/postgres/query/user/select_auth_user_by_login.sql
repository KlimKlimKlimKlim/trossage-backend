SELECT
    id,
    user_login,
    password_hash,
    display_name,
    created_at,
    updated_at
FROM
    users
WHERE
    user_login = $1
    AND deleted_at IS NULL;

