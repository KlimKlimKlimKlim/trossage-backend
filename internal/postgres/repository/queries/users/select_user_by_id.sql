SELECT
    id,
    login,
    password_hash,
    display_name,
    created_at,
    updated_at,
    deleted_at
FROM
    users
WHERE
    id = $1;

