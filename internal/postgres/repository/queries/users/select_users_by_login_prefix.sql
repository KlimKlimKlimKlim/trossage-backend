SELECT
    id,
    login,
    display_name,
    created_at
FROM
    users
WHERE
    deleted_at IS NULL
    AND login LIKE $1 || '%'
ORDER BY
    login
LIMIT $2 OFFSET $3;

