SELECT
    id,
    user_login,
    display_name,
    created_at
FROM
    users
WHERE
    id != $1
    AND deleted_at IS NULL
    AND user_login LIKE $2 || '%'
ORDER BY
    user_login
LIMIT $3 OFFSET $4;

