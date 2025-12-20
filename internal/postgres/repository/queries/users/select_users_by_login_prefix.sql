SELECT
    id,
    user_login,
    display_name,
    created_at
FROM
    users
WHERE
    deleted_at IS NULL
    AND user_login LIKE $1 || '%'
ORDER BY
    user_login
LIMIT $2 OFFSET $3;

