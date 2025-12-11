SELECT
    id,
    login,
    password_hash,
    display_name,
    created_at,
    updated_at
FROM
    users
WHERE
    login = $1;

