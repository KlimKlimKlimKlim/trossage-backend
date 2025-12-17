SELECT
    COUNT(*)
FROM
    users
WHERE
    deleted_at IS NULL
    AND login LIKE $1 || '%';

