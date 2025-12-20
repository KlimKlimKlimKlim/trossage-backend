SELECT
    COUNT(*) AS total
FROM
    users
WHERE
    deleted_at IS NULL
    AND user_login LIKE $1 || '%';

