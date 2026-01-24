UPDATE
    users
SET
    display_name = 'Deleted User',
    password_hash = '',
    deleted_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;

