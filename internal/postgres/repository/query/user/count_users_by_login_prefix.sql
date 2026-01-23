SELECT
    COUNT(*) AS total
FROM
    users
WHERE
    id != $1
    AND deleted_at IS NULL
    AND user_login LIKE $2 || '%';

